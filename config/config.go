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
		OuterNetPort  string `yaml:"outerNetPort"`
		NginxPort     string `yaml:"nginxPost"`
		ConsolePost   string `yaml:"consolePost"`
		ServerPost    string `yaml:"serverPost"`
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

	if GlobalConfig.Console.ServerPost == "" {
		GlobalConfig.Console.ServerPost = "8088"
	}
	if GlobalConfig.Console.ConsolePost == "" {
		GlobalConfig.Console.ConsolePost = "8080"
	}
	if GlobalConfig.Console.NginxPort == "" {
		GlobalConfig.Console.NginxPort = "80"
	}
	if GlobalConfig.Console.OuterNetPort == "" {
		GlobalConfig.Console.OuterNetPort = GlobalConfig.Console.NginxPort
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
