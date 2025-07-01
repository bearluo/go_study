package middleware

import (
	"net/http"
	"strings"

	"go-study/db/models"
	"go-study/services"
	"go-study/utils"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	authService services.IAuthService
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(authService services.IAuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth 要求认证的中间件
func (m *AuthMiddleware) RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 从请求头获取 Authorization
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "缺少认证令牌")
			}

			// 检查 Authorization 格式
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "认证令牌格式错误")
			}

			// 提取 token
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// 验证 Access Token
			claims, err := m.authService.ValidateAccessToken(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "认证令牌无效")
			}

			// 检查 token 是否过期
			if utils.IsAccessTokenExpiredWithClaims(claims) {
				return echo.NewHTTPError(http.StatusUnauthorized, "认证令牌已过期")
			}

			// 将用户信息存储到上下文中
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)
			c.Set("claims", claims)

			return next(c)
		}
	}
}

// RequireRole 要求特定角色的中间件
func (m *AuthMiddleware) RequireRole(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 先执行认证中间件
			authMiddleware := m.RequireAuth()
			authHandler := authMiddleware(next)

			// 检查角色
			role := c.Get("role").(string)
			if role != requiredRole {
				return echo.NewHTTPError(http.StatusForbidden, "权限不足")
			}

			return authHandler(c)
		}
	}
}

// RequireAdmin 要求管理员角色的中间件
func (m *AuthMiddleware) RequireAdmin() echo.MiddlewareFunc {
	return m.RequireRole(models.RoleAdmin)
}

// RequireUser 要求用户或更高权限的中间件
func (m *AuthMiddleware) RequireUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authMiddleware := m.RequireAuth()
			authHandler := authMiddleware(next)
			role := c.Get("role").(string)
			if !utils.IsUser(role) {
				return echo.NewHTTPError(http.StatusForbidden, "权限不足")
			}
			return authHandler(c)
		}
	}
}

// OptionalAuth 可选的认证中间件（不强制要求认证）
func (m *AuthMiddleware) OptionalAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 从请求头获取 Authorization
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				// 没有认证令牌，继续执行
				return next(c)
			}

			// 检查 Authorization 格式
			if !strings.HasPrefix(authHeader, "Bearer ") {
				// 格式错误，但不阻止继续执行
				return next(c)
			}

			// 提取 token
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// 尝试验证 Access Token
			claims, err := m.authService.ValidateAccessToken(token)
			if err != nil {
				// 验证失败，但不阻止继续执行
				return next(c)
			}

			// 检查 token 是否过期
			if utils.IsAccessTokenExpiredWithClaims(claims) {
				// 已过期，但不阻止继续执行
				return next(c)
			}

			// 将用户信息存储到上下文中
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)
			c.Set("claims", claims)

			return next(c)
		}
	}
}

// GetUserID 从上下文中获取用户ID
func GetUserID(c echo.Context) uint {
	if userID, ok := c.Get("user_id").(uint); ok {
		return userID
	}
	return 0
}

// GetUsername 从上下文中获取用户名
func GetUsername(c echo.Context) string {
	if username, ok := c.Get("username").(string); ok {
		return username
	}
	return ""
}

// GetEmail 从上下文中获取邮箱
func GetEmail(c echo.Context) string {
	if email, ok := c.Get("email").(string); ok {
		return email
	}
	return ""
}

// GetRole 从上下文中获取角色
func GetRole(c echo.Context) string {
	if role, ok := c.Get("role").(string); ok {
		return role
	}
	return ""
}

// GetClaims 从上下文中获取 JWT Claims
func GetClaims(c echo.Context) *utils.JWTClaims {
	if claims, ok := c.Get("claims").(*utils.JWTClaims); ok {
		return claims
	}
	return nil
}
