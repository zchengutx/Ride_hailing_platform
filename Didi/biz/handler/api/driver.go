package api

import (
	"Didi/biz/handler/request"
	"Didi/kitex_gen/Didi/driver"
	pb "Didi/kitex_gen/Didi/driver/driverserver"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
)

func CallACar(ctx context.Context, c *app.RequestContext) {
	var req request.CallACarReq
	if err := c.Bind(&req); err != nil {
		c.JSON(200, map[string]interface{}{
			"code":    400,
			"message": "查询失败",
			"data":    err.Error(),
		})
		return
	}
	newClient, _ := pb.NewClient("driver", client.WithHostPorts("127.0.0.1:50052"))
	car, _ := newClient.CallACar(ctx, &driver.CallACarReq{
		UserId:               int64(c.GetInt("userId")),
		Address:              req.Address,
		DrivingLicenseNumber: req.DrivingLicenseNumber,
		QuasiDrivingType:     req.QuasiDrivingType,
		DrivingAge:           int64(req.DrivingAge),
		CarNum:               req.CarNum,
		CarType:              req.CarType,
		VehicleMileage:       req.VehicleMileage,
		ServingTheCity:       req.ServingTheCity,
	})
	c.JSON(200, car)
}
