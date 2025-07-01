package handles

import (
	"go-study/services"
	"go-study/utils"
	"strings"

	"github.com/labstack/echo/v4"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register POST 用户注册
func (h *AuthHandler) Register(c echo.Context) error {
	var req services.RegisterRequest

	// 绑定请求数据
	if err := c.Bind(&req); err != nil {
		return utils.ParamError(c, "请求参数错误")
	}

	// 使用 Echo 验证器验证请求参数
	if err := c.Validate(&req); err != nil {
		// 获取详细的验证错误信息
		validationErrors := utils.GetValidationErrors(err)
		return utils.ValidationErrors(c, validationErrors)
	}

	// 执行注册
	response, err := h.authService.Register(&req)
	if err != nil {
		// 根据错误类型返回不同的响应
		if strings.Contains(err.Error(), "邮箱已存在") {
			return utils.Error(c, utils.CodeUserExists, "邮箱已存在")
		}
		if strings.Contains(err.Error(), "用户名已存在") {
			return utils.Error(c, utils.CodeUserExists, "用户名已存在")
		}
		return utils.SystemError(c, err)
	}

	return utils.Success(c, response, "注册成功")
}

// Login 用户登录
func (h *AuthHandler) Login(c echo.Context) error {
	var req services.LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.ParamError(c, "请求参数错误")
	}

	// 验证请求参数
	if err := c.Validate(&req); err != nil {
		// 获取详细的验证错误信息
		validationErrors := utils.GetValidationErrors(err)
		return utils.ValidationErrors(c, validationErrors)
	}

	// 执行登录
	response, err := h.authService.Login(&req)
	if err != nil {
		return utils.Unauthorized(c, err.Error())
	}

	return utils.Success(c, response, "登录成功")
}

// RefreshToken 刷新访问令牌
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req services.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return utils.ParamError(c, "请求参数错误")
	}

	// 验证请求参数
	if err := c.Validate(&req); err != nil {
		// 获取详细的验证错误信息
		validationErrors := utils.GetValidationErrors(err)
		return utils.ValidationErrors(c, validationErrors)
	}

	// 执行刷新令牌
	response, err := h.authService.RefreshToken(&req)
	if err != nil {
		return utils.Unauthorized(c, err.Error())
	}

	return utils.Success(c, response, "令牌刷新成功")
}

// Logout 用户登出
func (h *AuthHandler) Logout(c echo.Context) error {
	var req services.LogoutRequest
	if err := c.Bind(&req); err != nil {
		return utils.ParamError(c, "请求参数错误")
	}

	// 验证请求参数
	if err := c.Validate(&req); err != nil {
		// 获取详细的验证错误信息
		validationErrors := utils.GetValidationErrors(err)
		return utils.ValidationErrors(c, validationErrors)
	}

	// 执行登出
	if err := h.authService.Logout(&req); err != nil {
		return utils.SystemError(c, err)
	}

	return utils.Success(c, map[string]string{
		"message": "登出成功",
	}, "登出成功")
}

// LogoutAll 撤销用户的所有令牌
func (h *AuthHandler) LogoutAll(c echo.Context) error {
	// 从上下文中获取用户ID
	userID := c.Get("user_id").(uint)
	if userID == 0 {
		return utils.Unauthorized(c, "用户未认证")
	}

	// 执行撤销所有令牌
	if err := h.authService.LogoutAll(userID); err != nil {
		return utils.SystemError(c, err)
	}

	return utils.Success(c, map[string]string{
		"message": "已撤销所有令牌",
	}, "已撤销所有令牌")
}

// GetProfile 获取用户信息
func (h *AuthHandler) GetProfile(c echo.Context) error {
	// 从上下文中获取用户信息
	userID := c.Get("user_id").(uint)
	username := c.Get("username").(string)
	email := c.Get("email").(string)
	role := c.Get("role").(string)

	user := map[string]interface{}{
		"id":       userID,
		"username": username,
		"email":    email,
		"role":     role,
	}

	return utils.Success(c, user, "获取用户信息成功")
}

// ValidateToken 验证令牌
func (h *AuthHandler) ValidateToken(c echo.Context) error {
	// 从请求头获取 Authorization
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return utils.Unauthorized(c, "缺少认证令牌")
	}

	// 提取 token
	token := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	// 验证 Access Token
	claims, err := h.authService.ValidateAccessToken(token)
	if err != nil {
		return utils.Unauthorized(c, "认证令牌无效")
	}

	// 检查 token 是否过期
	if utils.IsAccessTokenExpiredWithClaims(claims) {
		return utils.Unauthorized(c, "认证令牌已过期")
	}

	return utils.Success(c, map[string]interface{}{
		"valid":    true,
		"user_id":  claims.UserID,
		"username": claims.Username,
		"email":    claims.Email,
		"role":     claims.Role,
	}, "令牌验证成功")
}
