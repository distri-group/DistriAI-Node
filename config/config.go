package config

import (
	logs "DistriAI-Node/utils/log_utils"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Base struct {
		PrivateKey    string `yaml:"privateKey"`
		SecurityLevel string `yaml:"securityLevel"`
	} `yaml:"base"`
	Console struct {
		WorkDirectory string `yaml:"workDirectory"`
		OuterNetIP    string `yaml:"outerNetIP"`
		Port          string `yaml:"port"`
	} `yaml:"console"`
}

var GlobalConfig Config

func InitializeConfig() {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		logs.Error(fmt.Sprintf("Error reading config file: %v", err))
	}
	err = yaml.Unmarshal([]byte(data), &GlobalConfig)
	if err != nil {
		logs.Error(fmt.Sprintf("Error reading config file: %v", err))
	}
}

type SolanaConfig struct {
	Key   string
	RPC   string
	WsRPC string
}

func NewConfig(key string, rpc string, wsRPC string) *SolanaConfig {
	return &SolanaConfig{
		Key:   key,
		RPC:   rpc,
		WsRPC: wsRPC,
	}
}
