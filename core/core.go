package core

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
	logs "DistriAI-Node/utils/log_utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
)

var oldDuration time.Time
var orderTimer *time.Timer

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

func OrderFailed(distri *distri.WrapperDistri, metadata string, buyer solana.PublicKey) error {
	logs.Normal("Order is failed")

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

func OrderRefunded(containerID string) error {
	logs.Normal("Order is refunded")
	if err := docker.StopWorkspaceContainer(containerID); err != nil {
		return err
	}
	orderTimer.Stop()
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

		isGPU := false
		// Easy debugging
		if hwInfo.GPUInfo.Number > 0 {
			isGPU = true
		}
		score, err := docker.RunScoreContainer(isGPU)
		if err != nil {
			return nil, nil, nil, err
		}

		imageWorkspace := pattern.ML_WORKSPACE_NAME
		if isGPU {
			imageWorkspace = pattern.ML_WORKSPACE_GPU_NAME
		}
		if err = docker.ImageExistOrPull(imageWorkspace); err != nil {
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

func CheckOrder(distri *distri.WrapperDistri, isGPU bool, containerID string) {
	newOrder, err := distri.GetOrder()
	if err != nil {
		logs.Error(fmt.Sprintf("GetOrder Error: %v", err))
		return
	}

	newDuration := time.Unix(newOrder.OrderTime, 0).Add(time.Hour * time.Duration(newOrder.Duration))

	logs.Normal(fmt.Sprintf("CheckOrder oldDuration: %v", oldDuration))
	logs.Normal(fmt.Sprintf("CheckOrder newDuration: %v", newDuration))

	if newDuration.After(oldDuration) {
		logs.Normal("Restart timer")
		oldDuration = newDuration
		orderTimer.Reset(time.Until(oldDuration))
	} else {
		if err = OrderComplete(distri, newOrder.Metadata, isGPU, containerID); err != nil {
			logs.Error(fmt.Sprintf("OrderComplete Error: %v", err))
			return
		}
	}
}

func StartTimer(distri *distri.WrapperDistri, order distri_ai.Order, isGPU bool, containerID string) {

	duration := time.Unix(order.OrderTime, 0).Add(time.Hour * time.Duration(order.Duration))
	logs.Normal(fmt.Sprintf("Order OrderTime: %v", time.Unix(order.OrderTime, 0)))
	logs.Normal(fmt.Sprintf("Order duration: %v", time.Hour*time.Duration(order.Duration)))
	logs.Normal(fmt.Sprintf("Order end time: %v", duration))

	oldDuration = duration
	orderTimer = time.AfterFunc(time.Until(oldDuration), func() {
		CheckOrder(distri, isGPU, containerID)
	})
}
