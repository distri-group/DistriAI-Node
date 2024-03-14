package control

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

func OrderFailed(distri *distri.WrapperDistri, metadata string, buyer solana.PublicKey) error {
	logs.Normal("Order is failed")

	var orderPlacedMetadata pattern.OrderPlacedMetadata

	err := json.Unmarshal([]byte(metadata), &orderPlacedMetadata)
	if err != nil {
		return fmt.Errorf("> json.Unmarshal: %v", err.Error())
	}

	orderPlacedMetadata.MachineAccounts = distri.ProgramDistriMachine.String()

	_, err = distri.OrderFailed(buyer, orderPlacedMetadata)
	if err != nil {
		return fmt.Errorf("> distri.OrderFailed: %v", err.Error())
	}
	return nil
}

func GetDistri(isHw bool) (*distri.WrapperDistri, *machine_info.MachineInfo, error) {

	key := config.GlobalConfig.Base.PrivateKey

	machineUUID, err := machine_uuid.GetInfoMachineUUID()
	if err != nil {
		return nil, nil, fmt.Errorf("> GetInfoMachineUUID: %v", err)
	}

	newConfig := config.NewConfig(
		key,
		pattern.RPC,
		pattern.WsRPC)

	var chainInfo *chain.InfoChain
	chainInfo, err = chain.GetChainInfo(newConfig, machineUUID)
	if err != nil {
		return nil, nil, fmt.Errorf("> GetChainInfo: %v", err)
	}

	var hwInfo machine_info.MachineInfo

	if isHw {
		hwInfo, err = machine_info.GetMachineInfo()
		if err != nil {
			return nil, nil, fmt.Errorf("> GetMachineInfo: %v", err)
		}

		diskInfo, err := disk.GetDiskInfo()
		if err != nil {
			return nil, nil, err
		}

		isGPU := false
		// Easy debugging
		if hwInfo.GPUInfo.Number > 0 {
			isGPU = true
		}
		score, err := docker.RunScoreContainer(isGPU)
		if err != nil {
			return nil, nil, err
		}

		imageWorkspace := pattern.ML_WORKSPACE_NAME
		if isGPU {
			imageWorkspace = pattern.ML_WORKSPACE_GPU_NAME
		}
		if err = docker.ImageExistOrPull(imageWorkspace); err != nil {
			return nil, nil, err
		}

		hwInfo.Score = score
		hwInfo.MachineAccounts = chainInfo.ProgramDistriMachine.String()
		hwInfo.DiskInfo = diskInfo
	}

	hwInfo.Addr = chainInfo.Wallet.Wallet.PublicKey().String()
	hwInfo.MachineUUID = machineUUID

	jsonData, _ := json.Marshal(hwInfo)
	logs.Normal(fmt.Sprintf("Hardware Info : %v", string(jsonData)))

	return distri.NewDistriWrapper(chainInfo), &hwInfo, nil
}

// var oldDuration time.Time
// var orderTimer *time.Timer
type OrderControl struct {
	distri      *distri.WrapperDistri
	oldDuration *time.Time
	orderTimer  *time.Timer
}

func NewOrderControl(distri *distri.WrapperDistri) *OrderControl {
	return &OrderControl{
		distri: distri,
	}
}

func (orderControl OrderControl) OrderRefunded(containerID string) error {
	logs.Normal("Order is refunded")
	if err := docker.StopWorkspaceContainer(containerID); err != nil {
		return err
	}
	orderControl.orderTimer.Stop()
	return nil
}

func (orderControl OrderControl) CheckOrder(isGPU bool, containerID string) {
	newOrder, err := orderControl.distri.GetOrder()
	if err != nil {
		logs.Error(fmt.Sprintf("GetOrder Error: %v", err))
		return
	}

	newDuration := time.Unix(newOrder.StartTime, 0).Add(time.Hour * time.Duration(newOrder.Duration))

	logs.Normal(fmt.Sprintf("CheckOrder oldDuration: %v", orderControl.oldDuration))
	logs.Normal(fmt.Sprintf("CheckOrder newDuration: %v", newDuration))

	if newDuration.After(*orderControl.oldDuration) {
		logs.Normal("Restart timer")
		orderControl.oldDuration = &newDuration
		orderControl.orderTimer.Reset(time.Until(*orderControl.oldDuration))
	} else {
		if err = OrderComplete(orderControl.distri, newOrder.Metadata, isGPU, containerID); err != nil {
			logs.Error(fmt.Sprintf("OrderComplete Error: %v", err))
			return
		}
	}
}

func (orderControl OrderControl) StartOrderTimer(order distri_ai.Order, isGPU bool, containerID string) {

	// duration := time.Unix(order.OrderTime, 0).Add(time.Hour * time.Duration(order.Duration))
	now := time.Now()
	duration := now.Add(time.Hour * time.Duration(order.Duration)).Add(time.Second * 10)
	logs.Normal(fmt.Sprintf("Order start: %v", now.Format("2006-01-02 15:04:05")))
	logs.Normal(fmt.Sprintf("Order duration: %v", time.Hour*time.Duration(order.Duration)))
	logs.Normal(fmt.Sprintf("Order end time: %v", duration))

	orderControl.oldDuration = &duration
	orderControl.orderTimer = time.AfterFunc(time.Until(*orderControl.oldDuration), func() {
		orderControl.CheckOrder(isGPU, containerID)
	})
}

// temp
func StartHeartbeatTask(distri *distri.WrapperDistri, machineID machine_uuid.MachineUUID) {
	ticker := time.NewTicker(6 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				taskID, err := utils.GenerateRandomString(16)
				if err != nil {
					logs.Error(err.Error())
				}
				taskUuid, err := utils.ParseTaskUUID(string(taskID))
				if err != nil {
					logs.Error(fmt.Sprintf("error parsing taskUuid: %v", err))
				}

				machineUuid, err := utils.ParseMachineUUID(string(machineID))
				if err != nil {
					logs.Error(fmt.Sprintf("error parsing machineUuid: %v", err))
				}

				hash, err := distri.SubmitTask(taskUuid, machineUuid, utils.CurrentPeriod(), pattern.TaskMetadata{})
				if err != nil {
					logs.Error(fmt.Sprintf("Error block : %v, msg : %v\n", hash, err))
				}
			}
		}
	}()
}
