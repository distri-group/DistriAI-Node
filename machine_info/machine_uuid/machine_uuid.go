package machine_uuid

import (
	"DistriAI-Node/utils/log_utils"
	"fmt"
	"os"
	"strings"
)

type MachineUUID string

func GetInfoMachineUUID() (MachineUUID, error) {
	logs.Normal("Getting machine ID...")

	// Easy debugging
	mID, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		return "", fmt.Errorf("> ReadFile: %v", err)
	}
	id := strings.TrimSpace(string(mID))
	// id = "ec25eab7e9b2f62b32f2ac249a0fd3b0"
	return MachineUUID(id), nil
}
