package location

import (
	"DistriAI-Node/utils/log_utils"
	"encoding/json"
	"io"
	"net/http"
)

type InfoLocation struct {
	Country string `json:"Country"`
	Region  string `json:"Region"`
	City    string `json:"City"`
	IP      string `json:"query"`
}

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
