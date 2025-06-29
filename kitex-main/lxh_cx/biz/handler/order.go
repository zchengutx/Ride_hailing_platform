package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	request2 "lxh_cx/biz/request"
	"lxh_cx/config"
	"lxh_cx/kitex_gen/lxh_cx/order"
	"net/http"
	"net/url"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/gin-gonic/gin"
)

func PathPlanning(ctx context.Context, c *app.RequestContext) {
	var req request2.PathPlanningReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{
			"code":    http.StatusBadRequest,
			"message": "Failed to obtain parameters",
			"data":    err.Error(),
		})
		return
	}
	// 此处填写您在控制台-应用管理-创建应用后获取的AK
	ak := "ihBwl08smb1nlXb8Ds1h1xGVvruZkx6w"

	// 服务地址
	host := "https://api.map.baidu.com"

	// 接口地址
	uri := "/direction/v2/driving"

	// 设置请求参数
	params := url.Values{
		"origin":      []string{req.Origin},
		"destination": []string{req.Destination},
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

	fmt.Println(string(body))

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": map[string]interface{}{
			"data": string(body),
		},
	})
	add, _ := config.OrderCli.RouteAdd(ctx, &order.RouteAddReq{
		string(body),
	}, callopt.WithConnectTimeout(time.Second*3))
	if add.BaseResp.Code != 200 {
		c.JSON(http.StatusOK, gin.H{
			"code": add.BaseResp.Code,
			"msg":  add.BaseResp.Msg,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": add.BaseResp.Code,
		"msg":  add.BaseResp.Msg,
	})
}
