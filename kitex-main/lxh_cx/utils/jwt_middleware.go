package utils

import (
	"context"
	"time"

	"lxh_cx/biz/request"
	"lxh_cx/config"
	"lxh_cx/kitex_gen/lxh_cx/user"
	"lxh_cx/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/hertz-contrib/jwt"
)

var identityKey = "id"

// NewJWTMiddleware 创建 JWT 中间件
func NewJWTMiddleware() (*jwt.HertzJWTMiddleware, error) {
	return jwt.New(&jwt.HertzJWTMiddleware{
		Realm:       "hertz jwt",
		Key:         []byte("2211a"), // 使用与之前相同的密钥
		Timeout:     time.Hour * 24,  // token 有效期 24 小时
		MaxRefresh:  time.Hour * 24,  // 刷新 token 有效期
		IdentityKey: identityKey,

		// PayloadFunc 设置登录时为 token 添加自定义负载信息
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.LxhPassenger); ok {
				return jwt.MapClaims{
					identityKey: int64(v.Id), // 统一转换为 int64 存储
					"nickname":  v.NickName,
					"user_name": v.Name,
					"mobile":    v.Mobile,
				}
			}
			return jwt.MapClaims{}
		},

		// IdentityHandler 从 token 提取用户信息
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)

			// 安全地从 claims 中提取 ID
			idClaim, ok := claims[identityKey]
			if !ok {
				return nil
			}

			// JWT claims 中的数字会被解析为 float64，需要正确转换
			var id int64
			switch v := idClaim.(type) {
			case float64:
				id = int64(v)
			case int64:
				id = v
			case int:
				id = int64(v)
			default:
				return nil
			}

			// 将 int64 的 ID 存储到 context 中，供后续使用
			c.Set("user_id", id)

			// 安全地提取其他字段
			nickname, _ := claims["nickname"].(string)
			userName, _ := claims["user_name"].(string)
			mobile, _ := claims["mobile"].(string)

			return &model.LxhPassenger{
				Id:       int32(id), // 数据库模型仍使用 int32
				NickName: nickname,
				Name:     userName,
				Mobile:   mobile,
			}
		},

		// Authenticator 验证用户登录信息
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var loginVals request.LoginReq
			if err := c.BindAndValidate(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			// 调用现有的登录验证逻辑
			resp, err := config.UserCli.Login(ctx, &user.LoginReq{
				Mobile:   loginVals.Mobile,
				SendCode: loginVals.SendCode,
			}, callopt.WithConnectTimeout(time.Second*3))

			if err != nil || resp.BaseResp.Code != 200 {
				return nil, jwt.ErrFailedAuthentication
			}

			// 返回用户信息，注意 resp.Id 是 int64 类型
			return &model.LxhPassenger{
				Id:       int32(resp.Id), // RPC 返回的是 int64，转换为 int32 存储到模型中
				Mobile:   loginVals.Mobile,
				NickName: "User", // 可以根据实际业务设置昵称
				Name:     "User", // 可以根据实际业务设置用户名
			}, nil
		},

		// Authorizator 设置用户访问权限
		Authorizator: func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
			// 这里可以添加更复杂的权限验证逻辑
			// 目前允许所有已认证用户访问
			if _, ok := data.(*model.LxhPassenger); ok {
				return true
			}
			return false
		},

		// Unauthorized 授权失败后的响应
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(code, utils.H{
				"code":    code,
				"message": message,
			})
		},

		// LoginResponse 登录成功响应
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			c.JSON(consts.StatusOK, utils.H{
				"code":    200,
				"message": "login successfully",
				"data": map[string]interface{}{
					"token":  token,
					"expire": expire.Format(time.RFC3339),
				},
			})
		},

		// RefreshResponse 刷新 token 响应
		RefreshResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			c.JSON(consts.StatusOK, utils.H{
				"code":   200,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},

		// TokenLookup 指定从哪里获取 token
		// 支持 header、query、cookie 等多种方式
		TokenLookup: "header: Authorization, header: x-token",

		// TokenHeadName 指定 token 前缀
		TokenHeadName: "Bearer",

		// TimeFunc 获取当前时间的函数
		TimeFunc: time.Now,
	})
}
