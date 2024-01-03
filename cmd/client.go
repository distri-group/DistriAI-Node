package cmd

import (
	"DistriAI-Node/chain"
	"DistriAI-Node/chain/distri"
	"DistriAI-Node/chain/subscribe"
	"DistriAI-Node/config"
	"DistriAI-Node/docker"
	"DistriAI-Node/machine_info"
	"DistriAI-Node/machine_info/disk"
	"DistriAI-Node/machine_info/machine_uuid"
	"DistriAI-Node/pattern"
	"DistriAI-Node/utils"
	logs "DistriAI-Node/utils/log_utils"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
)

var ClientCommand = cli.Command{
	Name:  "client",
	Usage: "Starting or terminating a client.",
	Subcommands: []cli.Command{
		{
			Name:  "execute",
			Usage: "Upload hardware configuration and initiate listening events.",
			Action: func(c *cli.Context) error {
				distriWrapper, hwInfo, chainInfo, err := getDistri(true)
				if err != nil {
					logs.Error(err.Error())
					return nil
				}

				machine, err := distriWrapper.GetMachine(*hwInfo)
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

				subscribeBlocks := subscribe.NewSubscribeWrapper(chainInfo)

				for {
					time.Sleep(3 * time.Second)

					logs.Result("=============== Start subscription ==================")
					order, err := subscribeBlocks.SubscribeEvents(hwInfo)
					logs.Normal("=============== End subscription ==================")
					if err != nil {
						logs.Error(err.Error())
						return nil
					}

					if order.Metadata == "" {
						logs.Error("order metadata is empty")
						continue
					}

					if order.Status.String() != "Training" {
						logs.Error(fmt.Sprintf("Order status is not training, status: %v", order.Status.String()))
						continue
					}

					logs.Normal(fmt.Sprintf("Start training model, orderId: %v", fmt.Sprintf("%#x", order.OrderId)))

					var orderPlacedMetadata pattern.OrderPlacedMetadata

					err = json.Unmarshal([]byte(order.Metadata), &orderPlacedMetadata)
					if err != nil {
						logs.Error(fmt.Sprintf("error unmarshaling order metadata: %v", err))
						return nil
					}

					orderPlacedMetadata.MachineAccounts = chainInfo.ProgramDistriMachine.String()

					// Easy debugging
					switch orderPlacedMetadata.FormData.LibType {
					case "lib":
					case "docker":
						imageName := orderPlacedMetadata.FormData.ImageName + ":" + orderPlacedMetadata.FormData.ImageTag
						logs.Normal(fmt.Sprintf("imageName: %v", imageName))
					default:
						err = fmt.Errorf("libType error: %v", orderPlacedMetadata.FormData.LibType)
					}

					if err != nil {
						logs.Error(fmt.Sprintf("Container operation failed\n%v", err))
						_, err = distriWrapper.OrderFailed(order.Buyer, orderPlacedMetadata)
						if err != nil {
							logs.Error(err.Error())
							return err
						}
						logs.Normal("OrderFailed done")
					} else {

						/*
							In practice, only the "lib" mode requires the uploading of models on the client side, while the uploading logic of the "docker" mode is handled by the user (buyer).
							For the convenience of debugging and demonstration, the given model download IP address is used uniformly.
						*/
						orderPlacedMetadata.FormData.ModelUrl = "https://ipfs.io/ipfs/QmPHdMGXiuzxeQB5xyn6fqjKLGPUo4rxircpMzzd9cVomF?filename=model.pt"

						_, err = distriWrapper.OrderCompleted(orderPlacedMetadata)
						if err != nil {
							return err
						}
						logs.Normal("OrderCompleted done")
					}
				}
			},
		},
		{
			Name:  "stop",
			Usage: "Stop the client.",
			Action: func(c *cli.Context) error {
				distriWrapper, hwInfo, _, err := getDistri(false)
				if err != nil {
					logs.Error(err.Error())
					return err
				}

				machine, err := distriWrapper.GetMachine(*hwInfo)
				if err != nil {
					logs.Error(fmt.Sprintf("Error: %v", err))
					return nil
				}
				if machine.Status.String() != "Idle" {
					logs.Error(fmt.Sprintf("Machine is not idle, status: %v", machine.Status.String()))
					return nil
				}

				hash, err := distriWrapper.RemoveMachine(*hwInfo)
				if err != nil {
					logs.Error(fmt.Sprintf("Error block : %v, msg : %v\n", hash, err))
					return err
				}
				return nil
			},
		},
	},
}

func getFreeSpace() (disk.InfoDisk, error) {
	logs.Normal("Getting free space info...")

	dirpath := config.GlobalConfig.Base.WorkDirectory
	if _, err := os.Stat(dirpath); os.IsNotExist(err) {
		logs.Normal(fmt.Sprintf("%s does not exist. Using default directory /home", dirpath))
		dirpath = "/home"
	}

	freeSpace, err := utils.GetFreeSpace(dirpath)
	if err != nil {
		return disk.InfoDisk{}, fmt.Errorf("error calculating free space: %v", err)
	}

	diskInfo := disk.InfoDisk{
		Path:       dirpath,
		TotalSpace: float64(freeSpace) / 1024 / 1024 / 1024,
	}

	return diskInfo, nil
}

func getDistri(isHw bool) (*distri.WrapperDistri, *machine_info.MachineInfo, *chain.InfoChain, error) {

	key := config.GlobalConfig.Base.PrivateKey

	// 获取机器ID信息
	machineUUID, err := machine_uuid.GetInfoMachineUUID()
	if err != nil {
		return nil, nil, nil, err
	}

	newConfig := config.NewConfig(
		key,
		pattern.RPC,
		pattern.WsRPC)

	var chainInfo *chain.InfoChain
	chainInfo, err = chain.GetChainInfo(newConfig, machineUUID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting chain info: %v", err)
	}

	var hwInfo machine_info.MachineInfo

	if isHw {
		hwInfo, err = machine_info.GetMachineInfo()
		if err != nil {
			return nil, nil, nil, fmt.Errorf("error getting hardware info: %v", err)
		}

		diskInfo, err := getFreeSpace()
		if err != nil {
			return nil, nil, nil, err
		}

		// Easy debugging
		score, err := docker.RunScoreContainer()
		if err != nil {
			return nil, nil, nil, err
		}

		hwInfo.Score = score
		hwInfo.MachineAccounts = chainInfo.ProgramDistriMachine.String()
		hwInfo.DiskInfo = diskInfo
	}

	hwInfo.Addr = chainInfo.Wallet.Wallet.PublicKey().String()
	hwInfo.MachineUUID = machineUUID

	jsonData, _ := json.Marshal(hwInfo)
	logs.Normal(fmt.Sprintf("Hardware Info : %v", string(jsonData)))

	return distri.NewDistriWrapper(chainInfo), &hwInfo, chainInfo, nil
}
