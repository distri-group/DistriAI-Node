package cmd

import (
	"DistriAI-Node/config"
	"DistriAI-Node/control"
	"DistriAI-Node/docker"
	"DistriAI-Node/nginx"
	"DistriAI-Node/pattern"
	"DistriAI-Node/server"
	"DistriAI-Node/utils"
	dbutils "DistriAI-Node/utils/db_utils"
	logs "DistriAI-Node/utils/log_utils"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/urfave/cli"
)

var ClientCommand = cli.Command{
	Name:  "node",
	Usage: "Starting or terminating a node program.",
	Subcommands: []cli.Command{
		{
			Name:  "start",
			Usage: "Upload hardware configuration and initiate listening events.",
			Action: func(c *cli.Context) error {
				distriWrapper, hwInfo, err := control.GetDistri(true)
				if err != nil {
					logs.Error(fmt.Sprintf("GetDistri: %v", err))
					return nil
				}

				if err = nginx.StartNginx(
					config.GlobalConfig.Console.NginxPort,
					config.GlobalConfig.Console.WorkPort,
					config.GlobalConfig.Console.ServerPort); err != nil {
					logs.Error(fmt.Sprintf("StartNginx error: %v", err))
					return nil
				}

				machine, err := distriWrapper.GetMachine()
				if err != nil {
					logs.Error(fmt.Sprintf("GetMachine: %v", err))
					return nil
				}

				if machine.Metadata == "" {
					logs.Normal("Machine does not exist")
					_, err := distriWrapper.AddMachine(*hwInfo)
					if err != nil {
						logs.Error(fmt.Sprintf("AddMachine: %v", err))
						return nil
					}
				} else {
					logs.Normal("Machine already exists")
				}

				go server.StartServer(config.GlobalConfig.Console.ServerPort)

				control.StartHeartbeatTask(distriWrapper, hwInfo.MachineUUID)

				for {
					time.Sleep(1 * time.Minute)

					machine, err = distriWrapper.GetMachine()
					if err != nil {
						logs.Error(fmt.Sprintf("GetMachine: %v", err))
						continue
					}

				ListenLoop:
					switch machine.Status.String() {
					case "Idle":
						// TODO: Add the logic of the Idle status.
						break ListenLoop
					case "ForRent":
						// TODO: Add the logic of the ForRent status.
						break ListenLoop
					case "Renting":

						logs.Normal(fmt.Sprintf("Machine is Renting, Details: %v", machine))

						orderID := machine.OrderPda
						if orderID.Equals(solana.SystemProgramID) {
							logs.Error(fmt.Sprintf("machine OrderPda error, OrderPda: %v", orderID))
							break ListenLoop
						}

						distriWrapper.ProgramDistriOrder = orderID
						newOrder, err := distriWrapper.GetOrder()
						if err != nil {
							logs.Error(fmt.Sprintf("GetOrder Error: %v", err))
							break ListenLoop
						}

						var orderPlacedMetadata pattern.OrderPlacedMetadata

						err = json.Unmarshal([]byte(newOrder.Metadata), &orderPlacedMetadata)
						if err != nil {
							logs.Error(fmt.Sprintf("json.Unmarshal: %v", err))
							break ListenLoop
						}

						isGPU := false
						var containerID string

						switch orderPlacedMetadata.OrderInfo.Intent {
						case "train":
							if hwInfo.GPUInfo.Number > 0 {
								isGPU = true
							}

							mlToken, err := dbutils.GenToken(newOrder.Buyer.String())
							if err != nil {
								logs.Error(fmt.Sprintf("GenToken: %v", err))
								break ListenLoop
							}
							logs.Normal(fmt.Sprintf("From buyer: %v ; mlToken: %v", newOrder.Buyer, mlToken))

							containerID, err = docker.RunWorkspaceContainer(isGPU, mlToken)
							if err != nil {
								if strings.Contains(err.Error(), "container already exists") {
									logs.Error(err.Error())
									if err = control.OrderFailed(distriWrapper, newOrder.Metadata, newOrder.Buyer); err != nil {
										logs.Error(fmt.Sprintf("control.OrderFailed: %v", err))
										break ListenLoop
									}
									break ListenLoop
								}
								logs.Error(fmt.Sprintln("RunWorkspaceContainer error: ", err))
								break ListenLoop
							}

							url := orderPlacedMetadata.OrderInfo.DownloadURL
							if len(url) > 0 {
								var downloadURL []utils.DownloadURL
								for _, u := range url {
									downloadURL = append(downloadURL, utils.DownloadURL{
										URL:      u,
										Checksum: "",
									})
								}

								logs.Normal("Downloading files...")
								err = utils.DownloadFiles(config.GlobalConfig.Console.WorkDirectory+"/ml-workspace", downloadURL)
								if err != nil {
									logs.Error(fmt.Sprintf("DownloadFiles: %v", err))
									break ListenLoop
								}
								logs.Normal("Download completed")
							}
						case "deploy":
							logs.Normal(fmt.Sprintf("Deploying, DownloadURL details: %v", orderPlacedMetadata.OrderInfo.DownloadURL))

							_, err := dbutils.GenToken(newOrder.Buyer.String())
							if err != nil {
								logs.Error(fmt.Sprintf("GenToken: %v", err))
								break ListenLoop
							}

							containerID, err = docker.RunDeployContainer(isGPU, orderPlacedMetadata.OrderInfo.DownloadURL)
							logs.Normal(fmt.Sprintf("DeployContainerID: %v", containerID))

							if err != nil {
								if strings.Contains(err.Error(), "container already exists") {
									logs.Error(err.Error())
									if err = control.OrderFailed(distriWrapper, newOrder.Metadata, newOrder.Buyer); err != nil {
										logs.Error(fmt.Sprintf("control.OrderFailed: %v", err))
										break ListenLoop
									}
									break ListenLoop
								}
								logs.Error(fmt.Sprintln("RunDeployContainer error: ", err))
								break ListenLoop
							}
						default:
							break ListenLoop
						}

						_, err = distriWrapper.OrderStart()
						if err != nil {
							logs.Error(fmt.Sprintf("OrderStart: %v", err))
							break ListenLoop
						}

						for {
							time.Sleep(1 * time.Minute)

							newOrder, err = distriWrapper.GetOrder()
							if err != nil {
								logs.Error(fmt.Sprintf("GetOrder Error: %v", err))
								break ListenLoop
							}

							switch newOrder.Status.String() {
							case "Preparing":
								logs.Error(fmt.Sprintf("Order error, ID: %v\norder: %v", distriWrapper.ProgramDistriOrder, newOrder))
								break ListenLoop
							case "Training":
								orderEndTime := time.Unix(newOrder.StartTime, 0).Add(time.Hour * time.Duration(newOrder.Duration))
								timeNow := time.Now()
								if timeNow.After(orderEndTime) {

									logs.Normal(fmt.Sprintf("Order completed, Details: %v", newOrder))

									if err = control.OrderComplete(distriWrapper, newOrder.Metadata, isGPU, containerID); err != nil {
										logs.Error(fmt.Sprintf("OrderComplete: %v", err))
									}
									break ListenLoop
								}
								continue
							case "Completed":
								logs.Error(fmt.Sprintf("Order error, ID: %v\norder: %v", distriWrapper.ProgramDistriOrder, newOrder))
								break ListenLoop
							case "Failed":
								logs.Error(fmt.Sprintf("Order error, ID: %v\norder: %v", distriWrapper.ProgramDistriOrder, newOrder))
								break ListenLoop
							case "Refunded":
								err = control.OrderRefunded(containerID)
								if err != nil {
									logs.Error(fmt.Sprintf("OrderRefunded: %v", err))
								}
								break ListenLoop
							}
						}
					default:
						logs.Error(fmt.Sprintf("machine status error, Status: %v", machine.Status))
						break ListenLoop
					}
				}
			},
		},
		{
			Name:  "stop",
			Usage: "Stop the client.",
			Action: func(c *cli.Context) error {
				distriWrapper, hwInfo, err := control.GetDistri(false)
				if err != nil {
					logs.Error(err.Error())
					return nil
				}

				hash, err := distriWrapper.RemoveMachine(*hwInfo)
				if err != nil {
					logs.Error(fmt.Sprintf("Error block : %v, msg : %v\n", hash, err))
				}

				db, err := dbutils.NewDB()
				if err != nil {
					logs.Error(err.Error())
				}
				db.Delete([]byte("buyer"))
				db.Delete([]byte("token"))
				db.Close()

				nginx.StopNginx()
				return nil
			},
		},
	},
}
