package location

import (
	"DistriAI-Node/utils/log_utils"
	"encoding/json"
	"io"
	"net/http"
)

// InfoLocation 定义 InfoLocation 结构体
type InfoLocation struct {
	Country string `json:"Country"` // IP 国家
	Region  string `json:"Region"`  // IP 地区
	City    string `json:"City"`    // IP 城市
	IP      string `json:"query"`   // IP 地址
}

// GetLocationInfo 获取IP对应的并返回 InfoMemory 结构体
func GetLocationInfo() (InfoLocation, error) {
	logs.Normal("Getting location info...")

	resp, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return InfoLocation{}, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return InfoLocation{}, err
	}

	var response InfoLocation
	json.Unmarshal(body, &response)
	return response, nil
}
