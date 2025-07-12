package router

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"kitex_main/gateway/biz/handler/response"
	"kitex_main/pkg"
	"net/http"
	"strconv"
)

func Middleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		peek := string(ctx.Request.Header.Peek("Authorization"))
		if peek == "" {
			resp := response.Response{
				Code:    401,
				Message: "缺少授权信息",
				Data:    nil,
			}
			ctx.JSON(http.StatusUnauthorized, resp)
			ctx.Abort()
			return
		}

		fmt.Println("peek:", peek)

		//if len(peek) < 7 || peek[:7] != "Bearer " {
		//	resp := response.Response{
		//		Code:    402,
		//		Message: "Token 格式错误",
		//		Data:    nil,
		//	}
		//	ctx.JSON(http.StatusUnauthorized, resp)
		//	ctx.Abort()
		//	return
		//}
		//tokenStr := peek[7:]

		token, s := pkg.GetToken(peek)

		if token == nil || s != "" {
			resp := response.Response{
				Code:    403,
				Message: "登录过期或无效 Token",
				Data:    nil,
			}
			ctx.JSON(http.StatusUnauthorized, resp)
			ctx.Abort()
			return
		}
		userIDStr, ok := token["user"].(string)
		if !ok {
			resp := response.Response{
				Code:    403,
				Message: "无效的用户信息",
				Data:    nil,
			}
			ctx.JSON(http.StatusUnauthorized, resp)
			ctx.Abort()
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			resp := response.Response{
				Code:    403,
				Message: "无效的用户 ID",
				Data:    nil,
			}
			ctx.JSON(http.StatusUnauthorized, resp)
			ctx.Abort()
			return
		}

		ctx.Set("userID", userID)

		ctx.Next(c)
	}
}
