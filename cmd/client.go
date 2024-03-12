package cmd

import (
	"DistriAI-Node/chain/subscribe"
	"DistriAI-Node/config"
	"DistriAI-Node/core"
	"DistriAI-Node/docker"
	"DistriAI-Node/nginx"
	"DistriAI-Node/pattern"
	"DistriAI-Node/server"
	"DistriAI-Node/utils"
	dbutils "DistriAI-Node/utils/db_utils"
	logs "DistriAI-Node/utils/log_utils"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc/ws"
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
				distriWrapper, hwInfo, chainInfo, err := core.GetDistri(true)
				if err != nil {
					logs.Error(fmt.Sprintf("GetDistri: %v", err))
					return nil
				}

				if err = nginx.StartNginx(
					config.GlobalConfig.Console.NginxPort,
					config.GlobalConfig.Console.ConsolePost,
					config.GlobalConfig.Console.ServerPost); err != nil {
					logs.Error(fmt.Sprintf("StartNginx error: %v", err))
					return nil
				}

				machine, err := distriWrapper.GetMachine()
				if err != nil {
					logs.Error(fmt.Sprintf("Error: %v", err))
					return nil
				}

				if machine.Metadata == "" {
					logs.Normal("Machine does not exist")
					hash, err := distriWrapper.AddMachine(*hwInfo)
					if err != nil {
						logs.Error(fmt.Sprintf("Error block : %v, msg : %v\n", hash, err))
						return nil
					}
				} else {
					logs.Normal("Machine already exists")
					if machine.Status.String() == "Renting" {
						logs.Error(fmt.Sprintf("Machine is Renting, status: %v", machine.Status.String()))
						return nil
					}
				}

				go server.StartServer(config.GlobalConfig.Console.ServerPost)

				core.StartHeartbeatTask(distriWrapper, hwInfo.MachineUUID)

				subscribeBlocks := subscribe.NewSubscribeWrapper(chainInfo)

				var oldOrder solana.PublicKey
				var containerID string
				for {

					logs.Normal("=============== Start subscription")
					order, err := subscribeBlocks.SubscribeEvents(hwInfo.MachineUUID)
					logs.Normal("=============== End subscription")
					if err != nil {
						logs.Error(err.Error())
						logs.Normal("Restart subscription")
						subscribeBlocks.Conn.WsClient.Close()
						subscribeBlocks.Conn.WsClient = nil
						time.Sleep(1 * time.Second)
						subscribeBlocks.Conn.WsClient, err = ws.Connect(context.Background(), pattern.WsRPC)
						if err != nil {
							logs.Error(fmt.Sprintf("ws.Connect restart error: %v", err))
							return nil
						}
						continue
					}

					if order.Metadata == "" {
						logs.Error("order metadata is empty")
						continue
					}

					spew.Dump(order)

					switch order.Status.String() {
					case "Training":
						if oldOrder.Equals(subscribeBlocks.ProgramDistriOrder) {
							continue
						}
						logs.Vital(fmt.Sprintf("Start order. OrderAccount: %v", subscribeBlocks.ProgramDistriOrder))

						isGPU := false
						if hwInfo.GPUInfo.Number > 0 {
							isGPU = true
						}

						mlToken, err := utils.GenerateRandomString(16)
						if err != nil {
							logs.Error(err.Error())
							return nil
						}

						db, err := dbutils.NewDB()
						if err != nil {
							logs.Error(err.Error())
							return nil
						}
						db.Update([]byte("buyer"), []byte(order.Buyer.String()))
						db.Update([]byte("token"), []byte(mlToken))
						db.Close()
						logs.Normal(fmt.Sprintf("From buyer: %v ; mlToken: %v", order.Buyer, mlToken))

						containerID, err = docker.RunWorkspaceContainer(isGPU, mlToken)
						if err != nil {
							if strings.Contains(err.Error(), "container already exists") {
								logs.Error(err.Error())
								if err = core.OrderFailed(distriWrapper, order.Metadata, order.Buyer); err != nil {
									logs.Error(err.Error())
									return nil
								}
								continue
							}
							logs.Error(fmt.Sprintln("RunWorkspaceContainer error: ", err))
							return nil
						}

						_, err = distriWrapper.OrderStart()
						if err != nil {
							logs.Error(fmt.Sprintf("OrderStart: %v", err))
							return nil
						}

						oldOrder = subscribeBlocks.ProgramDistriOrder
						subscribeBlocks.IsRunning = true
						core.StartOrderTimer(distriWrapper, order, isGPU, containerID)
					case "Refunded":
						if containerID == "" {
							continue
						}
						logs.Vital(fmt.Sprintf("Refunded order. OrderAccount: %v", subscribeBlocks.ProgramDistriOrder))

						if err = core.OrderRefunded(containerID); err != nil {
							return nil
						}
						containerID = ""
					default:
						logs.Error(fmt.Sprintf("Order status is not training or refunded, status: %v", order.Status.String()))
						continue
					}
				}
			},
		},
		{
			Name:  "stop",
			Usage: "Stop the client.",
			Action: func(c *cli.Context) error {
				distriWrapper, hwInfo, _, err := core.GetDistri(false)
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
