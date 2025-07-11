package main

import (
	pb "cart/kitex_gen/cart/mapservice"
	"cart/rpc/basic/global"
	"cart/rpc/basic/model"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"cart/rpc/basic/dal"
)

// BaiduMapClient 百度地图客户端
type BaiduMapClient struct {
	APIKey  string
	APIHost string
}

// NewBaiduMapClient 创建百度地图客户端
func NewBaiduMapClient() *BaiduMapClient {
	// 从配置文件获取百度地图配置
	baiduConfig := &global.AppConf.BaiduMapConfig
	return &BaiduMapClient{
		APIKey:  baiduConfig.APIKey,
		APIHost: baiduConfig.APIHost,
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

// MapServiceImpl implements the last service interface defined in the IDL.
type MapServiceImpl struct {
	mapClient *BaiduMapClient
}

var mongoService dal.MongoService

// 初始化MongoDB服务
func init() {
	// 延迟初始化，确保global.MongoDB已经连接
	go func() {
		time.Sleep(2 * time.Second)
		mongoService = dal.NewMongoService()

		// 创建索引
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := mongoService.CreateIndexes(ctx); err != nil {
			log.Printf("创建MongoDB索引失败: %v", err)
		} else {
			log.Println("MongoDB索引创建成功")
		}
	}()
}

// logAPICall 记录API调用日志到MongoDB
func (s *MapServiceImpl) logAPICall(ctx context.Context, apiType string, requestParams, responseData interface{}, responseCode, responseTime int, isSuccess bool, errorMessage string) {
	if mongoService == nil {
		return // 如果MongoDB服务未初始化，跳过日志记录
	}

	// 将请求参数和响应数据转换为map
	var reqParamsMap, respDataMap map[string]interface{}

	if requestParams != nil {
		reqBytes, _ := json.Marshal(requestParams)
		json.Unmarshal(reqBytes, &reqParamsMap)
	}

	if responseData != nil {
		respBytes, _ := json.Marshal(responseData)
		json.Unmarshal(respBytes, &respDataMap)
	}

	// 创建日志记录
	logRecord := &model.MongoMapApiLog{
		ApiType:       apiType,
		RequestParams: reqParamsMap,
		ResponseData:  respDataMap,
		ResponseCode:  responseCode,
		ResponseTime:  responseTime,
		IsSuccess:     isSuccess,
		ErrorMessage:  errorMessage,
		CreatedAt:     time.Now(),
	}

	// 异步记录日志，不影响主业务逻辑
	go func() {
		logCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := mongoService.MapApiLog().Create(logCtx, logRecord); err != nil {
			log.Printf("记录API日志失败: %v", err)
		}
	}()
}

// convertRegionToRegionInfo 将数据库Region模型转换为Thrift RegionInfo
func convertRegionToRegionInfo(region *model.Region) *pb.RegionInfo {
	return &pb.RegionInfo{
		Code:    region.Code,
		Name:    region.Name,
		Pcode:   region.Pcode,
		Sname:   region.Sname,
		Level:   region.Level,
		Mername: region.Mername,
		Pinyin:  region.Pinyin,
	}
}

// GetLocationByAddress 根据地址获取坐标（地理编码）
func (c *BaiduMapClient) GetLocationByAddress(address string) (*GeoCodingResponse, error) {
	apiURL := fmt.Sprintf("%s/geocoding/v3/?address=%s&output=json&ak=%s",
		c.APIHost, url.QueryEscape(address), c.APIKey)

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
func (c *BaiduMapClient) GetAddressByLocation(lat, lng float64) (*ReverseGeoCodingResponse, error) {
	apiURL := fmt.Sprintf("%s/reverse_geocoding/v3/?ak=%s&output=json&coordtype=wgs84ll&location=%f,%f",
		c.APIHost, c.APIKey, lat, lng)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var result ReverseGeoCodingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &result, nil
}

// GetLocationByIP 根据IP获取位置信息
func (c *BaiduMapClient) GetLocationByIP(ip string) (*IPLocationResponse, error) {
	apiURL := fmt.Sprintf("%s/location/ip?ip=%s&ak=%s&coor=bd09ll",
		c.APIHost, ip, c.APIKey)

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

// calculateDistance 计算两点间距离（哈弗森公式）
func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371000 // 地球半径，单位：米

	// 转换为弧度
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

// GeoCoding 地理编码接口，将地址转换为坐标
func (s *MapServiceImpl) GeoCoding(ctx context.Context, req *pb.GeoCodingReq) (*pb.GeoCodingResp, error) {
	ctx = context.Background()
	startTime := time.Now()
	var responseCode int
	var isSuccess bool
	var errorMessage string

	if s.mapClient == nil {
		s.mapClient = NewBaiduMapClient()
	}

	if req.Address == "" {
		responseCode = 400
		isSuccess = false
		errorMessage = "地址参数不能为空"

		// 记录API调用日志
		s.logAPICall(ctx, "geocoding", req, nil, responseCode, int(time.Since(startTime).Milliseconds()), isSuccess, errorMessage)

		return &pb.GeoCodingResp{
			Code:    400,
			Message: "地址参数不能为空",
		}, nil
	}

	result, err := s.mapClient.GetLocationByAddress(req.Address)
	if err != nil {
		responseCode = 500
		isSuccess = false
		errorMessage = fmt.Sprintf("地理编码失败: %v", err)

		// 记录API调用日志
		s.logAPICall(ctx, "geocoding", req, nil, responseCode, int(time.Since(startTime).Milliseconds()), isSuccess, errorMessage)

		return &pb.GeoCodingResp{
			Code:    500,
			Message: fmt.Sprintf("地理编码失败: %v", err),
		}, nil
	}

	if result.Status != 0 {
		responseCode = 500
		isSuccess = false
		errorMessage = fmt.Sprintf("百度地图API返回错误状态: %d", result.Status)

		// 记录API调用日志
		s.logAPICall(ctx, "geocoding", req, result, responseCode, int(time.Since(startTime).Milliseconds()), isSuccess, errorMessage)

		return &pb.GeoCodingResp{
			Code:    500,
			Message: fmt.Sprintf("百度地图API返回错误状态: %d", result.Status),
		}, nil
	}

	responseCode = 200
	isSuccess = true
	errorMessage = ""

	response := &pb.GeoCodingResp{
		Code:    200,
		Message: "地理编码成功",
		Data: &pb.GeoCodingData{
			Address: req.Address,
			Lng:     result.Result.Location.Lng,
			Lat:     result.Result.Location.Lat,
			Precise: int32(result.Result.Precise),
			Level:   result.Result.Level,
		},
	}

	// 记录成功的API调用日志
	s.logAPICall(ctx, "geocoding", req, response, responseCode, int(time.Since(startTime).Milliseconds()), isSuccess, errorMessage)

	return response, nil
}

// ReverseGeoCoding 逆地理编码接口，将坐标转换为地址
func (s *MapServiceImpl) ReverseGeoCoding(ctx context.Context, req *pb.ReverseGeoCodingReq) (*pb.ReverseGeoCodingResp, error) {
	ctx = context.Background()

	if s.mapClient == nil {
		s.mapClient = NewBaiduMapClient()
	}

	if req.Lat == 0 || req.Lng == 0 {
		return &pb.ReverseGeoCodingResp{
			Code:    400,
			Message: "经纬度参数不能为空",
		}, nil
	}

	result, err := s.mapClient.GetAddressByLocation(req.Lat, req.Lng)
	if err != nil {
		return &pb.ReverseGeoCodingResp{
			Code:    500,
			Message: fmt.Sprintf("逆地理编码失败: %v", err),
		}, nil
	}

	if result.Status != 0 {
		return &pb.ReverseGeoCodingResp{
			Code:    500,
			Message: fmt.Sprintf("百度地图API返回错误状态: %d", result.Status),
		}, nil
	}

	return &pb.ReverseGeoCodingResp{
		Code:    200,
		Message: "逆地理编码成功",
		Data: &pb.ReverseGeoCodingData{
			Lng:              result.Result.Location.Lng,
			Lat:              result.Result.Location.Lat,
			FormattedAddress: result.Result.FormattedAddress,
			Business:         result.Result.Business,
			Citycode:         strconv.Itoa(result.Result.Citycode),
		},
	}, nil
}

// IPLocation IP定位接口，根据IP获取位置信息
func (s *MapServiceImpl) IPLocation(ctx context.Context, req *pb.IPLocationReq) (*pb.IPLocationResp, error) {
	ctx = context.Background()

	if s.mapClient == nil {
		s.mapClient = NewBaiduMapClient()
	}

	if req.Ip == "" {
		return &pb.IPLocationResp{
			Code:    400,
			Message: "IP参数不能为空",
		}, nil
	}

	result, err := s.mapClient.GetLocationByIP(req.Ip)
	if err != nil {
		return &pb.IPLocationResp{
			Code:    500,
			Message: fmt.Sprintf("IP定位失败: %v", err),
		}, nil
	}

	if result.Status != 0 {
		return &pb.IPLocationResp{
			Code:    500,
			Message: fmt.Sprintf("百度地图API返回错误状态: %d", result.Status),
		}, nil
	}

	// 解析坐标
	x, _ := strconv.ParseFloat(result.Content.Point.X, 64)
	y, _ := strconv.ParseFloat(result.Content.Point.Y, 64)

	return &pb.IPLocationResp{
		Code:    200,
		Message: "IP定位成功",
		Data: &pb.IPLocationData{
			Ip:            req.Ip,
			Address:       result.Content.Address,
			AddressDetail: result.Content.AddressDetail.Province + result.Content.AddressDetail.City + result.Content.AddressDetail.District,
			Point: &pb.Point{
				X: x,
				Y: y,
			},
		},
	}, nil
}

// DistanceCalculate 距离计算接口，计算两点间距离
func (s *MapServiceImpl) DistanceCalculate(ctx context.Context, req *pb.DistanceCalculateReq) (*pb.DistanceCalculateResp, error) {
	ctx = context.Background()

	if req.OriginLat == 0 || req.OriginLng == 0 || req.DestLat == 0 || req.DestLng == 0 {
		return &pb.DistanceCalculateResp{
			Code:    400,
			Message: "起点和终点经纬度参数不能为空",
		}, nil
	}

	// 计算距离
	distanceMeters := calculateDistance(req.OriginLat, req.OriginLng, req.DestLat, req.DestLng)
	distanceKm := distanceMeters / 1000

	return &pb.DistanceCalculateResp{
		Code:    200,
		Message: "距离计算成功",
		Data: &pb.DistanceData{
			Origin: &pb.LocationPoint{
				Lat: req.OriginLat,
				Lng: req.OriginLng,
			},
			Destination: &pb.LocationPoint{
				Lat: req.DestLat,
				Lng: req.DestLng,
			},
			DistanceMeters: distanceMeters,
			DistanceKm:     distanceKm,
		},
	}, nil
}

// GetCurrentLocation 获取当前位置接口，根据IP获取当前位置
func (s *MapServiceImpl) GetCurrentLocation(ctx context.Context, req *pb.GetCurrentLocationReq) (*pb.GetCurrentLocationResp, error) {
	ctx = context.Background()

	if s.mapClient == nil {
		s.mapClient = NewBaiduMapClient()
	}

	if req.ClientIp == "" {
		return &pb.GetCurrentLocationResp{
			Code:    400,
			Message: "客户端IP参数不能为空",
		}, nil
	}

	result, err := s.mapClient.GetLocationByIP(req.ClientIp)
	if err != nil {
		return &pb.GetCurrentLocationResp{
			Code:    500,
			Message: fmt.Sprintf("获取当前位置失败: %v", err),
		}, nil
	}

	if result.Status != 0 {
		return &pb.GetCurrentLocationResp{
			Code:    500,
			Message: fmt.Sprintf("百度地图API返回错误状态: %d", result.Status),
		}, nil
	}

	// 解析坐标
	x, _ := strconv.ParseFloat(result.Content.Point.X, 64)
	y, _ := strconv.ParseFloat(result.Content.Point.Y, 64)

	return &pb.GetCurrentLocationResp{
		Code:    200,
		Message: "获取当前位置成功",
		Data: &pb.IPLocationData{
			Ip:            req.ClientIp,
			Address:       result.Content.Address,
			AddressDetail: result.Content.AddressDetail.Province + result.Content.AddressDetail.City + result.Content.AddressDetail.District,
			Point: &pb.Point{
				X: x,
				Y: y,
			},
		},
	}, nil
}

// GetProvinces 获取省份列表接口，获取所有省份信息
func (s *MapServiceImpl) GetProvinces(ctx context.Context, req *pb.GetProvincesReq) (*pb.GetProvincesResp, error) {
	ctx = context.Background()

	var regions []model.Region

	// 查询所有省份（level = 1）
	if err := global.DB.Debug().Where("level = ?", 1).Find(&regions).Error; err != nil {
		return &pb.GetProvincesResp{
			Code:    500,
			Message: fmt.Sprintf("查询省份失败: %v", err),
		}, nil
	}

	// 转换为RegionInfo
	var regionInfos []*pb.RegionInfo
	for _, region := range regions {
		regionInfos = append(regionInfos, convertRegionToRegionInfo(&region))
	}

	return &pb.GetProvincesResp{
		Code:    200,
		Message: "获取省份列表成功",
		Data:    regionInfos,
	}, nil
}

// GetCities 获取城市列表接口，获取指定省份下的城市信息
func (s *MapServiceImpl) GetCities(ctx context.Context, req *pb.GetCitiesReq) (*pb.GetCitiesResp, error) {
	ctx = context.Background()

	if req.ProvinceCode == 0 {
		return &pb.GetCitiesResp{
			Code:    400,
			Message: "省份代码不能为空",
		}, nil
	}

	var regions []model.Region

	// 查询指定省份下的所有城市（level = 2）
	if err := global.DB.Debug().Where("pcode = ? AND level = ?", req.ProvinceCode, 2).Find(&regions).Error; err != nil {
		return &pb.GetCitiesResp{
			Code:    500,
			Message: fmt.Sprintf("查询城市失败: %v", err),
		}, nil
	}

	// 转换为RegionInfo
	var regionInfos []*pb.RegionInfo
	for _, region := range regions {
		regionInfos = append(regionInfos, convertRegionToRegionInfo(&region))
	}

	return &pb.GetCitiesResp{
		Code:    200,
		Message: "获取城市列表成功",
		Data:    regionInfos,
	}, nil
}

// GetDistricts 获取区县列表接口，获取指定城市下的区县信息
func (s *MapServiceImpl) GetDistricts(ctx context.Context, req *pb.GetDistrictsReq) (*pb.GetDistrictsResp, error) {
	ctx = context.Background()

	if req.CityCode == 0 {
		return &pb.GetDistrictsResp{
			Code:    400,
			Message: "城市代码不能为空",
		}, nil
	}

	var regions []model.Region

	// 查询指定城市下的所有区县（level = 3）
	if err := global.DB.Debug().Where("pcode = ? AND level = ?", req.CityCode, 3).Find(&regions).Error; err != nil {
		return &pb.GetDistrictsResp{
			Code:    500,
			Message: fmt.Sprintf("查询区县失败: %v", err),
		}, nil
	}

	// 转换为RegionInfo
	var regionInfos []*pb.RegionInfo
	for _, region := range regions {
		regionInfos = append(regionInfos, convertRegionToRegionInfo(&region))
	}

	return &pb.GetDistrictsResp{
		Code:    200,
		Message: "获取区县列表成功",
		Data:    regionInfos,
	}, nil
}

// GetRegionPath 获取区域路径接口，获取指定区域的完整路径
func (s *MapServiceImpl) GetRegionPath(ctx context.Context, req *pb.GetRegionPathReq) (*pb.GetRegionPathResp, error) {
	ctx = context.Background()

	if req.RegionCode == 0 {
		return &pb.GetRegionPathResp{
			Code:    400,
			Message: "区域代码不能为空",
		}, nil
	}

	// 获取当前区域
	var currentRegion model.Region
	if err := global.DB.Debug().Where("code = ?", req.RegionCode).First(&currentRegion).Error; err != nil {
		return &pb.GetRegionPathResp{
			Code:    404,
			Message: "未找到指定区域",
		}, nil
	}

	regionPath := &pb.RegionPath{}
	fullPathParts := []string{}

	// 根据级别获取完整路径
	switch currentRegion.Level {
	case 3: // 区县级
		// 获取区县信息
		regionPath.District = convertRegionToRegionInfo(&currentRegion)
		fullPathParts = append([]string{currentRegion.Name}, fullPathParts...)

		// 获取城市信息
		var cityRegion model.Region
		if err := global.DB.Where("code = ?", currentRegion.Pcode).First(&cityRegion).Error; err == nil {
			regionPath.City = convertRegionToRegionInfo(&cityRegion)
			fullPathParts = append([]string{cityRegion.Name}, fullPathParts...)

			// 获取省份信息
			var provinceRegion model.Region
			if err := global.DB.Where("code = ?", cityRegion.Pcode).First(&provinceRegion).Error; err == nil {
				regionPath.Province = convertRegionToRegionInfo(&provinceRegion)
				fullPathParts = append([]string{provinceRegion.Name}, fullPathParts...)
			}
		}

	case 2: // 市级
		// 获取城市信息
		regionPath.City = convertRegionToRegionInfo(&currentRegion)
		fullPathParts = append([]string{currentRegion.Name}, fullPathParts...)

		// 获取省份信息
		var provinceRegion model.Region
		if err := global.DB.Where("code = ?", currentRegion.Pcode).First(&provinceRegion).Error; err == nil {
			regionPath.Province = convertRegionToRegionInfo(&provinceRegion)
			fullPathParts = append([]string{provinceRegion.Name}, fullPathParts...)
		}

	case 1: // 省级
		// 获取省份信息
		regionPath.Province = convertRegionToRegionInfo(&currentRegion)
		fullPathParts = append([]string{currentRegion.Name}, fullPathParts...)
	}

	regionPath.FullPath = strings.Join(fullPathParts, "/")

	return &pb.GetRegionPathResp{
		Code:    200,
		Message: "获取区域路径成功",
		Data:    regionPath,
	}, nil
}

// SearchRegion 搜索区域接口，根据关键词搜索区域信息
func (s *MapServiceImpl) SearchRegion(ctx context.Context, req *pb.SearchRegionReq) (*pb.SearchRegionResp, error) {
	ctx = context.Background()

	if req.Keyword == "" {
		return &pb.SearchRegionResp{
			Code:    400,
			Message: "搜索关键词不能为空",
		}, nil
	}

	var regions []model.Region
	query := global.DB.Debug().Model(&model.Region{})

	// 如果指定了父级区域代码，则在该区域内搜索
	if req.ParentCode != nil && *req.ParentCode != 0 {
		query = query.Where("pcode = ?", *req.ParentCode)
	}

	// 按名称、简称、拼音搜索
	query = query.Where("name LIKE ? OR sname LIKE ? OR pinyin LIKE ?",
		"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")

	// 限制结果数量，避免返回过多数据
	if err := query.Limit(50).Find(&regions).Error; err != nil {
		return &pb.SearchRegionResp{
			Code:    500,
			Message: fmt.Sprintf("搜索区域失败: %v", err),
		}, nil
	}

	// 转换为RegionInfo
	var regionInfos []*pb.RegionInfo
	for _, region := range regions {
		regionInfos = append(regionInfos, convertRegionToRegionInfo(&region))
	}

	return &pb.SearchRegionResp{
		Code:    200,
		Message: "搜索区域成功",
		Data:    regionInfos,
	}, nil
}
