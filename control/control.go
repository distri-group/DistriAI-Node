package control

import (
	"DistriAI-Node/chain"
	"DistriAI-Node/chain/distri"
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

func OrderFailed(distri *distri.WrapperDistri, orderPlacedMetadata pattern.OrderPlacedMetadata, buyer solana.PublicKey) error {
	logs.Normal("Order is failed")

	orderPlacedMetadata.MachineAccounts = distri.ProgramDistriMachine.String()

	_, err := distri.OrderFailed(buyer, orderPlacedMetadata)
	if err != nil {
		return fmt.Errorf("> distri.OrderFailed: %v", err.Error())
	}
	return nil
}

func GetDistri(longTime bool) (*distri.WrapperDistri, *machine_info.MachineInfo, error) {

	var hwInfo machine_info.MachineInfo

	hwInfo, err := machine_info.GetMachineInfo(longTime)
	if err != nil {
		return nil, nil, fmt.Errorf("> GetMachineInfo: %v", err)
	}

	diskInfo, err := disk.GetDiskInfo()
	if err != nil {
		return nil, nil, err
	}
	hwInfo.DiskInfo = diskInfo

	if longTime {
		// Easy debugging
		isGPU := false
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
	}

	key := config.GlobalConfig.Base.PrivateKey

	machineUUID, err := machine_uuid.GetInfoMachineUUID(
		hwInfo.CPUInfo.ModelName,
		hwInfo.GPUInfo.Model,
		hwInfo.IpInfo.IP,
		hwInfo.LocationInfo.Country,
		hwInfo.LocationInfo.Region,
		hwInfo.LocationInfo.City)
	if err != nil {
		return nil, nil, fmt.Errorf("> GetInfoMachineUUID: %v", err)
	}

	newConfig := config.NewConfig(
		key,
		config.GlobalConfig.Base.Rpc)

	var chainInfo *chain.InfoChain
	chainInfo, err = chain.GetChainInfo(newConfig, machineUUID)
	if err != nil {
		return nil, nil, fmt.Errorf("> GetChainInfo: %v", err)
	}

	hwInfo.MachineAccounts = chainInfo.ProgramDistriMachine.String()
	hwInfo.Addr = chainInfo.Wallet.Wallet.PublicKey().String()
	hwInfo.MachineUUID = machineUUID

	jsonData, _ := json.Marshal(hwInfo)
	logs.Normal(fmt.Sprintf("Hardware Info : %v", string(jsonData)))

	modleCreatePath := pattern.ModleCreatePath
	err = os.MkdirAll(modleCreatePath, 0755)
	if err != nil {
		return nil, nil, fmt.Errorf("> MkdirAll: %v", err)
	}
	var modelURL []utils.DownloadURL
	modelURL = append(modelURL, utils.DownloadURL{
		URL:      config.GlobalConfig.Console.IpfsNodeUrl + "/ipfs" + utils.EnsureLeadingSlash("QmZQpwwUTne3rR1ZHfSTAwMQAsGChBBc7Mm8yHCb3QsEhE"),
		Checksum: "",
		Name:     "DistriAI-Model-Create.zip",
	})
	err = utils.DownloadFiles(modleCreatePath, modelURL)
	if err != nil {
		return nil, nil, fmt.Errorf("> DownloadFiles: %v", err)
	}
	_, err = utils.Unzip(modleCreatePath+"/DistriAI-Model-Create.zip", modleCreatePath)
	if err != nil {
		return nil, nil, fmt.Errorf("> Unzip: %v", err)
	}
	logs.Normal("Model upload web static resources have been downloaded")

	return distri.NewDistriWrapper(chainInfo), &hwInfo, nil
}

// var oldDuration time.Time
// var orderTimer *time.Timer

func OrderRefunded(containerID string) error {
	logs.Normal("Order is refunded")
	if err := docker.StopWorkspaceContainer(containerID); err != nil {
		return err
	}
	return nil
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
