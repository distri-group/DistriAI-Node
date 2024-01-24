package ip

import (
	"DistriAI-Node/config"
	"DistriAI-Node/utils"
	logs "DistriAI-Node/utils/log_utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type InfoIP struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

func GetIpInfo() (InfoIP, error) {
	logs.Normal("Getting outer net ip info...")

	var response InfoIP
	response.IP = config.GlobalConfig.Console.OuterNetIP
	if response.IP == "" {
		resp, err := http.Get("https://ipinfo.io")
		if err != nil {
			return InfoIP{}, err
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return InfoIP{}, err
		}
		json.Unmarshal(body, &response)
	}

	if config.GlobalConfig.Console.Port == "" {
		config.GlobalConfig.Console.Port = "8080"
	}
	if !utils.CheckPort(config.GlobalConfig.Console.Port) {
		return InfoIP{}, fmt.Errorf("port %s is not available", config.GlobalConfig.Console.Port)
	}
	if config.GlobalConfig.Console.OuterNetPort == "" {
		config.GlobalConfig.Console.OuterNetPort = config.GlobalConfig.Console.Port
	}
	response.Port = config.GlobalConfig.Console.OuterNetPort
	return response, nil
}
