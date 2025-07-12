package main

import (
	pb "cart/kitex_gen/cart/mapservice"
	"cart/rpc/basic/global"
	"cart/rpc/basic/model"
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"cart/biz/utils"
)

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
	mapClient *model.BaiduMapClient
}

// 相关调用处改为 utils.NewBaiduMapClient、utils.ConvertRegionToRegionInfo、utils.CalculateDistance

// 计算两点间距离
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

// GeoCoding 地理编码接口
// 将地址转换为经纬度坐标
func (s *MapServiceImpl) GeoCoding(ctx context.Context, req *pb.GeoCodingReq) (*pb.GeoCodingResp, error) {
	if s.mapClient == nil {
		cfg := &global.AppConf.BaiduMapConfig
		s.mapClient = utils.NewBaiduMapClient(cfg.APIKey, cfg.APIHost)
	}

	if req.Address == "" {
		return &pb.GeoCodingResp{
			Code:    400,
			Message: "地址参数不能为空",
		}, nil
	}

	// 修改为调用 utils.GetLocationByAddress
	result, err := utils.GetLocationByAddress((*model.BaiduMapClient)(s.mapClient), req.Address)
	if err != nil {
		return &pb.GeoCodingResp{
			Code:    500,
			Message: fmt.Sprintf("地理编码失败: %v", err),
		}, nil
	}

	if result.Status != 0 {
		return &pb.GeoCodingResp{
			Code:    500,
			Message: fmt.Sprintf("百度地图API返回错误状态: %d", result.Status),
		}, nil
	}

	return &pb.GeoCodingResp{
		Code:    200,
		Message: "成功",
		Data: &pb.GeoCodingData{
			Address: req.Address,
			Lng:     result.Result.Location.Lng,
			Lat:     result.Result.Location.Lat,
			Precise: int32(result.Result.Precise),
			Level:   result.Result.Level,
		},
	}, nil
}

// ReverseGeoCoding 逆地理编码接口
// 将经纬度坐标转换为地址信息
func (s *MapServiceImpl) ReverseGeoCoding(ctx context.Context, req *pb.ReverseGeoCodingReq) (*pb.ReverseGeoCodingResp, error) {
	if s.mapClient == nil {
		cfg := &global.AppConf.BaiduMapConfig
		s.mapClient = utils.NewBaiduMapClient(cfg.APIKey, cfg.APIHost)
	}

	if req.Lat == 0 || req.Lng == 0 {
		return &pb.ReverseGeoCodingResp{
			Code:    400,
			Message: "经纬度参数不能为空",
		}, nil
	}

	// 修改为调用 utils.GetAddressByLocation
	result, err := utils.GetAddressByLocation((*model.BaiduMapClient)(s.mapClient), req.Lat, req.Lng)
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
		Message: "成功",
		Data: &pb.ReverseGeoCodingData{
			Lng:              result.Result.Location.Lng,
			Lat:              result.Result.Location.Lat,
			FormattedAddress: result.Result.FormattedAddress,
			Business:         result.Result.Business,
			Citycode:         strconv.Itoa(result.Result.Citycode),
		},
	}, nil
}

