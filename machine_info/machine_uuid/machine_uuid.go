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
	return MachineUUID(strings.TrimSpace(string(mID))), nil
	// return MachineUUID("105fc89c16980c81fd642af19eacbb18"), nil
}
