package utils

// 系统常量
const (
	// DefaultPageSize 默认分页大小
	DefaultPageSize = 20

	// MaxPageSize 最大分页大小
	MaxPageSize = 100

	// MinPasswordLength 最小密码长度
	MinPasswordLength = 6

	// MaxPasswordLength 最大密码长度
	MaxPasswordLength = 18

	// MinUsernameLength 最小用户名长度
	MinUsernameLength = 2

	// MaxUsernameLength 最大用户名长度
	MaxUsernameLength = 10

	// MaxEmailLength 最大邮箱长度
	MaxEmailLength = 50
)

// 时间常量
const (
	// DefaultAccessTokenDuration 默认 Access Token 有效期
	DefaultAccessTokenDuration = 15 * 60 // 15分钟（秒）

	// DefaultRefreshTokenDuration 默认 Refresh Token 有效期
	DefaultRefreshTokenDuration = 7 * 24 * 60 * 60 // 7天（秒）

	// TokenCleanupInterval 令牌清理间隔（小时）
	TokenCleanupInterval = 24
)

// 正则表达式常量
const (
	// RegexUsername 用户名正则表达式（字母、数字）
	RegexUsername = `^[a-zA-Z0-9]+$`

	// RegexPassword 密码正则表达式（字母、数字、特殊字符）
	RegexPassword = `^[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]+$`

	// RegexEmail 邮箱正则表达式
	RegexEmail = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

// 缓存键前缀
const (
	// CacheKeyUser 用户缓存键前缀
	CacheKeyUser = "user:"

	// CacheKeyToken 令牌缓存键前缀
	CacheKeyToken = "token:"

	// CacheKeySession 会话缓存键前缀
	CacheKeySession = "session:"
)

// 日志常量
const (
	// LogLevelDebug 调试级别
	LogLevelDebug = "debug"

	// LogLevelInfo 信息级别
	LogLevelInfo = "info"

	// LogLevelWarn 警告级别
	LogLevelWarn = "warn"

	// LogLevelError 错误级别
	LogLevelError = "error"

	// LogLevelFatal 致命错误级别
	LogLevelFatal = "fatal"
)

// 环境常量
const (
	// EnvDevelopment 开发环境
	EnvDevelopment = "development"

	// EnvTesting 测试环境
	EnvTesting = "testing"

	// EnvStaging 预发布环境
	EnvStaging = "staging"

	// EnvProduction 生产环境
	EnvProduction = "production"
)