// IPLocation IP定位接口
// 根据IP地址获取地理位置信息
func (s *MapServiceImpl) IPLocation(ctx context.Context, req *pb.IPLocationReq) (*pb.IPLocationResp, error) {
	if s.mapClient == nil {
		cfg := &global.AppConf.BaiduMapConfig
		s.mapClient = utils.NewBaiduMapClient(cfg.APIKey, cfg.APIHost)
	}

	if req.Ip == "" {
		return &pb.IPLocationResp{
			Code:    400,
			Message: "IP参数不能为空",
		}, nil
	}

	// 修改为调用 utils.GetLocationByIP
	result, err := utils.GetLocationByIP((*model.BaiduMapClient)(s.mapClient), req.Ip)
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
		Message: "成功",
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

// DistanceCalculate 距离计算接口
// 计算两个坐标点之间的距离
func (s *MapServiceImpl) DistanceCalculate(ctx context.Context, req *pb.DistanceCalculateReq) (*pb.DistanceCalculateResp, error) {
	if req.OriginLat == 0 || req.OriginLng == 0 || req.DestLat == 0 || req.DestLng == 0 {
		return &pb.DistanceCalculateResp{
			Code:    400,
			Message: "起点和终点经纬度参数不能为空",
		}, nil
	}

	// 计算距离
	distanceMeters := utils.CalculateDistance(req.OriginLat, req.OriginLng, req.DestLat, req.DestLng)
	distanceKm := distanceMeters / 1000

	return &pb.DistanceCalculateResp{
		Code:    200,
		Message: "成功",
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

// GetCurrentLocation 获取当前位置接口
// 根据客户端IP获取当前位置信息
func (s *MapServiceImpl) GetCurrentLocation(ctx context.Context, req *pb.GetCurrentLocationReq) (*pb.GetCurrentLocationResp, error) {
	if s.mapClient == nil {
		cfg := &global.AppConf.BaiduMapConfig
		s.mapClient = utils.NewBaiduMapClient(cfg.APIKey, cfg.APIHost)
	}

	if req.ClientIp == "" {
		return &pb.GetCurrentLocationResp{
			Code:    400,
			Message: "客户端IP参数不能为空",
		}, nil
	}

	// 修改为调用 utils.GetLocationByIP
	result, err := utils.GetLocationByIP((*model.BaiduMapClient)(s.mapClient), req.ClientIp)
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
		Message: "成功",
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

// GetProvinces 获取省份列表接口
// 查询所有省份信息
func (s *MapServiceImpl) GetProvinces(ctx context.Context, req *pb.GetProvincesReq) (*pb.GetProvincesResp, error) {
	var regions []model.Region

	// 查询所有省份（level = 1）
	if err := global.DB.Where("level = ?", 1).Find(&regions).Error; err != nil {
		return &pb.GetProvincesResp{
			Code:    500,
			Message: fmt.Sprintf("查询省份失败: %v", err),
		}, nil
	}

	// 转换为RegionInfo
	var regionInfos []*pb.RegionInfo
	for _, region := range regions {
		regionInfos = append(regionInfos, utils.ConvertRegionToRegionInfo(&region))
	}

	return &pb.GetProvincesResp{
		Code:    200,
		Message: "获取省份列表成功",
		Data:    regionInfos,
	}, nil
}

// GetCities 获取城市列表接口
// 根据省份代码查询该省份下的所有城市
func (s *MapServiceImpl) GetCities(ctx context.Context, req *pb.GetCitiesReq) (*pb.GetCitiesResp, error) {
	if req.ProvinceCode == 0 {
		return &pb.GetCitiesResp{
			Code:    400,
			Message: "省份代码不能为空",
		}, nil
	}

	var regions []model.Region

	// 查询指定省份下的所有城市（level = 2）
	if err := global.DB.Where("pcode = ? AND level = ?", req.ProvinceCode, 2).Find(&regions).Error; err != nil {
		return &pb.GetCitiesResp{
			Code:    500,
			Message: fmt.Sprintf("查询城市失败: %v", err),
		}, nil
	}

	// 转换为RegionInfo
	var regionInfos []*pb.RegionInfo
	for _, region := range regions {
		regionInfos = append(regionInfos, utils.ConvertRegionToRegionInfo(&region))
	}

	return &pb.GetCitiesResp{
		Code:    200,
		Message: "获取城市列表成功",
		Data:    regionInfos,
	}, nil
}

// GetDistricts 获取区县列表接口
// 根据城市代码查询该城市下的所有区县
func (s *MapServiceImpl) GetDistricts(ctx context.Context, req *pb.GetDistrictsReq) (*pb.GetDistrictsResp, error) {
	if req.CityCode == 0 {
		return &pb.GetDistrictsResp{
			Code:    400,
			Message: "城市代码不能为空",
		}, nil
	}

	var regions []model.Region

	// 查询指定城市下的所有区县（level = 3）
	if err := global.DB.Where("pcode = ? AND level = ?", req.CityCode, 3).Find(&regions).Error; err != nil {
		return &pb.GetDistrictsResp{
			Code:    500,
			Message: fmt.Sprintf("查询区县失败: %v", err),
		}, nil
	}

	// 转换为RegionInfo
	var regionInfos []*pb.RegionInfo
	for _, region := range regions {
		regionInfos = append(regionInfos, utils.ConvertRegionToRegionInfo(&region))
	}

	return &pb.GetDistrictsResp{
		Code:    200,
		Message: "获取区县列表成功",
		Data:    regionInfos,
	}, nil
}

// GetRegionPath 获取区域路径接口
// 根据区域代码获取完整的省市区路径信息
func (s *MapServiceImpl) GetRegionPath(ctx context.Context, req *pb.GetRegionPathReq) (*pb.GetRegionPathResp, error) {
	if req.RegionCode == 0 {
		return &pb.GetRegionPathResp{
			Code:    400,
			Message: "区域代码不能为空",
		}, nil
	}

	// 获取当前区域
	var currentRegion model.Region
	if err := global.DB.Where("code = ?", req.RegionCode).First(&currentRegion).Error; err != nil {
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
		regionPath.District = utils.ConvertRegionToRegionInfo(&currentRegion)
		fullPathParts = append([]string{currentRegion.Name}, fullPathParts...)

		// 获取城市信息
		var cityRegion model.Region
		if err := global.DB.Where("code = ?", currentRegion.Pcode).First(&cityRegion).Error; err == nil {
			regionPath.City = utils.ConvertRegionToRegionInfo(&cityRegion)
			fullPathParts = append([]string{cityRegion.Name}, fullPathParts...)

			// 获取省份信息
			var provinceRegion model.Region
			if err := global.DB.Where("code = ?", cityRegion.Pcode).First(&provinceRegion).Error; err == nil {
				regionPath.Province = utils.ConvertRegionToRegionInfo(&provinceRegion)
				fullPathParts = append([]string{provinceRegion.Name}, fullPathParts...)
			}
		}

	case 2: // 市级
		// 获取城市信息
		regionPath.City = utils.ConvertRegionToRegionInfo(&currentRegion)
		fullPathParts = append([]string{currentRegion.Name}, fullPathParts...)

		// 获取省份信息
		var provinceRegion model.Region
		if err := global.DB.Where("code = ?", currentRegion.Pcode).First(&provinceRegion).Error; err == nil {
			regionPath.Province = utils.ConvertRegionToRegionInfo(&provinceRegion)
			fullPathParts = append([]string{provinceRegion.Name}, fullPathParts...)
		}

	case 1: // 省级
		// 获取省份信息
		regionPath.Province = utils.ConvertRegionToRegionInfo(&currentRegion)
		fullPathParts = append([]string{currentRegion.Name}, fullPathParts...)
	}

	regionPath.FullPath = strings.Join(fullPathParts, "/")

	return &pb.GetRegionPathResp{
		Code:    200,
		Message: "获取区域路径成功",
		Data:    regionPath,
	}, nil
}

// SearchRegion 搜索区域接口
// 根据关键词搜索区域信息，支持模糊匹配
func (s *MapServiceImpl) SearchRegion(ctx context.Context, req *pb.SearchRegionReq) (*pb.SearchRegionResp, error) {
	if req.Keyword == "" {
		return &pb.SearchRegionResp{
			Code:    400,
			Message: "搜索关键词不能为空",
		}, nil
	}

	var regions []model.Region
	query := global.DB.Model(&model.Region{})

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
		regionInfos = append(regionInfos, utils.ConvertRegionToRegionInfo(&region))
	}

	return &pb.SearchRegionResp{
		Code:    200,
		Message: "搜索区域成功",
		Data:    regionInfos,
	}, nil
}
