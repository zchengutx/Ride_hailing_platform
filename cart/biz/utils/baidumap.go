package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"

	mapservicepb "cart/kitex_gen/cart/mapservice"
	"cart/rpc/basic/model"
)

// 百度地图API配置
const (
	BaiduMapAPIHost = "https://api.map.baidu.com"
	// 需要在配置文件中设置您的百度地图API Key
	DefaultBaiduMapAPIKey = "mGQEH6rmm11Em096OUIIY4SUfUI5AD3u" // 您的百度地图API Key
)

// NewBaiduMapClient 创建百度地图客户端
func NewBaiduMapClient(apiKey, apiHost string) *model.BaiduMapClient {
	return &model.BaiduMapClient{
		APIKey:  apiKey,
		APIHost: apiHost,
	}
}

// Location 坐标结构
type Location struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

// AddressComponent 地址组件
type AddressComponent struct {
	Country      string `json:"country"`
	Province     string `json:"province"`
	City         string `json:"city"`
	District     string `json:"district"`
	Street       string `json:"street"`
	StreetNumber string `json:"street_number"`
	Adcode       string `json:"adcode"`
}

// GeoCodingResponse 地理编码响应
type GeoCodingResponse struct {
	Status int `json:"status"`
	Result struct {
		Location   Location `json:"location"`
		Precise    int      `json:"precise"`
		Confidence int      `json:"confidence"`
		Level      string   `json:"level"`
	} `json:"result"`
}

// ReverseGeoCodingResponse 逆地理编码响应
type ReverseGeoCodingResponse struct {
	Status int `json:"status"`
	Result struct {
		Location           Location         `json:"location"`
		FormattedAddress   string           `json:"formatted_address"`
		Business           string           `json:"business"`
		AddressComponent   AddressComponent `json:"addressComponent"`
		Pois               []interface{}    `json:"pois"`
		Roads              []interface{}    `json:"roads"`
		Poiregions         []interface{}    `json:"poiregions"`
		SematicDescription string           `json:"sematic_description"`
		Citycode           int              `json:"citycode"`
	} `json:"result"`
}

// IPLocationResponse IP定位响应
type IPLocationResponse struct {
	Status  int    `json:"status"`
	Address string `json:"address"`
	Content struct {
		Address       string `json:"address"`
		AddressDetail struct {
			Province     string `json:"province"`
			City         string `json:"city"`
			District     string `json:"district"`
			Street       string `json:"street"`
			StreetNumber string `json:"street_number"`
			CityCode     int    `json:"city_code"`
		} `json:"address_detail"`
		Point struct {
			X string `json:"x"`
			Y string `json:"y"`
		} `json:"point"`
	} `json:"content"`
}

// GetLocationByAddress 根据地址获取坐标（地理编码）
func GetLocationByAddress(c *model.BaiduMapClient, address string) (*GeoCodingResponse, error) {
	apiURL := fmt.Sprintf("%s/geocoding/v3/?address=%s&output=json&ak=%s",
		BaiduMapAPIHost, url.QueryEscape(address), c.APIKey)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var result GeoCodingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &result, nil
}

// GetAddressByLocation 根据坐标获取地址（逆地理编码）
func GetAddressByLocation(c *model.BaiduMapClient, lat, lng float64) (*ReverseGeoCodingResponse, error) {
	apiURL := fmt.Sprintf("%s/reverse_geocoding/v3/?ak=%s&output=json&coordtype=wgs84ll&location=%f,%f",
		BaiduMapAPIHost, c.APIKey, lat, lng)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 新增日志输出百度返回内容
	fmt.Println("[百度地图API返回]", string(body))

	var result ReverseGeoCodingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &result, nil
}

// GetLocationByIP 根据IP获取位置信息
func GetLocationByIP(c *model.BaiduMapClient, ip string) (*IPLocationResponse, error) {
	apiURL := fmt.Sprintf("%s/location/ip?ip=%s&ak=%s&coor=bd09ll",
		BaiduMapAPIHost, ip, c.APIKey)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var result IPLocationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &result, nil
}

// GetDefaultBaiduMapClient 获取默认百度地图客户端
func GetDefaultBaiduMapClient() *model.BaiduMapClient {
	return NewBaiduMapClient("", "")
}

// ConvertRegionToRegionInfo 将数据库Region模型转换为Thrift RegionInfo
func ConvertRegionToRegionInfo(region *model.Region) *mapservicepb.RegionInfo {
	return &mapservicepb.RegionInfo{
		Code:    region.Code,
		Name:    region.Name,
		Pcode:   region.Pcode,
		Sname:   region.Sname,
		Level:   region.Level,
		Mername: region.Mername,
		Pinyin:  region.Pinyin,
	}
}

// CalculateDistance 计算两点间距离（哈弗辛公式）
func CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371000 // 地球半径，单位：米
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
