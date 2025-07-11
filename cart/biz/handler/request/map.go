package request

// GeoCodingReq 地理编码请求结构体
type GeoCodingReq struct {
	Address string `json:"address" form:"address" vd:"required"`
}

// ReverseGeoCodingReq 逆地理编码请求结构体
type ReverseGeoCodingReq struct {
	Lat float64 `json:"lat" form:"lat" vd:"required"`
	Lng float64 `json:"lng" form:"lng" vd:"required"`
}

// IPLocationReq IP定位请求结构体
type IPLocationReq struct {
	IP string `json:"ip" form:"ip" vd:"required"`
}

// DistanceCalculateReq 距离计算请求结构体
type DistanceCalculateReq struct {
	Lat       float64 `json:"lat" form:"lat" vd:"required"`
	Lng       float64 `json:"lng" form:"lng" vd:"required"`
	DestLat   float64 `json:"dest_lat" form:"dest_lat"`
	DestLng   float64 `json:"dest_lng" form:"dest_lng"`
	Distance  float64 `json:"distance" form:"distance"`
	OriginLat float64 `json:"origin_lat" form:"origin_lat"` // 新增
	OriginLng float64 `json:"origin_lng" form:"origin_lng"` // 新增
}

// CalculateDistanceReq 计算距离请求结构体
type CalculateDistanceReq struct {
	OriginLat float64 `json:"origin_lat" form:"origin_lat" vd:"required"`
	OriginLng float64 `json:"origin_lng" form:"origin_lng" vd:"required"`
	DestLat   float64 `json:"dest_lat" form:"dest_lat" vd:"required"`
	DestLng   float64 `json:"dest_lng" form:"dest_lng" vd:"required"`
}

// GetProvincesReq 获取省份列表请求
type GetProvincesReq struct {
	ProvinceCode int64 `json:"province_code" form:"province_code" vd:"required"`
}

// GetCitiesReq 获取城市列表请求
type GetCitiesReq struct {
	ProvinceCode int64 `json:"province_code" form:"province_code" vd:"required"` // 新增
}

// GetDistrictsReq 获取区县列表请求
type GetDistrictsReq struct {
	CityCode int64 `json:"city_code" form:"city_code" vd:"required"` // 新增
}

// GetRegionPathReq 获取区域路径请求
type GetRegionPathReq struct {
	RegionCode int64 `json:"region_code" form:"region_code" vd:"required"`
}

// SearchRegionReq 搜索区域请求
type SearchRegionReq struct {
	Keyword    string `json:"keyword" form:"keyword" vd:"required"`
	ParentCode *int64 `json:"parent_code" form:"parent_code"`
}
