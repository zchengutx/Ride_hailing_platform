namespace go cart.mapservice

// 地理编码请求
struct GeoCodingReq {
    1: string address
}

// 地理编码响应
struct GeoCodingResp {
    1: i16 code
    2: string message
    3: optional GeoCodingData data
}

struct GeoCodingData {
    1: string address
    2: double lng
    3: double lat
    4: i32 precise
    5: string level
}

// 逆地理编码请求
struct ReverseGeoCodingReq {
    1: double lat
    2: double lng
}

// 逆地理编码响应
struct ReverseGeoCodingResp {
    1: i16 code
    2: string message
    3: optional ReverseGeoCodingData data
}

struct ReverseGeoCodingData {
    1: double lng
    2: double lat
    3: string formatted_address
    4: string business
    5: string citycode
}

// IP定位请求
struct IPLocationReq {
    1: string ip
}

// IP定位响应
struct IPLocationResp {
    1: i16 code
    2: string message
    3: optional IPLocationData data
}

struct IPLocationData {
    1: string ip
    2: string address
    3: string address_detail
    4: optional Point point
}

struct Point {
    1: double x
    2: double y
}

// 距离计算请求
struct DistanceCalculateReq {
    1: double origin_lat
    2: double origin_lng
    3: double dest_lat
    4: double dest_lng
}

// 距离计算响应
struct DistanceCalculateResp {
    1: i16 code
    2: string message
    3: optional DistanceData data
}

struct DistanceData {
    1: LocationPoint origin
    2: LocationPoint destination
    3: double distance_meters
    4: double distance_km
}

struct LocationPoint {
    1: double lat
    2: double lng
}

// 获取当前位置请求
struct GetCurrentLocationReq {
    1: string client_ip
}

// 获取当前位置响应
struct GetCurrentLocationResp {
    1: i16 code
    2: string message
    3: optional IPLocationData data
}

// 行政区划信息
struct RegionInfo {
    1: i64 code         // 行政区划代码
    2: string name      // 行政区划名称
    3: i64 pcode        // 上级区划代码
    4: string sname     // 地名简称
    5: i64 level        // 行政区划等级（1：省、直辖市；2：市州；3：区县）
    6: string mername   // 组合名称
    7: string pinyin    // 拼音
}

// 获取省份列表请求
struct GetProvincesReq {
    // 空请求，获取所有省份
}

// 获取省份列表响应
struct GetProvincesResp {
    1: i16 code
    2: string message
    3: optional list<RegionInfo> data
}

// 获取城市列表请求
struct GetCitiesReq {
    1: i64 province_code // 省份代码
}

// 获取城市列表响应
struct GetCitiesResp {
    1: i16 code
    2: string message
    3: optional list<RegionInfo> data
}

// 获取区县列表请求
struct GetDistrictsReq {
    1: i64 city_code // 城市代码
}

// 获取区县列表响应
struct GetDistrictsResp {
    1: i16 code
    2: string message
    3: optional list<RegionInfo> data
}

// 根据区域代码获取完整路径请求
struct GetRegionPathReq {
    1: i64 region_code // 区域代码
}

// 获取区域路径响应
struct GetRegionPathResp {
    1: i16 code
    2: string message
    3: optional RegionPath data
}

struct RegionPath {
    1: optional RegionInfo province  // 省份信息
    2: optional RegionInfo city      // 城市信息
    3: optional RegionInfo district  // 区县信息
    4: string full_path              // 完整路径，如"北京市/朝阳区"
}

// 搜索区域请求
struct SearchRegionReq {
    1: string keyword   // 搜索关键词
    2: optional i64 parent_code // 父级区域代码，限制搜索范围
}

// 搜索区域响应
struct SearchRegionResp {
    1: i16 code
    2: string message
    3: optional list<RegionInfo> data
}

// 地图服务接口
service MapService {
    // 地理编码 - 根据地址获取坐标
    GeoCodingResp GeoCoding(1: GeoCodingReq req)
    
    // 逆地理编码 - 根据坐标获取地址
    ReverseGeoCodingResp ReverseGeoCoding(1: ReverseGeoCodingReq req)
    
    // IP定位 - 根据IP获取位置信息
    IPLocationResp IPLocation(1: IPLocationReq req)
    
    // 距离计算 - 计算两点间距离
    DistanceCalculateResp DistanceCalculate(1: DistanceCalculateReq req)
    
    // 获取当前位置（基于IP）
    GetCurrentLocationResp GetCurrentLocation(1: GetCurrentLocationReq req)
    
    // 获取省份列表
    GetProvincesResp GetProvinces(1: GetProvincesReq req)
    
    // 获取城市列表
    GetCitiesResp GetCities(1: GetCitiesReq req)
    
    // 获取区县列表
    GetDistrictsResp GetDistricts(1: GetDistrictsReq req)
    
    // 获取区域完整路径
    GetRegionPathResp GetRegionPath(1: GetRegionPathReq req)
    
    // 搜索区域
    SearchRegionResp SearchRegion(1: SearchRegionReq req)
} 