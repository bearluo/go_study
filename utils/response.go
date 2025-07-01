package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ValidationError 自定义验证错误结构
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`              // 业务错误码，0表示成功
	Message string      `json:"message"`           // 响应消息
	Data    interface{} `json:"data,omitempty"`    // 响应数据
	Error   string      `json:"error,omitempty"`   // 错误信息
	Details interface{} `json:"details,omitempty"` // 详细错误信息
}

// 业务错误码定义
const (
	CodeSuccess         = 0    // 成功
	CodeParamError      = 1001 // 参数错误
	CodeValidationError = 1002 // 验证错误
	CodeUserNotFound    = 2001 // 用户不存在
	CodeUserExists      = 2002 // 用户已存在
	CodePasswordError   = 2003 // 密码错误
	CodeTokenError      = 2004 // 令牌错误
	CodeTokenExpired    = 2005 // 令牌过期
	CodeUnauthorized    = 3001 // 未授权
	CodeForbidden       = 3002 // 禁止访问
	CodeNotFound        = 4001 // 资源不存在
	CodeSystemError     = 5001 // 系统错误
)

// Success 成功响应
func Success(c echo.Context, data interface{}, message string) error {
	if message == "" {
		message = "操作成功"
	}

	return c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// Error 业务错误响应
func Error(c echo.Context, code int, message string) error {
	return c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Error:   message,
	})
}

// ValidationErrorResponse 参数验证错误响应
func ValidationErrorResponse(c echo.Context, details interface{}) error {
	return c.JSON(http.StatusBadRequest, Response{
		Code:    CodeValidationError,
		Message: "参数验证失败",
		Error:   "参数验证失败",
		Details: details,
	})
}

// ValidationErrors 验证错误数组响应
func ValidationErrors(c echo.Context, errors []ValidationError) error {
	return c.JSON(http.StatusBadRequest, Response{
		Code:    CodeValidationError,
		Message: "参数验证失败",
		Error:   "参数验证失败",
		Details: errors,
	})
}

// ParamError 参数错误响应
func ParamError(c echo.Context, message string) error {
	if message == "" {
		message = "参数错误"
	}

	return c.JSON(http.StatusBadRequest, Response{
		Code:    CodeParamError,
		Message: message,
		Error:   message,
	})
}

// SystemError 系统错误响应
func SystemError(c echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, Response{
		Code:    CodeSystemError,
		Message: "系统内部错误",
		Error:   err.Error(),
	})
}

// NotFound 资源不存在响应
func NotFound(c echo.Context, message string) error {
	if message == "" {
		message = "资源不存在"
	}

	return c.JSON(http.StatusNotFound, Response{
		Code:    CodeNotFound,
		Message: message,
		Error:   message,
	})
}

// Unauthorized 未授权响应
func Unauthorized(c echo.Context, message string) error {
	if message == "" {
		message = "未授权访问"
	}

	return c.JSON(http.StatusUnauthorized, Response{
		Code:    CodeUnauthorized,
		Message: message,
		Error:   message,
	})
}

// Forbidden 禁止访问响应
func Forbidden(c echo.Context, message string) error {
	if message == "" {
		message = "禁止访问"
	}

	return c.JSON(http.StatusForbidden, Response{
		Code:    CodeForbidden,
		Message: message,
		Error:   message,
	})
}

// UserNotFound 用户不存在响应
func UserNotFound(c echo.Context) error {
	return Error(c, CodeUserNotFound, "用户不存在")
}

// UserExists 用户已存在响应
func UserExists(c echo.Context) error {
	return Error(c, CodeUserExists, "用户已存在")
}

// PasswordError 密码错误响应
func PasswordError(c echo.Context) error {
	return Error(c, CodePasswordError, "密码错误")
}

// TokenError 令牌错误响应
func TokenError(c echo.Context) error {
	return Error(c, CodeTokenError, "令牌错误")
}

// TokenExpired 令牌过期响应
func TokenExpired(c echo.Context) error {
	return Error(c, CodeTokenExpired, "令牌过期")
}
