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

    // Updates machine account information in the order metadata
    orderPlacedMetadata.MachineAccounts = distri.ProgramDistriMachine.String()

    // Notifies the distribution system that the order is completed
    _, err = distri.OrderCompleted(orderPlacedMetadata, isGPU)
    if err != nil {
        return err
    }
    return nil
}

func OrderFailed(distri *distri.WrapperDistri, orderPlacedMetadata pattern.OrderPlacedMetadata, buyer solana.PublicKey) error {
    // Logs the event of the order failing.
    logs.Normal("Order is failed")

    // Updates the order placement metadata with the machine account details on order failure.
    orderPlacedMetadata.MachineAccounts = distri.ProgramDistriMachine.String()

    // Attempts to flag the order as failed using the distri module's method.
    _, err := distri.OrderFailed(buyer, orderPlacedMetadata)
    // If the order fails to be flagged as failed, returns an error with details.
    if err != nil {
        return fmt.Errorf("> distri.OrderFailed: %v", err.Error())
    }
    // Returns nil indicating successful processing of the order failure.
    return nil
}

// GetDistri retrieves a distribution wrapper, machine information, and potential errors based on machine details and a flag indicating a long runtime.
func GetDistri(longTime bool) (*distri.WrapperDistri, *machine_info.MachineInfo, error) {
    // Initialize machine information structure
    var hwInfo machine_info.MachineInfo

    // Fetch basic machine info which may include CPU, GPU, disk, etc.
    hwInfo, err := machine_info.GetMachineInfo(longTime)
    if err != nil {
        return nil, nil, fmt.Errorf("> GetMachineInfo: %v", err)
    }

    // Obtain disk information and incorporate it into machine info
    diskInfo, err := disk.GetDiskInfo()
    if err != nil {
        return nil, nil, err
    }
    hwInfo.DiskInfo = diskInfo

    // For a long runtime, perform extra GPU checks and scoring
    if longTime {
        // Determine if GPU is present and run scoring container accordingly
        // Simplified debugging logic
        isGPU := hwInfo.GPUInfo.Number > 0

        score, err := docker.RunScoreContainer(isGPU)
        if err != nil {
            return nil, nil, err
        }

        // Select the appropriate ML workspace image based on GPU presence
        imageWorkspace := pattern.ML_WORKSPACE_NAME
        if isGPU {
            imageWorkspace = pattern.ML_WORKSPACE_GPU_NAME
        }
        // Ensure the required Docker image is available
        if err = docker.ImageExistOrPull(imageWorkspace); err != nil {
            return nil, nil, err
        }

        // Include scoring in machine information
        hwInfo.Score = score
    }

    // Create a new configuration object using the global private key and RPC settings
    key := config.GlobalConfig.Base.PrivateKey
    newConfig := config.NewConfig(key, config.GlobalConfig.Base.Rpc)

    // Retrieve chain information using machine hardware details and UUID
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

    chainInfo, err := chain.GetChainInfo(newConfig, machineUUID)
    if err != nil {
        return nil, nil, fmt.Errorf("> GetChainInfo: %v", err)
    }

    // Update machine info to include chain-related details
    hwInfo.MachineAccounts = chainInfo.ProgramDistriMachine.String()
    hwInfo.Addr = chainInfo.Wallet.Wallet.PublicKey().String() // Truncated for brevity
    hwInfo.MachineUUID = machineUUID

    // Log machine information after marshaling to JSON
    jsonData, _ := json.Marshal(hwInfo)
    logs.Normal(fmt.Sprintf("Hardware Info : %v", string(jsonData)))

    // Prepare directory for model creation
    modleCreatePath := pattern.ModleCreatePath
    if err = os.MkdirAll(modleCreatePath, 0755); err != nil {
        return nil, nil, fmt.Errorf("> MkdirAll: %v", err)
    }

    // Download model upload web static resources (ZIP file)
    modelURL := []utils.DownloadURL{
        {URL: config.GlobalConfig.Console.IpfsNodeUrl + "/ipfs" + utils.EnsureLeadingSlash("QmZQpwwUTne3rR1ZHfSTAwMQAsGChBBc7Mm8yHCb3QsEhE"),
         Checksum: "", Name: "DistriAI-Model-Create.zip"},
    }
    if err = utils.DownloadFiles(modleCreatePath, modelURL); err != nil {
        return nil, nil, fmt.Errorf("> DownloadFiles: %v", err)
    }

    // Unzip the downloaded model resources
    _, err = utils.Unzip(modleCreatePath+"/DistriAI-Model-Create.zip", modleCreatePath)
    if err != nil {
        return nil, nil, fmt.Errorf("> Unzip: %v", err)
    }
    logs.Normal("Model upload web static resources have been downloaded")

    // Return the distribution wrapper, machine info, and no error upon successful completion
    return distri.NewDistriWrapper(chainInfo), &hwInfo, nil
}

// var oldDuration time.Time
// var orderTimer *time.Timer

func OrderRefunded(containerID string) error {
    // Logs that the order refund process has started.
	logs.Normal("Order is refunded")
	
    // Calls the Docker service to stop the specified container.
    // If the stop operation fails, an error is returned.
	if err := docker.StopWorkspaceContainer(containerID); err != nil {
		return err
	}
	
    // Returns nil indicating a successful container stop operation.
	return nil
}

// temp
// StartHeartbeatTask initiates a periodic task that sends heartbeat signals to the distributed system.
// Heartbeat signals maintain the machine's presence and availability within the distributed system.
// Parameters:
// - distri (*distri.WrapperDistri): A wrapper for interacting with the distributed system.
// - machineID (machine_uuid.MachineUUID): A unique identifier for the machine sending the heartbeat.
func StartHeartbeatTask(distri *distri.WrapperDistri, machineID machine_uuid.MachineUUID) {
	// Create a new ticker that triggers every 6 hours.
	ticker := time.NewTicker(6 * time.Hour)
	// Launch a goroutine to handle the heartbeat task.
	go func() {
		for {
			select {
			case <-ticker.C:
				// Generate a random task ID.
				taskID, err := utils.GenerateRandomString(16)
				if err != nil {
					logs.Error(err.Error())
				}
				// Parse the generated task ID into a TaskUUID.
				taskUuid, err := utils.ParseTaskUUID(string(taskID))
				if err != nil {
					logs.Error(fmt.Sprintf("Error parsing taskUuid: %v", err))
				}
				// Parse the machine ID into a MachineUUID.
				machineUuid, err := utils.ParseMachineUUID(string(machineID))
				if err != nil {
					logs.Error(fmt.Sprintf("Error parsing machineUuid: %v", err))
				}
				// Submit the heartbeat task to the distributed system.
				hash, err := distri.SubmitTask(taskUuid, machineUuid, utils.CurrentPeriod(), pattern.TaskMetadata{})
				if err != nil {
					logs.Error(fmt.Sprintf("Error submitting task: Hash %v, Error: %v", hash, err))
				}
			}
			// The select loop continues indefinitely to keep the heartbeat task running.
		}
	}()
}