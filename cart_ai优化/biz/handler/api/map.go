// api 包下的 map.go 实现地图相关接口，如地理编码、逆地理编码、IP定位、距离计算等
package api

import (
	"cart/biz/handler/request"
	"cart/biz/utils"
	pb "cart/kitex_gen/cart/mapservice"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// MapClient 是与地图服务交互的RPC客户端
var (
	MapClient = utils.GetDefaultMapClient()
)

// GeoCoding 地理编码接口，根据地址获取坐标
func GeoCoding(ctx context.Context, c *app.RequestContext) {
	var req request.GeoCodingReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	result, err := MapClient.GeoCoding(ctx, &pb.GeoCodingReq{
		Address: req.Address,
	})
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "服务调用失败",
			"data": err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code": result.Code,
		"msg":  result.Message,
		"data": result.Data,
	})
}

// ReverseGeoCoding 逆地理编码接口，根据坐标获取地址
func ReverseGeoCoding(ctx context.Context, c *app.RequestContext) {
	var req request.ReverseGeoCodingReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	result, err := MapClient.ReverseGeoCoding(ctx, &pb.ReverseGeoCodingReq{
		Lat: req.Lat,
		Lng: req.Lng,
	})
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "服务调用失败",
			"data": err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code": result.Code,
		"msg":  result.Message,
		"data": result.Data,
	})
}

// IPLocation IP定位接口，根据IP获取位置信息
func IPLocation(ctx context.Context, c *app.RequestContext) {
	var req request.IPLocationReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	// 如果IP为空，尝试获取客户端IP
	if req.IP == "" {
		clientIP := c.ClientIP()
		if clientIP != "" {
			req.IP = clientIP
		}
	}

	result, err := MapClient.IPLocation(ctx, &pb.IPLocationReq{
		Ip: req.IP,
	})
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "服务调用失败",
			"data": err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code": result.Code,
		"msg":  result.Message,
		"data": result.Data,
	})
}

// DistanceCalculate 距离计算接口，计算两点间距离
func DistanceCalculate(ctx context.Context, c *app.RequestContext) {
	var req request.DistanceCalculateReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	result, err := MapClient.DistanceCalculate(ctx, &pb.DistanceCalculateReq{
		OriginLat: req.OriginLat,
		OriginLng: req.OriginLng,
		DestLat:   req.DestLat,
		DestLng:   req.DestLng,
	})
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "服务调用失败",
			"data": err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code": result.Code,
		"msg":  result.Message,
		"data": result.Data,
	})
}

// GetCurrentLocation 获取当前位置接口，基于IP获取当前位置信息
func GetCurrentLocation(ctx context.Context, c *app.RequestContext) {
	clientIP := c.ClientIP()
	if clientIP == "" {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "无法获取客户端IP",
			"data": nil,
		})
		return
	}

	result, err := MapClient.GetCurrentLocation(ctx, &pb.GetCurrentLocationReq{
		ClientIp: clientIP,
	})
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "服务调用失败",
			"data": err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code": result.Code,
		"msg":  result.Message,
		"data": result.Data,
	})
}

// GetProvinces 获取省份列表接口，获取所有省份信息
func GetProvinces(ctx context.Context, c *app.RequestContext) {
	result, err := MapClient.GetProvinces(ctx, &pb.GetProvincesReq{})
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "服务调用失败",
			"data": err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code": result.Code,
		"msg":  result.Message,
		"data": result.Data,
	})
}

// GetCities 获取城市列表接口，获取指定省份下的城市信息
func GetCities(ctx context.Context, c *app.RequestContext) {
	var req request.GetCitiesReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	result, err := MapClient.GetCities(ctx, &pb.GetCitiesReq{
		ProvinceCode: req.ProvinceCode,
	})
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "服务调用失败",
			"data": err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code": result.Code,
		"msg":  result.Message,
		"data": result.Data,
	})
}

// GetDistricts 获取区县列表接口，获取指定城市下的区县信息
func GetDistricts(ctx context.Context, c *app.RequestContext) {
	var req request.GetDistrictsReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	result, err := MapClient.GetDistricts(ctx, &pb.GetDistrictsReq{
		CityCode: req.CityCode,
	})
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "服务调用失败",
			"data": err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code": result.Code,
		"msg":  result.Message,
		"data": result.Data,
	})
}

// GetRegionPath 获取区域路径接口，获取指定区域的完整路径
func GetRegionPath(ctx context.Context, c *app.RequestContext) {
	var req request.GetRegionPathReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	result, err := MapClient.GetRegionPath(ctx, &pb.GetRegionPathReq{
		RegionCode: req.RegionCode,
	})
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "服务调用失败",
			"data": err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code": result.Code,
		"msg":  result.Message,
		"data": result.Data,
	})
}

// SearchRegion 搜索区域接口，根据关键词搜索区域信息
func SearchRegion(ctx context.Context, c *app.RequestContext) {
	var req request.SearchRegionReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	result, err := MapClient.SearchRegion(ctx, &pb.SearchRegionReq{
		Keyword:    req.Keyword,
		ParentCode: req.ParentCode,
	})
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "服务调用失败",
			"data": err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code": result.Code,
		"msg":  result.Message,
		"data": result.Data,
	})
}
