package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client/callopt"
	"lxh_cx/biz/request"
	"lxh_cx/config"
	"lxh_cx/kitex_gen/lxh_cx/driver"
	"net/http"
	"time"
)

func DriverRegister(ctx context.Context, c *app.RequestContext) {
	var req request.DriverRegisterReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"code":    http.StatusBadRequest,
			"message": "Failed to obtain parameters",
			"data":    err.Error(),
		})
		return
	}

	resp, _ := config.DriverCli.DriverRegister(ctx, &driver.DriverRegisterReq{
		IdCardFileId:         string(req.IdCardFileId),
		DriverLicenseFileId:  string(req.DriverLicenseFileId),
		DrivingLicenseFileId: string(req.DrivingLicenseFileId),
		AvatorFileId:         string(req.AvatorFileId),
		DriverId:             req.DriverId,
	}, callopt.WithConnectTimeout(time.Second*3))

	if resp.BaseResp.Code != 0 {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"code":    resp.BaseResp.Code,
			"message": resp.BaseResp.Msg,
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"code": resp.BaseResp.Code,
		"msg":  resp.BaseResp.Msg,
	})
}

func DriverInfoList(ctx context.Context, c *app.RequestContext) {

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"code":    consts.StatusUnauthorized,
			"message": "无法获取用户信息",
		})
		return
	}

	id, ok := userID.(int64)
	if !ok {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"code":    consts.StatusUnauthorized,
			"message": "用户ID类型错误",
		})
		return
	}

	list, _ := config.DriverCli.DriverInfoList(ctx, &driver.DriverInfoListReq{DriverId: int32(id)})
	if list.BaseResp.Code != 200 {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"code":    list.BaseResp.Code,
			"message": list.BaseResp.Msg,
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"code": list.BaseResp.Code,
		"msg":  list.BaseResp.Msg,
		"data": map[string]interface{}{
			"info": list.DriverInfo,
		},
	})

}

func DriverAdd(ctx context.Context, c *app.RequestContext) {

	var req request.DriverAddReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"code":    http.StatusBadRequest,
			"message": "Failed to obtain parameters",
			"data":    err.Error(),
		})
		return
	}

	resp, _ := config.DriverCli.DriverAdd(ctx, &driver.DriverAddReq{
		IdCardFileId:         string(req.IdCardFileId),
		DriverLicenseFileId:  string(req.DriverLicenseFileId),
		DrivingLicenseFileId: string(req.DrivingLicenseFileId),
		AvatorFileId:         string(req.AvatorFileId),
		DriverId:             req.DriverId,
		PassengerId:          req.PassengerId,
	}, callopt.WithConnectTimeout(time.Second*3))

	if resp.BaseResp.Code != 0 {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"code":    resp.BaseResp.Code,
			"message": resp.BaseResp.Msg,
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"code": resp.BaseResp.Code,
		"msg":  resp.BaseResp.Msg,
	})
}

func DriverOverOrder(ctx context.Context, c *app.RequestContext) {

	var req request.DriverOverOrderReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"code":    http.StatusBadRequest,
			"message": "Failed to obtain parameters",
			"data":    err.Error(),
		})
		return
	}

	resp, _ := config.DriverCli.DriverOverOrder(ctx, &driver.DriverOverOrderReq{
		OrderId: req.OrderId,
	}, callopt.WithConnectTimeout(time.Second*3))

	if resp.BaseResp.Code != 0 {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"code":    resp.BaseResp.Code,
			"message": resp.BaseResp.Msg,
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"code": resp.BaseResp.Code,
		"msg":  resp.BaseResp.Msg,
	})

}

func DriverCancelOrder(ctx context.Context, c *app.RequestContext) {
	var req request.DriverCancelOrderReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"code":    http.StatusBadRequest,
			"message": "Failed to obtain parameters",
			"data":    err.Error(),
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"code":    consts.StatusUnauthorized,
			"message": "无法获取用户信息",
		})
		return
	}

	id, ok := userID.(int64)
	if !ok {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"code":    consts.StatusUnauthorized,
			"message": "用户ID类型错误",
		})
		return
	}
	resp, _ := config.DriverCli.DriverCancelOrder(ctx, &driver.DriverCancelOrderReq{
		DriverId:      int32(id),
		ConfirmPerson: req.ConfirmPerson,
		ConfirmReason: string(req.ConfirmReason),
		ConfirmRemark: req.ConfirmRemark,
		OrderId:       req.OrderId,
	}, callopt.WithConnectTimeout(time.Second*3))
	if resp.BaseResp.Code != 0 {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"code":    resp.BaseResp.Code,
			"message": resp.BaseResp.Msg,
		})
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"code": resp.BaseResp.Code,
		"msg":  resp.BaseResp.Msg,
	})
}
