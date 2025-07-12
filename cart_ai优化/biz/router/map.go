package router

import (
	"cart/biz/dal/global"
	"cart/biz/handler/api"
	"cart/biz/middleware"
	"log"

	"github.com/cloudwego/hertz/pkg/route"
)

// MapModel 地图模块路由
func MapModel(r *route.RouterGroup) {
	log.Println("百度地图服务")
	mapModel := r.Group("/map")

	// 公开的地图服务 - 不需要认证
	{
		// 地理编码 - 根据地址获取坐标
		mapModel.POST("/geocoding", api.GeoCoding)
		mapModel.GET("/geocoding", api.GeoCoding)

		// 逆地理编码 - 根据坐标获取地址
		mapModel.POST("/reverse-geocoding", api.ReverseGeoCoding)
		mapModel.GET("/reverse-geocoding", api.ReverseGeoCoding)

		// IP定位 - 根据IP获取位置信息
		mapModel.POST("/ip-location", api.IPLocation)
		mapModel.GET("/ip-location", api.IPLocation)

		// 距离计算 - 计算两点间距离
		mapModel.POST("/distance", api.DistanceCalculate)
		mapModel.GET("/distance", api.DistanceCalculate)

		// Region 行政区划相关接口
		// 获取省份列表
		mapModel.GET("/provinces", api.GetProvinces)

		// 获取城市列表
		mapModel.GET("/cities", api.GetCities)
		mapModel.POST("/cities", api.GetCities)

		// 获取区县列表
		mapModel.GET("/districts", api.GetDistricts)
		mapModel.POST("/districts", api.GetDistricts)

		// 获取区域完整路径
		mapModel.GET("/region-path", api.GetRegionPath)
		mapModel.POST("/region-path", api.GetRegionPath)

		// 搜索区域
		mapModel.GET("/search-region", api.SearchRegion)
		mapModel.POST("/search-region", api.SearchRegion)
	}

	// 需要认证的个人位置服务
	mapAuth := mapModel.Group("")
	mapAuth.Use(middleware.JWTAuth(global.JWT_SELECT_KEY))
	{
		// 获取当前位置（基于IP）- 个人隐私信息，需要认证
		mapAuth.GET("/current-location", api.GetCurrentLocation)
	}
}
