package utils

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitJWTConfig(t *testing.T) {
	// 清理环境变量
	os.Unsetenv(EnvJWTSecretKey)
	os.Unsetenv(EnvJWTAccessTokenDuration)
	os.Unsetenv(EnvJWTRefreshTokenDuration)

	// 测试默认配置
	InitJWTConfig()
	config := GetJWTConfig()

	assert.NotNil(t, config)
	assert.Equal(t, DefaultJWTSecretKey, config.AccessTokenSecret)
	assert.Equal(t, time.Duration(DefaultJWTAccessTokenDuration)*time.Second, config.AccessTokenDuration)
	assert.Equal(t, time.Duration(DefaultJWTRefreshTokenDuration)*time.Second, config.RefreshTokenDuration)
}

func TestInitJWTConfig_WithEnvironmentVariables(t *testing.T) {
	// 设置环境变量
	os.Setenv(EnvJWTSecretKey, "test-secret-key")
	os.Setenv(EnvJWTAccessTokenDuration, "1800")   // 30分钟
	os.Setenv(EnvJWTRefreshTokenDuration, "86400") // 1天

	// 测试环境变量配置
	InitJWTConfig()
	config := GetJWTConfig()

	assert.NotNil(t, config)
	assert.Equal(t, "test-secret-key", config.AccessTokenSecret)
	assert.Equal(t, time.Duration(1800)*time.Second, config.AccessTokenDuration)
	assert.Equal(t, time.Duration(86400)*time.Second, config.RefreshTokenDuration)

	// 清理环境变量
	os.Unsetenv(EnvJWTSecretKey)
	os.Unsetenv(EnvJWTAccessTokenDuration)
	os.Unsetenv(EnvJWTRefreshTokenDuration)
}

func TestInitJWTConfig_InvalidDuration(t *testing.T) {
	// 设置无效的环境变量
	os.Setenv(EnvJWTAccessTokenDuration, "invalid")
	os.Setenv(EnvJWTRefreshTokenDuration, "also-invalid")

	// 测试无效配置时使用默认值
	InitJWTConfig()
	config := GetJWTConfig()

	assert.NotNil(t, config)
	assert.Equal(t, DefaultJWTSecretKey, config.AccessTokenSecret)
	assert.Equal(t, time.Duration(DefaultJWTAccessTokenDuration)*time.Second, config.AccessTokenDuration)
	assert.Equal(t, time.Duration(DefaultJWTRefreshTokenDuration)*time.Second, config.RefreshTokenDuration)

	// 清理环境变量
	os.Unsetenv(EnvJWTAccessTokenDuration)
	os.Unsetenv(EnvJWTRefreshTokenDuration)
}

func TestGetJWTConfig_LazyInitialization(t *testing.T) {
	// 重置配置
	DefaultJWTConfig = nil

	// 测试懒加载初始化
	config := GetJWTConfig()

	assert.NotNil(t, config)
	assert.NotNil(t, DefaultJWTConfig)
	assert.Equal(t, DefaultJWTConfig, config)
}

func TestGetEnvOrDefault(t *testing.T) {
	// 测试环境变量存在
	os.Setenv("TEST_KEY", "test-value")
	assert.Equal(t, "test-value", getEnvOrDefault("TEST_KEY", "default"))

	// 测试环境变量不存在
	assert.Equal(t, "default", getEnvOrDefault("NONEXISTENT_KEY", "default"))

	// 清理
	os.Unsetenv("TEST_KEY")
}

func TestGetEnvIntOrDefault(t *testing.T) {
	// 测试有效整数
	os.Setenv("TEST_INT", "123")
	assert.Equal(t, 123, getEnvIntOrDefault("TEST_INT", 456))

	// 测试无效整数
	os.Setenv("TEST_INVALID", "not-a-number")
	assert.Equal(t, 456, getEnvIntOrDefault("TEST_INVALID", 456))

	// 测试环境变量不存在
	assert.Equal(t, 456, getEnvIntOrDefault("NONEXISTENT_INT", 456))

	// 清理
	os.Unsetenv("TEST_INT")
	os.Unsetenv("TEST_INVALID")
}
