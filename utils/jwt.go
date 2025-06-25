package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v4"
)

// JWTConfig JWT配置结构
type JWTConfig struct {
	SecretKey       string        // JWT密钥
	ExpiresIn       time.Duration // 过期时间
	RefreshExpiresIn time.Duration // 刷新token过期时间
	Issuer          string        // 签发者
}

// Claims JWT声明结构
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	UserType string `json:"user_type"` // passenger, driver, admin
	jwt.RegisteredClaims
}

// JWTManager JWT管理器
type JWTManager struct {
	config *JWTConfig
}

// NewJWTManager 创建JWT管理器
func NewJWTManager(config *JWTConfig) *JWTManager {
	// 设置默认配置
	if config.SecretKey == "" {
		config.SecretKey = "ride_hailing_platform_secret_key_2024"
	}
	if config.ExpiresIn == 0 {
		config.ExpiresIn = 24 * time.Hour // 默认24小时
	}
	if config.RefreshExpiresIn == 0 {
		config.RefreshExpiresIn = 7 * 24 * time.Hour // 默认7天
	}
	if config.Issuer == "" {
		config.Issuer = "ride_hailing_platform"
	}
	
	return &JWTManager{config: config}
}

// GenerateToken 生成访问token
func (jm *JWTManager) GenerateToken(userID int64, username, userType string) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:   userID,
		Username: username,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(jm.config.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    jm.config.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.config.SecretKey))
}

// GenerateRefreshToken 生成刷新token
func (jm *JWTManager) GenerateRefreshToken(userID int64, username, userType string) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:   userID,
		Username: username,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(jm.config.RefreshExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    jm.config.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.config.SecretKey))
}

// ValidateToken 验证token
func (jm *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jm.config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ExtractTokenFromHeader 从请求头中提取token
func (jm *JWTManager) ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}

// JWTMiddleware JWT中间件
func (jm *JWTManager) JWTMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		authHeader := string(c.GetHeader("Authorization"))
		
		// 提取token
		tokenString, err := jm.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code": 401,
				"msg":  "未提供有效的认证token",
				"data": nil,
			})
			c.Abort()
			return
		}

		// 验证token
		claims, err := jm.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code": 401,
				"msg":  "token验证失败: " + err.Error(),
				"data": nil,
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_type", claims.UserType)
		c.Set("claims", claims)

		c.Next(ctx)
	}
}

// OptionalJWTMiddleware 可选的JWT中间件（不强制要求token）
func (jm *JWTManager) OptionalJWTMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		authHeader := string(c.GetHeader("Authorization"))
		
		if authHeader != "" {
			// 提取token
			tokenString, err := jm.ExtractTokenFromHeader(authHeader)
			if err == nil {
				// 验证token
				claims, err := jm.ValidateToken(tokenString)
				if err == nil {
					// 将用户信息存储到上下文中
					c.Set("user_id", claims.UserID)
					c.Set("username", claims.Username)
					c.Set("user_type", claims.UserType)
					c.Set("claims", claims)
				}
			}
		}

		c.Next(ctx)
	}
}

// RoleMiddleware 角色验证中间件
func (jm *JWTManager) RoleMiddleware(allowedRoles ...string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code": 401,
				"msg":  "用户未认证",
				"data": nil,
			})
			c.Abort()
			return
		}

		userTypeStr := userType.(string)
		allowed := false
		for _, role := range allowedRoles {
			if userTypeStr == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, map[string]interface{}{
				"code": 403,
				"msg":  "权限不足",
				"data": nil,
			})
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}

// GetUserFromContext 从上下文中获取用户信息
func GetUserFromContext(c *app.RequestContext) (*Claims, error) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, errors.New("用户未认证")
	}

	userClaims, ok := claims.(*Claims)
	if !ok {
		return nil, errors.New("用户信息格式错误")
	}

	return userClaims, nil
}

// GetUserIDFromContext 从上下文中获取用户ID
func GetUserIDFromContext(c *app.RequestContext) (int64, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("用户未认证")
	}

	id, ok := userID.(int64)
	if !ok {
		return 0, errors.New("用户ID格式错误")
	}

	return id, nil
}

// GetUserTypeFromContext 从上下文中获取用户类型
func GetUserTypeFromContext(c *app.RequestContext) (string, error) {
	userType, exists := c.Get("user_type")
	if !exists {
		return "", errors.New("用户未认证")
	}

	typeStr, ok := userType.(string)
	if !ok {
		return "", errors.New("用户类型格式错误")
	}

	return typeStr, nil
}

// TokenPair token对结构
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// GenerateTokenPair 生成token对（访问token和刷新token）
func (jm *JWTManager) GenerateTokenPair(userID int64, username, userType string) (*TokenPair, error) {
	accessToken, err := jm.GenerateToken(userID, username, userType)
	if err != nil {
		return nil, fmt.Errorf("生成访问token失败: %w", err)
	}

	refreshToken, err := jm.GenerateRefreshToken(userID, username, userType)
	if err != nil {
		return nil, fmt.Errorf("生成刷新token失败: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(jm.config.ExpiresIn.Seconds()),
	}, nil
}

// RefreshAccessToken 刷新访问token
func (jm *JWTManager) RefreshAccessToken(refreshToken string) (*TokenPair, error) {
	claims, err := jm.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("刷新token验证失败: %w", err)
	}

	// 生成新的token对
	return jm.GenerateTokenPair(claims.UserID, claims.Username, claims.UserType)
}

// 默认JWT管理器实例
var defaultJWTManager *JWTManager

// InitDefaultJWTManager 初始化默认JWT管理器
func InitDefaultJWTManager(config *JWTConfig) {
	defaultJWTManager = NewJWTManager(config)
}

// GetDefaultJWTManager 获取默认JWT管理器
func GetDefaultJWTManager() *JWTManager {
	if defaultJWTManager == nil {
		// 使用默认配置初始化
		defaultJWTManager = NewJWTManager(&JWTConfig{})
	}
	return defaultJWTManager
}

// 便捷函数，使用默认JWT管理器
func GenerateToken(userID int64, username, userType string) (string, error) {
	return GetDefaultJWTManager().GenerateToken(userID, username, userType)
}

func GenerateTokenPair(userID int64, username, userType string) (*TokenPair, error) {
	return GetDefaultJWTManager().GenerateTokenPair(userID, username, userType)
}

func ValidateToken(tokenString string) (*Claims, error) {
	return GetDefaultJWTManager().ValidateToken(tokenString)
}

func JWTMiddleware() app.HandlerFunc {
	return GetDefaultJWTManager().JWTMiddleware()
}

func OptionalJWTMiddleware() app.HandlerFunc {
	return GetDefaultJWTManager().OptionalJWTMiddleware()
}

func RoleMiddleware(allowedRoles ...string) app.HandlerFunc {
	return GetDefaultJWTManager().RoleMiddleware(allowedRoles...)
}
