package request

// GeoCodingReq 地理编码请求
type GeoCodingReq struct {
	Address string `json:"address" form:"address" binding:"required"` // 地址信息
}

// ReverseGeoCodingReq 逆地理编码请求
type ReverseGeoCodingReq struct {
	Lat float64 `json:"lat" form:"lat" binding:"required"` // 纬度
	Lng float64 `json:"lng" form:"lng" binding:"required"` // 经度
}

// IPLocationReq IP定位请求
type IPLocationReq struct {
	IP string `json:"ip" form:"ip" binding:"required"` // IP地址
}

// NearbySearchReq 周边搜索请求
type NearbySearchReq struct {
	Lat      float64 `json:"lat" form:"lat" binding:"required"` // 纬度
	Lng      float64 `json:"lng" form:"lng" binding:"required"` // 经度
	Query    string  `json:"query" form:"query"`                // 搜索关键词
	Radius   int     `json:"radius" form:"radius"`              // 搜索半径，默认1000米
	PageSize int     `json:"page_size" form:"page_size"`        // 每页显示条数，默认10
	PageNum  int     `json:"page_num" form:"page_num"`          // 页码，默认0
}

// DistanceCalculateReq 距离计算请求
type DistanceCalculateReq struct {
	OriginLat float64 `json:"origin_lat" form:"origin_lat" binding:"required"` // 起点纬度
	OriginLng float64 `json:"origin_lng" form:"origin_lng" binding:"required"` // 起点经度
	DestLat   float64 `json:"dest_lat" form:"dest_lat" binding:"required"`     // 终点纬度
	DestLng   float64 `json:"dest_lng" form:"dest_lng" binding:"required"`     // 终点经度
}

// GetCitiesReq 获取城市列表请求
type GetCitiesReq struct {
	ProvinceCode int64 `json:"province_code" form:"province_code" binding:"required"` // 省份代码
}

// GetDistrictsReq 获取区县列表请求
type GetDistrictsReq struct {
	CityCode int64 `json:"city_code" form:"city_code" binding:"required"` // 城市代码
}

// GetRegionPathReq 获取区域路径请求
type GetRegionPathReq struct {
	RegionCode int64 `json:"region_code" form:"region_code" binding:"required"` // 区域代码
}

// SearchRegionReq 搜索区域请求
type SearchRegionReq struct {
	Keyword    string `json:"keyword" form:"keyword" binding:"required"` // 搜索关键词
	ParentCode *int64 `json:"parent_code" form:"parent_code"`            // 父级区域代码，可选
}
