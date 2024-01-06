package cmd

import (
	"DistriAI-Node/chain"
	"DistriAI-Node/chain/distri"
	"DistriAI-Node/chain/subscribe"
	"DistriAI-Node/config"
	"DistriAI-Node/core_task"
	"DistriAI-Node/docker"
	"DistriAI-Node/machine_info"
	"DistriAI-Node/machine_info/disk"
	"DistriAI-Node/machine_info/machine_uuid"
	"DistriAI-Node/pattern"
	logs "DistriAI-Node/utils/log_utils"
	"encoding/json"
	"fmt"
	"time"

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
				distriWrapper, hwInfo, chainInfo, err := getDistri(true)
				if err != nil {
					logs.Error(err.Error())
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

				subscribeBlocks := subscribe.NewSubscribeWrapper(chainInfo)

				for {
					time.Sleep(1 * time.Second)

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

					logs.Normal(fmt.Sprintf("Start workspace container, orderId: %v", fmt.Sprintf("%#x", order.OrderId)))

					isGPU := false
					if hwInfo.GPUInfo.Number > 0 {
						isGPU = true
					}
					containerID, err := docker.RunWorkspaceContainer(isGPU)
					if err != nil {
						logs.Error(fmt.Sprintln("RunWorkspaceContainer error: ", err))
						return nil
					}
					if core_task.StartTimer(distriWrapper, order) {
						core_task.OrderComplete(distriWrapper, order.Metadata, isGPU, containerID)
					} else {
						core_task.OrderFailed(distriWrapper, order.Metadata, order.Buyer, containerID)
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

				machine, err := distriWrapper.GetMachine()
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

func getDistri(isHw bool) (*distri.WrapperDistri, *machine_info.MachineInfo, *chain.InfoChain, error) {

	key := config.GlobalConfig.Base.PrivateKey

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

		diskInfo, err := disk.GetDiskInfo()
		if err != nil {
			return nil, nil, nil, err
		}

		isGPU := false
		if hwInfo.GPUInfo.Number > 0 {
			isGPU = true
		}
		// Easy debugging
		score, err := docker.RunScoreContainer(isGPU)
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
