package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/dgrijalva/jwt-go"
)

// CustomClaims 自定义JWT声明结构体
// 包含用户的基本信息和标准JWT声明
type CustomClaims struct {
	ID                 int    // 用户ID
	NickName           string // 用户昵称
	AuthorityId        int    // 权限ID，用于角色权限控制
	jwt.StandardClaims        // JWT标准声明（包含过期时间等）
}

// JWTAuth JWT认证中间件
// 用于验证HTTP请求头中的JWT token，确保用户已登录
// 参数secretKey: JWT签名密钥
// 返回值: Hertz框架的中间件处理函数
func JWTAuth(secretKey string) app.HandlerFunc {
	// 检查密钥是否为空，为空则抛出异常
	if len(secretKey) == 0 {
		panic("secretKey is not null")
	}

	return func(ctx context.Context, c *app.RequestContext) {
		// 从HTTP请求头中获取JWT token
		// 前端需要在请求头的"x-token"字段中携带登录时获得的token
		// 前端通常将token存储在cookie或localStorage中
		token := string(c.Request.Header.Peek("x-token"))

		// 检查token是否存在
		if token == "" {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"msg": "请登录",
			})
			c.Abort() // 终止请求处理
			return
		}

		// 创建JWT处理器实例
		j := NewJWT(secretKey)

		// 解析token，获取其中包含的用户信息
		claims, err := j.ParseToken(token)
		if err != nil {
			// 处理token过期的情况
			if err == TokenExpired {
				c.JSON(http.StatusUnauthorized, map[string]string{
					"msg": "授权已过期",
				})
				c.Abort()
				return
			}

			// 处理其他token验证失败的情况
			c.JSON(http.StatusUnauthorized, "未登陆")
			c.Abort()
			return
		}

		// 将解析出的用户信息存储到请求上下文中，供后续处理使用
		c.Set("claims", claims)    // 存储完整的声明信息
		c.Set("userId", claims.ID) // 存储用户ID，方便后续业务逻辑直接使用

		// 继续执行后续的处理链
		c.Next(ctx)
	}
}

// JWT JWT处理器结构体
// 用于创建、解析和刷新JWT token
type JWT struct {
	SigningKey []byte // JWT签名密钥，用于token的签名和验证
}

// JWT相关错误定义
var (
	TokenExpired     = errors.New("Token is expired")            // token已过期
	TokenNotValidYet = errors.New("Token not active yet")        // token尚未生效
	TokenMalformed   = errors.New("That's not even a token")     // token格式错误
	TokenInvalid     = errors.New("Couldn't handle this token:") // token无效
)

// NewJWT 创建JWT处理器实例
// 参数secretKey: JWT签名密钥，用于token的生成和验证
// 返回值: JWT处理器指针
func NewJWT(secretKey string) *JWT {
	// 检查密钥是否为空，为空则抛出异常
	if len(secretKey) == 0 {
		panic("secretKey is not null")
	}
	return &JWT{
		SigningKey: []byte(secretKey), // 将字符串密钥转换为字节数组
	}
}

// CreateToken 创建JWT token
// 使用HS256算法对用户声明进行签名，生成JWT token字符串
// 参数claims: 包含用户信息的自定义声明
// 返回值: (token字符串, 错误信息)
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	// 使用HS256签名算法创建token，并包含用户声明信息
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用签名密钥对token进行签名，返回最终的token字符串
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析JWT token
// 验证token的签名并解析出其中包含的用户声明信息
// 参数tokenString: 需要解析的JWT token字符串
// 返回值: (用户声明信息指针, 错误信息)
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	// 使用签名密钥解析token，并将声明信息解析为CustomClaims结构体
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})

	// 处理解析过程中的各种错误情况
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			// 检查具体的验证错误类型
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// token格式不正确
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// token已过期
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				// token尚未生效（nbf字段限制）
				return nil, TokenNotValidYet
			} else {
				// 其他验证错误
				return nil, TokenInvalid
			}
		}
	}

	// 检查token是否成功解析且有效
	if token != nil {
		// 尝试获取声明信息并验证token的有效性
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		// token为空，返回无效错误
		return nil, TokenInvalid
	}
}

// RefreshToken 刷新JWT token
// 为即将过期的token生成新的token，延长用户的登录时间
// 参数tokenString: 需要刷新的原始token字符串
// 返回值: (新的token字符串, 错误信息)
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	// 暂时禁用时间验证，允许解析已过期的token
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	// 解析原始token以获取用户声明信息
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}

	// 验证token并提取声明信息
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		// 恢复正常的时间函数
		jwt.TimeFunc = time.Now
		// 更新token的过期时间为当前时间+1小时
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		// 使用新的过期时间创建新token
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
