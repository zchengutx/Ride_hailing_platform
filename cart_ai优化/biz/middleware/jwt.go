package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	ID          int    `json:"id"`
	NickName    string `json:"nickname"`
	AuthorityId int    `json:"authority_id"`
	jwt.RegisteredClaims
}

func JWTAuth(secretKey string) app.HandlerFunc {
	if len(secretKey) == 0 {
		panic("secretKey cannot be empty")
	}

	return func(ctx context.Context, c *app.RequestContext) {
		// 获取token，支持多种方式
		token := string(c.Request.Header.Peek("x-token"))
		if token == "" {
			token = string(c.Request.Header.Peek("Authorization"))
			if token != "" && len(token) > 7 && token[:7] == "Bearer " {
				token = token[7:]
			}
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code": 401,
				"msg":  "请登录",
				"data": nil,
			})
			c.Abort()
			return
		}

		j := NewJWT(secretKey)
		claims, err := j.ParseToken(token)
		if err != nil {
			var errorMsg string
			switch {
			case errors.Is(err, TokenExpired):
				errorMsg = "授权已过期，请重新登录"
			case errors.Is(err, TokenNotValidYet):
				errorMsg = "Token尚未生效"
			case errors.Is(err, TokenMalformed):
				errorMsg = "Token格式错误"
			default:
				errorMsg = "Token验证失败"
			}

			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code": 401,
				"msg":  errorMsg,
				"data": nil,
			})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Next(ctx)
	}
}

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

func NewJWT(secretKey string) *JWT {
	if len(secretKey) == 0 {
		panic("secretKey cannot be empty")
	}
	return &JWT{
		SigningKey: []byte(secretKey),
	}
}

// CreateToken 创建token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	// 设置token过期时间
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "cart-service",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, TokenMalformed
		}
		return j.SigningKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, TokenExpired
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, TokenNotValidYet
		} else if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, TokenMalformed
		}
		return nil, TokenInvalid
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, TokenInvalid
}
