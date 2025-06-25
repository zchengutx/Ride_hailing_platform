package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BaiduRouteResp struct {
	Status int `json:"status"`
	Result struct {
		Routes []struct {
			Distance int `json:"distance"`
			Duration int `json:"duration"`
		} `json:"routes"`
	} `json:"result"`
}

// GetRouteFromBaidu 调用百度地图API获取路径距离和时长
func GetRouteFromBaidu(startLat, startLng, endLat, endLng float64, ak string) (distance, duration int, err error) {
	url := fmt.Sprintf(
		"https://api.map.baidu.com/directionlite/v1/driving?origin=%f,%f&destination=%f,%f&ak=%s",
		startLat, startLng, endLat, endLng, ak,
	)
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	var result BaiduRouteResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, 0, err
	}
	if result.Status != 0 || len(result.Result.Routes) == 0 {
		return 0, 0, fmt.Errorf("baidu api error: %d", result.Status)
	}
	return result.Result.Routes[0].Distance, result.Result.Routes[0].Duration, nil
}
