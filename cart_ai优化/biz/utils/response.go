package utils

import (
	"github.com/cloudwego/hertz/pkg/app"
)

// 错误码定义
const (
	// 成功
	CodeSuccess = 200

	// 客户端错误 4xx
	CodeBadRequest       = 400 // 请求参数错误
	CodeUnauthorized     = 401 // 未授权
	CodeForbidden        = 403 // 禁止访问
	CodeNotFound         = 404 // 资源不存在
	CodeMethodNotAllowed = 405 // 方法不允许
	CodeTooManyRequests  = 429 // 请求过于频繁

	// 业务错误 6xx
	CodeValidationFailed    = 600 // 参数验证失败
	CodeSMSCodeExpired      = 601 // 验证码过期
	CodeSMSCodeInvalid      = 602 // 验证码错误
	CodeUserNotExist        = 603 // 用户不存在
	CodeUserAlreadyExist    = 604 // 用户已存在
	CodeOrderNotFound       = 605 // 订单不存在
	CodeOrderStatusError    = 606 // 订单状态错误
	CodeInsufficientBalance = 607 // 余额不足
	CodeServiceUnavailable  = 608 // 服务不可用

	// 服务器错误 5xx
	CodeInternalError   = 500 // 服务器内部错误
	CodeDatabaseError   = 501 // 数据库错误
	CodeRedisError      = 502 // Redis错误
	CodeThirdPartyError = 503 // 第三方服务错误
)

// 错误信息映射
var ErrorMessages = map[int]string{
	CodeSuccess: "成功",

	// 客户端错误
	CodeBadRequest:       "请求参数错误",
	CodeUnauthorized:     "未授权，请先登录",
	CodeForbidden:        "禁止访问",
	CodeNotFound:         "资源不存在",
	CodeMethodNotAllowed: "请求方法不允许",
	CodeTooManyRequests:  "请求过于频繁，请稍后再试",

	// 业务错误
	CodeValidationFailed:    "参数验证失败",
	CodeSMSCodeExpired:      "验证码已过期，请重新获取",
	CodeSMSCodeInvalid:      "验证码错误",
	CodeUserNotExist:        "用户不存在",
	CodeUserAlreadyExist:    "用户已存在",
	CodeOrderNotFound:       "订单不存在",
	CodeOrderStatusError:    "订单状态错误",
	CodeInsufficientBalance: "余额不足",
	CodeServiceUnavailable:  "服务暂不可用",

	// 服务器错误
	CodeInternalError:   "服务器内部错误",
	CodeDatabaseError:   "数据库服务异常",
	CodeRedisError:      "缓存服务异常",
	CodeThirdPartyError: "第三方服务异常",
}

// APIResponse 统一响应结构
type APIResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// APIError 自定义错误类型
type APIError struct {
	Code    int
	Message string
	Details map[string]interface{}
	Err     error
}

func (e APIError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// NewAPIError 创建API错误
func NewAPIError(code int, message string, err error) *APIError {
	if message == "" {
		message = ErrorMessages[code]
	}
	return &APIError{
		Code:    code,
		Message: message,
		Err:     err,
		Details: make(map[string]interface{}),
	}
}

// WithDetails 添加错误详情
func (e *APIError) WithDetails(details map[string]interface{}) *APIError {
	for k, v := range details {
		e.Details[k] = v
	}
	return e
}

// Success 成功响应
func Success(c *app.RequestContext, data interface{}) {
	requestID, _ := c.Get("request_id")
	response := APIResponse{
		Code:      CodeSuccess,
		Message:   ErrorMessages[CodeSuccess],
		Data:      data,
		RequestID: getString(requestID),
		Timestamp: GetCurrentTimestamp(),
	}

	c.JSON(200, response)
}

// Error 错误响应
func ErrorResponse(c *app.RequestContext, apiErr *APIError) {
	requestID, _ := c.Get("request_id")

	// 记录错误日志
	LogError("api", "response", apiErr.Err, map[string]interface{}{
		"request_id": requestID,
		"code":       apiErr.Code,
		"message":    apiErr.Message,
		"details":    apiErr.Details,
		"path":       string(c.Path()),
		"method":     string(c.Method()),
		"client_ip":  c.ClientIP(),
	})

	response := APIResponse{
		Code:      apiErr.Code,
		Message:   apiErr.Message,
		RequestID: getString(requestID),
		Timestamp: GetCurrentTimestamp(),
	}

	// 生产环境隐藏敏感错误详情
	if GetEnv("APP_ENV", "development") == "production" && apiErr.Code >= 500 {
		response.Message = "服务器内部错误，请稍后重试"
	}

	// 根据错误码设置HTTP状态码
	httpStatus := getHTTPStatus(apiErr.Code)
	c.JSON(httpStatus, response)
}

// ErrorWithDetails 带详情的错误响应
func ErrorWithDetails(c *app.RequestContext, code int, message string, details map[string]interface{}, err error) {
	apiErr := NewAPIError(code, message, err).WithDetails(details)
	ErrorResponse(c, apiErr)
}

// getHTTPStatus 根据业务错误码获取HTTP状态码
func getHTTPStatus(code int) int {
	switch {
	case code >= 200 && code < 300:
		return 200
	case code >= 400 && code < 500:
		return code
	case code >= 600 && code < 700:
		return 400 // 业务错误返回400
	case code >= 500:
		return 500
	default:
		return 500
	}
}

// getString 安全地获取字符串值
func getString(v interface{}) string {
	if str, ok := v.(string); ok {
		return str
	}
	return ""
}

// 便捷方法
func BadRequest(c *app.RequestContext, message string, err error) {
	ErrorResponse(c, NewAPIError(CodeBadRequest, message, err))
}

func Unauthorized(c *app.RequestContext, message string) {
	ErrorResponse(c, NewAPIError(CodeUnauthorized, message, nil))
}

func Forbidden(c *app.RequestContext, message string) {
	ErrorResponse(c, NewAPIError(CodeForbidden, message, nil))
}

func NotFound(c *app.RequestContext, message string) {
	ErrorResponse(c, NewAPIError(CodeNotFound, message, nil))
}

func ValidationFailed(c *app.RequestContext, details map[string]interface{}) {
	apiErr := NewAPIError(CodeValidationFailed, "", nil).WithDetails(details)
	ErrorResponse(c, apiErr)
}

func InternalError(c *app.RequestContext, err error) {
	ErrorResponse(c, NewAPIError(CodeInternalError, "", err))
}

func DatabaseError(c *app.RequestContext, err error) {
	ErrorResponse(c, NewAPIError(CodeDatabaseError, "", err))
}

func RedisError(c *app.RequestContext, err error) {
	ErrorResponse(c, NewAPIError(CodeRedisError, "", err))
}

func ThirdPartyError(c *app.RequestContext, service string, err error) {
	details := map[string]interface{}{
		"service": service,
	}
	apiErr := NewAPIError(CodeThirdPartyError, "", err).WithDetails(details)
	ErrorResponse(c, apiErr)
}
