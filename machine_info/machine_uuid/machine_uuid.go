package machine_uuid

import (
	"DistriAI-Node/utils/log_utils"
	"os"
	"strings"
)

type MachineUUID string

func GetInfoMachineUUID() (MachineUUID, error) {
	logs.Normal("Getting machine ID...")

	// Easy debugging
	mID, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		return "", err
	}

	return MachineUUID(strings.TrimSpace(string(mID))), nil
	// return MachineUUID("b711ebf34e474f4db43198e23a59d411"), nil
}
