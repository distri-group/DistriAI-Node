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

// OrderComplete marks the completion of an order process.
func OrderComplete(distri *distri.WrapperDistri, metadata string, isGPU bool, containerID string) error {
    logs.Normal("Order is complete")

    // Stops the container associated with the order
    if err := docker.StopWorkspaceContainer(containerID); err != nil {
        return err
    }

    // Parses the order placement metadata
    var orderPlacedMetadata pattern.OrderPlacedMetadata
    err := json.Unmarshal([]byte(metadata), &orderPlacedMetadata)
    if err != nil {
        return err
    }

    // Updates the machine account information in the order metadata
    orderPlacedMetadata.MachineAccounts = distri.ProgramDistriMachine.String()

    // Notifies the distribution system that the order is completed
    _, err = distri.OrderCompleted(orderPlacedMetadata, isGPU)
    if err != nil {
        return err
    }
    return nil
}

// OrderFailed Handles the scenario where an order has failed.
func OrderFailed(distri *distri.WrapperDistri, orderPlacedMetadata pattern.OrderPlacedMetadata, buyer solana.PublicKey) error {
    // Logs the information that the order has failed
    logs.Normal("Order is failed")

    // Updates the order placement metadata with the machine account details at the time of failure
    orderPlacedMetadata.MachineAccounts = distri.ProgramDistriMachine.String()

    // Attempts to mark the order as failed using the method from the distri module
    _, err := distri.OrderFailed(buyer, orderPlacedMetadata)
    if err != nil {
        return fmt.Errorf("> distri.OrderFailed: %v", err.Error())
    }
    return nil
}

// GetDistri retrieves distribution wrapper and machine information based on whether detailed hardware info is required.
// It collects machine stats, disk info, and optionally (if longTime is true) GPU details, container performance scores.
func GetDistri(longTime bool) (*distri.WrapperDistri, *machine_info.MachineInfo, error) {
	var hwInfo machine_info.MachineInfo

	// Retrieve basic machine information; return error if unsuccessful.
	hwInfo, err := machine_info.GetMachineInfo(longTime)
	if err != nil {
		return nil, nil, fmt.Errorf("> GetMachineInfo: %v", err)
	}

	// Gather disk information and attach it to the hardware info.
	diskInfo, err := disk.GetDiskInfo()
	if err != nil {
		return nil, nil, err
	}
	hwInfo.DiskInfo = diskInfo

	// Perform extended operations when longTime is true for debugging and scoring.
	if longTime {
		// Determine if the machine has a GPU.
		isGPU := hwInfo.GPUInfo.Number > 0

		// Run a container to assess performance and update hardware info with the score.
		score, err := docker.RunScoreContainer(isGPU)
		if err != nil {
			return nil, nil, err
		}
		hwInfo.Score = score

		// Handle Docker images based on GPU presence.
		// Ensure the required Docker image is available.
		imageWorkspace := pattern.ML_WORKSPACE_NAME
		if isGPU {
			imageWorkspace = pattern.ML_WORKSPACE_GPU_NAME
		}
		if err = docker.ImageExistOrPull(imageWorkspace); err != nil {
			return nil, nil, err
		}
	}

	// Generate a unique machine UUID using hardware specifics.
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

	// Configure a new instance with the private key and RPC settings.
	newConfig := config.NewConfig(
		config.GlobalConfig.Base.PrivateKey,
		config.GlobalConfig.Base.Rpc)

	// Fetch chain information using the new configuration and machine UUID.
	var chainInfo *chain.InfoChain
	chainInfo, err = chain.GetChainInfo(newConfig, machineUUID)
	if err != nil {
		return nil, nil, fmt.Errorf("> GetChainInfo: %v", err)
	}

	// Update machine info with chain details and UUID.
	hwInfo.MachineAccounts = chainInfo.ProgramDistriMachine.String()
	hwInfo.Addr = chainInfo.Wallet.Wallet.PublicKey().String()
	hwInfo.MachineUUID = machineUUID

	// Log the collected hardware information.
	jsonData, _ := json.Marshal(hwInfo)
	logs.Normal(fmt.Sprintf("Hardware Info : %v", string(jsonData)))

	// Prepare directory for model files and ensure its existence.
	modleCreatePath := pattern.ModleCreatePath
	if err = os.MkdirAll(modleCreatePath, 0755); err != nil {
		return nil, nil, fmt.Errorf("> MkdirAll: %v", err)
	}

	// Define and download model-related resources.
	var modelURL []utils.DownloadURL
	modelURL = append(modelURL, utils.DownloadURL{
		URL:      config.GlobalConfig.Console.IpfsNodeUrl + "/ipfs" + utils.EnsureLeadingSlash("QmZQpwwUTne3rR1ZHfSTAwMQAsGChBBc7Mm8yHCb3QsEhE"),
		Checksum: "",
		Name:     "DistriAI-Model-Create.zip",
	})
	if err = utils.DownloadFiles(modleCreatePath, modelURL); err != nil {
		return nil, nil, fmt.Errorf("> DownloadFiles: %v", err)
	}

	// Unzip the downloaded model file.
	if _, err = utils.Unzip(modleCreatePath+"/DistriAI-Model-Create.zip", modleCreatePath); err != nil {
		return nil, nil, fmt.Errorf("> Unzip: %v", err)
	}
	logs.Normal("Model upload web static resources have been downloaded")

	// Instantiate and return the distribution wrapper with the chain info and updated machine info.
	return distri.NewDistriWrapper(chainInfo), &hwInfo, nil
}

// var oldDuration time.Time
// var orderTimer *time.Timer

// OrderRefunded Handles the logic for processing a refunded order.
func OrderRefunded(containerID string) error {
    // Logs the event indicating that the order refund process has started.
	logs.Normal("Order is refunded")
	
    // Invokes the Docker service to stop the specified container.
    // If the stop operation fails, an error is returned.
	if err := docker.StopWorkspaceContainer(containerID); err != nil {
		return err
	}
	
    // Returns nil if the container stop operation is successful.
	return nil
}

// temp
// StartHeartbeatTask initiates a periodic heartbeat task to maintain the connection between the machine and the distributed system.
func StartHeartbeatTask(distri *distri.WrapperDistri, machineID machine_uuid.MachineUUID) {
	// Sets up a ticker that triggers the heartbeat task every 6 hours.
	ticker := time.NewTicker(6 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				// Generates a random task ID for the heartbeat task.
				taskID, err := utils.GenerateRandomString(16)
				if err != nil {
					// Logs an error if the task ID generation fails.
					logs.Error(err.Error())
				}
				// Parses the task ID into UUID format for task submission.
				taskUuid, err := utils.ParseTaskUUID(string(taskID))
				if err != nil {
					// Logs an error if parsing the task UUID fails.
					logs.Error(fmt.Sprintf("error parsing taskUuid: %v", err))
				}

				// Parses the machine ID into UUID format for task submission.
				machineUuid, err := utils.ParseMachineUUID(string(machineID))
				if err != nil {
					// Logs an error if parsing the machine UUID fails.
					logs.Error(fmt.Sprintf("error parsing machineUuid: %v", err))
				}

				// Submits the heartbeat task to the distributed system.
				// Parameters include the task UUID, machine UUID, current period, and empty task metadata.
				hash, err := distri.SubmitTask(taskUuid, machineUuid, utils.CurrentPeriod(), pattern.TaskMetadata{})
				if err != nil {
					// Logs an error if the task submission fails.
					logs.Error(fmt.Sprintf("Error block : %v, msg : %v\n", hash, err))
				}
			}
		}
	}()
}