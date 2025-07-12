package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"kitex_main/gateway/biz/handler/response"
	"kitex_main/kitex_gen/Ride_hailing_platform/ride"
	"kitex_main/kitex_gen/Ride_hailing_platform/ride/trips"
	"kitex_main/pkg"
	"net/http"
)

func CreateOrder(ctx context.Context, c *app.RequestContext) {
	cil, err := trips.NewClient("Ride_hailing_platform.ride", client.WithHostPorts("127.0.0.1:50053"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{
			Code:    400,
			Message: "init client err",
			Data:    err.Error(),
		})
		return
	}

	addr, end := pkg.CoordinateConversion(Coordinates)

	var userId int

	if value, exists := c.Get("userID"); exists {
		userId = value.(int)
	} else {
		c.JSON(http.StatusInternalServerError, response.Response{
			Code:    500,
			Message: "get user id error",
		})
		return
	}

	order, _ := cil.CreatedOrder(ctx, &ride.CreatedOrderReq{
		Amount:      int64(Price),
		PassengerId: int64(userId),
		StartAddr:   addr,
		EndEnd:      end,
	})

	if order.Code != 200 {
		c.JSON(http.StatusInternalServerError, response.Response{
			Code:    500,
			Message: "order error",
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "order success",
	})
}
