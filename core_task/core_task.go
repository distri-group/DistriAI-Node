package core_task

import (
	"DistriAI-Node/chain"
	"DistriAI-Node/chain/distri"
	"DistriAI-Node/chain/distri/distri_ai"
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
	"time"

	"github.com/gagliardetto/solana-go"
)

func OrderComplete(distri *distri.WrapperDistri, metadata string, isGPU bool, containerID string) error {
	logs.Normal("Order is complete")

	if err := docker.StopWorkspaceContainer(containerID); err != nil {
		return err
	}

	var orderPlacedMetadata pattern.OrderPlacedMetadata

	err := json.Unmarshal([]byte(metadata), &orderPlacedMetadata)
	if err != nil {
		return err
	}

	orderPlacedMetadata.MachineAccounts = distri.ProgramDistriMachine.String()

	_, err = distri.OrderCompleted(orderPlacedMetadata, isGPU)
	if err != nil {
		return err
	}
	return nil
}

func OrderFailed(distri *distri.WrapperDistri, metadata string, buyer solana.PublicKey, containerID string) error {
	logs.Normal("Order is failed")

	if err := docker.StopWorkspaceContainer(containerID); err != nil {
		return err
	}

	var orderPlacedMetadata pattern.OrderPlacedMetadata

	err := json.Unmarshal([]byte(metadata), &orderPlacedMetadata)
	if err != nil {
		return err
	}

	orderPlacedMetadata.MachineAccounts = distri.ProgramDistriMachine.String()

	_, err = distri.OrderFailed(buyer, orderPlacedMetadata)
	if err != nil {
		return err
	}
	return nil
}

func GetDistri(isHw bool) (*distri.WrapperDistri, *machine_info.MachineInfo, *chain.InfoChain, error) {

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

		if config.GlobalConfig.Console.Port != "" {
			config.GlobalConfig.Console.Port = "8080"
		}
		if !utils.CheckPort(config.GlobalConfig.Console.Port) {
			return nil, nil, nil, fmt.Errorf("port %s is not available", config.GlobalConfig.Console.Port)
		}

		isGPU := false
		// Easy debugging
		if hwInfo.GPUInfo.Number > 0 {
			isGPU = true
		}
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


func CheckOrder(done chan bool, distri *distri.WrapperDistri, oldDuration time.Time) {
	newOrder, err := distri.GetOrder()
	if err != nil {
		logs.Error(fmt.Sprintf("GetOrder Error: %v", err))
		done <- false
		return
	}

	newDuration := time.Unix(newOrder.OrderTime, 0).Add(time.Hour * time.Duration(newOrder.Duration))

	logs.Normal(fmt.Sprintf("CheckOrder newDuration: %v", newDuration))
	logs.Normal(fmt.Sprintf("CheckOrder oldDuration: %v", oldDuration))

	if newDuration.After(oldDuration) {
		logs.Normal("Restart timer")
		if !StartTimer(distri, newOrder) {
			done <- false
			return
		}
	}
	done <- true
}

func StartTimer(distri *distri.WrapperDistri, order distri_ai.Order) bool {
	done := make(chan bool)

	duration := time.Unix(order.OrderTime, 0).Add(time.Hour * time.Duration(order.Duration))
	logs.Normal(fmt.Sprintf("Order OrderTime: %v", time.Unix(order.OrderTime, 0)))
	logs.Normal(fmt.Sprintf("Order duration: %v", time.Hour*time.Duration(order.Duration)))
	logs.Normal(fmt.Sprintf("Order end time: %v", duration))
	time.AfterFunc(time.Until(duration), func() {
		CheckOrder(done, distri, duration)
	})
	return <-done
}
