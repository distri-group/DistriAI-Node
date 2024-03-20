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
			return InfoIP{}, fmt.Errorf("> http.Get: %v", err)
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return InfoIP{}, fmt.Errorf("> io.ReadAll: %v", err)
		}
		json.Unmarshal(body, &response)
	}

	if !utils.CheckPort(config.GlobalConfig.Console.NginxPort) {
		return InfoIP{}, fmt.Errorf("> port %s is not available", config.GlobalConfig.Console.NginxPort)
	}
	if !utils.CheckPort(config.GlobalConfig.Console.WorkPort) {
		return InfoIP{}, fmt.Errorf("> port %s is not available", config.GlobalConfig.Console.WorkPort)
	}
	if !utils.CheckPort(config.GlobalConfig.Console.ServerPort) {
		return InfoIP{}, fmt.Errorf("> port %s is not available", config.GlobalConfig.Console.ServerPort)
	}

	response.Port = config.GlobalConfig.Console.OuterNetPort
	return response, nil
}
