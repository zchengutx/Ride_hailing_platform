package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"io/ioutil"
	request2 "kitex_main/gateway/biz/handler/request"
	"kitex_main/gateway/biz/handler/response"
	home2 "kitex_main/kitex_gen/Ride_hailing_platform/home"
	"kitex_main/kitex_gen/Ride_hailing_platform/home/home"
	"kitex_main/kitex_gen/Ride_hailing_platform/ride"
	"kitex_main/kitex_gen/Ride_hailing_platform/ride/trips"
	"kitex_main/pkg"
	"net/http"
	"net/url"
)

type RouteResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Result  struct {
		Origin struct {
			Lng float64 `json:"lng"`
			Lat float64 `json:"lat"`
		} `json:"origin"`
		Destination struct {
			Lng float64 `json:"lng"`
			Lat float64 `json:"lat"`
		} `json:"destination"`
		Routes []struct {
			RouteMd5         string `json:"route_md5"`
			Distance         int    `json:"distance"`
			Duration         int    `json:"duration"`
			TrafficCondition int    `json:"traffic_condition"`
			Toll             int    `json:"toll"`
			RestrictionInfo  struct {
				Status int `json:"status"`
			} `json:"restriction_info"`
			Steps []struct {
				LegIndex         int    `json:"leg_index"`
				Distance         int    `json:"distance"`
				Duration         int    `json:"duration"`
				Direction        int    `json:"direction"`
				Turn             int    `json:"turn"`
				RoadType         int    `json:"road_type"`
				RoadTypes        string `json:"road_types"`
				Instruction      string `json:"instruction"`
				Path             string `json:"path"`
				TrafficCondition []struct {
					Status int `json:"status"`
					GeoCnt int `json:"geo_cnt"`
				} `json:"traffic_condition"`
				StartLocation struct {
					Lng string `json:"lng"`
					Lat string `json:"lat"`
				} `json:"start_location"`
				EndLocation struct {
					Lng string `json:"lng"`
					Lat string `json:"lat"`
				} `json:"end_location"`
			} `json:"steps"`
		} `json:"routes"`
	} `json:"result"`
}

var Coordinates []string

var Price float64

// 获取当前经纬度
func Driving(ctx context.Context, c *app.RequestContext) {

	// 此处填写您在控制台-应用管理-创建应用后获取的AK
	ak := "fbZPlNxNhl49M4bg6zHwLTP0DpJym7eS"

	// 服务地址
	host := "https://api.map.baidu.com"

	// 接口地址
	uri := "/location/ip"

	ip, _ := pkg.GetExternalIP()

	// 设置请求参数
	params := url.Values{
		"coor": []string{"bd09ll"},
		"ip":   []string{ip},
		"ak":   []string{ak},
	}

	fmt.Println(params.Encode())

	// 发起请求
	request, err := url.Parse(host + uri + "?" + params.Encode())
	if nil != err {
		fmt.Printf("host error: %v", err)
		return
	}

	resp, err1 := http.Get(request.String())
	fmt.Printf("url: %s\n", request.String())
	defer resp.Body.Close()
	if err1 != nil {
		fmt.Printf("request error: %v", err1)
		return
	}
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Printf("response error: %v", err2)
	}
	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "driving success",
		Data:    body,
	})
}

// 规划行程接口
func Directionlite(ctx context.Context, c *app.RequestContext) {

	// 此处填写您在控制台-应用管理-创建应用后获取的AK
	ak := "fbZPlNxNhl49M4bg6zHwLTP0DpJym7eS"

	// 服务地址
	host := "https://api.map.baidu.com"

	// 接口地址
	uri := "/directionlite/v1/driving"

	var req request2.Driving
	err := c.BindAndValidate(&req)
	if err != nil {
		resp := response.Response{
			Code:    400,
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// 设置请求参数
	params := url.Values{
		"origin":      []string{req.Origins},
		"destination": []string{req.Destinations},
		"ak":          []string{ak},
	}

	// 发起请求
	request, err := url.Parse(host + uri + "?" + params.Encode())
	if nil != err {
		fmt.Printf("host error: %v", err)
		return
	}

	resp, err1 := http.Get(request.String())
	fmt.Printf("url: %s\n", request.String())
	defer resp.Body.Close()
	if err1 != nil {
		fmt.Printf("request error: %v", err1)
		return
	}
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Printf("response error: %v", err2)
	}

	cil, err := trips.NewClient("Ride_hailing_platform.ride", client.WithHostPorts("127.0.0.1:50053"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{
			Code:    400,
			Message: "init ride error",
		})
		return
	}

	Trip, _ := cil.CreateTrips(ctx, &ride.CreateTripsReq{
		Data: string(body),
	})

	if Trip.Code != 200 {
		c.JSON(http.StatusInternalServerError, response.Response{
			Code:    400,
			Message: "trips error",
		})
		return
	}

	var Data RouteResponse

	json.Unmarshal(body, &Data)

	Price, err = pkg.ComputePrices(float64(Data.Result.Routes[0].Distance)/1000.0, req.CarType)

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{
			Code:    500,
			Message: "compute prices error",
		})
		return
	}

	Coordinates = append(Coordinates, req.Origins)
	Coordinates = append(Coordinates, req.Destinations)

	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "success",
	})

}

//保存历史记录

func CreateHistoricalSearch(ctx context.Context, c *app.RequestContext) {

	cil, err := home.NewClient("Ride_hailing_platform.home", client.WithHostPorts("127.0.0.1:50052"))
	if err != nil {
		resp := response.Response{
			Code:    400,
			Message: "init historical search error",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	var userId int

	if value, exists := c.Get("userID"); exists {
		userId = value.(int)
	} else {
		resp := response.Response{
			Code:    400,
			Message: "user not exists",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	var req request2.CreateHistoricalSearch

	if err = c.BindAndValidate(&req); err != nil {
		resp := response.Response{
			Code:    400,
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	search, _ := cil.CreateHistoricalSearch(ctx, &home2.CreateHistoricalSearchReq{
		UserId:     int64(userId),
		AddrName:   req.AddrName,
		IpAddrName: req.IdAddr,
	})
	if search.Code != 200 {
		resp := response.Response{
			Code:    400,
			Message: "historical search error",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp := response.Response{
		Code:    200,
		Message: "historical search success",
	}
	c.JSON(http.StatusOK, resp)

}

// 历史记录展示
func HistoricalSearchList(ctx context.Context, c *app.RequestContext) {
	cil, err := home.NewClient("Ride_hailing_platform.home", client.WithHostPorts("127.0.0.1:50052"))
	if err != nil {
		resp := response.Response{
			Code:    400,
			Message: "init historical search error",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	var userId int

	if value, exists := c.Get("userID"); exists {
		userId = value.(int)
	} else {
		resp := response.Response{
			Code:    400,
			Message: "user not exists",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	var req request2.HistoricalSearchList
	if err = c.BindAndValidate(&req); err != nil {
		resp := response.Response{
			Code:    400,
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	list, _ := cil.HistoricalSearchList(ctx, &home2.HistoricalSearchListReq{
		UserId:     int64(userId),
		IpAddrName: req.IdAddr,
	})
	if list.Code != 200 {
		resp := response.Response{
			Code:    400,
			Message: "historical search error",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := response.Response{
		Code:    200,
		Message: "historical search success",
		Data:    list.List,
	}
	c.JSON(http.StatusOK, resp)

}
