package config

import (
	"DistriAI-Node/pattern"
	"DistriAI-Node/utils"
	logs "DistriAI-Node/utils/log_utils"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Base struct {
		Rpc           string `yaml:"rpc"`
		PrivateKey    string `yaml:"privateKey"`
		SecurityLevel string `yaml:"securityLevel"`
	} `yaml:"base"`
	Console struct {
		WorkDirectory string `yaml:"workDirectory"`
		IpfsNodeUrl   string `yaml:"ipfsNodeUrl"`
		OuterNetIP    string `yaml:"outerNetIP"`
		OuterNetPort  string `yaml:"outerNetPort"`
		NginxPort     string `yaml:"nginxPort"`
		WorkPort      string `yaml:"workPort"`
		ServerPort    string `yaml:"serverPort"`
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

	if GlobalConfig.Console.IpfsNodeUrl == "" {
		GlobalConfig.Console.IpfsNodeUrl = pattern.DefaultIpfsNode
	} else {
		GlobalConfig.Console.IpfsNodeUrl = utils.EnsureTrailingSlash(GlobalConfig.Console.IpfsNodeUrl)
	}
	if GlobalConfig.Console.ServerPort == "" {
		GlobalConfig.Console.ServerPort = "8088"
	}
	if GlobalConfig.Console.WorkPort == "" {
		GlobalConfig.Console.WorkPort = "8080"
	}
	if GlobalConfig.Console.NginxPort == "" {
		GlobalConfig.Console.NginxPort = "80"
	}
	if GlobalConfig.Console.OuterNetPort == "" {
		GlobalConfig.Console.OuterNetPort = GlobalConfig.Console.NginxPort
	}
	if GlobalConfig.Base.Rpc == "" {
		GlobalConfig.Base.Rpc = pattern.RPC
	}
}

type SolanaConfig struct {
	Key   string
	RPC   string
}

func NewConfig(key string, rpc string) *SolanaConfig {
	return &SolanaConfig{
		Key:   key,
		RPC:   rpc,
	}
}
