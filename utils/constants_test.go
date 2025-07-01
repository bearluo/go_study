package utils

import (
	"testing"
	"time"
)

func TestTimeConstants(t *testing.T) {
	// 测试时间常量
	if DefaultAccessTokenDuration != 15*60 {
		t.Errorf("DefaultAccessTokenDuration 应该是 15*60 秒，实际是 %d", DefaultAccessTokenDuration)
	}

	if DefaultRefreshTokenDuration != 7*24*60*60 {
		t.Errorf("DefaultRefreshTokenDuration 应该是 7*24*60*60 秒，实际是 %d", DefaultRefreshTokenDuration)
	}

	if TokenCleanupInterval != 24 {
		t.Errorf("TokenCleanupInterval 应该是 24 小时，实际是 %d", TokenCleanupInterval)
	}
}

func TestSystemConstants(t *testing.T) {
	// 测试系统常量
	if DefaultPageSize != 20 {
		t.Errorf("DefaultPageSize 应该是 20，实际是 %d", DefaultPageSize)
	}

	if MaxPageSize != 100 {
		t.Errorf("MaxPageSize 应该是 100，实际是 %d", MaxPageSize)
	}

	if MinPasswordLength != 6 {
		t.Errorf("MinPasswordLength 应该是 6，实际是 %d", MinPasswordLength)
	}

	if MaxPasswordLength != 18 {
		t.Errorf("MaxPasswordLength 应该是 18，实际是 %d", MaxPasswordLength)
	}

	if MinUsernameLength != 2 {
		t.Errorf("MinUsernameLength 应该是 2，实际是 %d", MinUsernameLength)
	}

	if MaxUsernameLength != 10 {
		t.Errorf("MaxUsernameLength 应该是 10，实际是 %d", MaxUsernameLength)
	}

	if MaxEmailLength != 50 {
		t.Errorf("MaxEmailLength 应该是 50，实际是 %d", MaxEmailLength)
	}
}

func TestRegexConstants(t *testing.T) {
	// 测试正则表达式常量
	if RegexUsername == "" {
		t.Error("RegexUsername 不能为空")
	}

	if RegexPassword == "" {
		t.Error("RegexPassword 不能为空")
	}

	if RegexEmail == "" {
		t.Error("RegexEmail 不能为空")
	}
}

func TestCacheKeyConstants(t *testing.T) {
	// 测试缓存键常量
	if CacheKeyUser == "" {
		t.Error("CacheKeyUser 不能为空")
	}

	if CacheKeyToken == "" {
		t.Error("CacheKeyToken 不能为空")
	}

	if CacheKeySession == "" {
		t.Error("CacheKeySession 不能为空")
	}

	// 检查缓存键前缀
	if CacheKeyUser != "user:" {
		t.Errorf("CacheKeyUser 应该是 'user:'，实际是 '%s'", CacheKeyUser)
	}

	if CacheKeyToken != "token:" {
		t.Errorf("CacheKeyToken 应该是 'token:'，实际是 '%s'", CacheKeyToken)
	}

	if CacheKeySession != "session:" {
		t.Errorf("CacheKeySession 应该是 'session:'，实际是 '%s'", CacheKeySession)
	}
}

func TestLogLevelConstants(t *testing.T) {
	// 测试日志级别常量
	logLevels := []string{LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError, LogLevelFatal}
	expectedLevels := []string{"debug", "info", "warn", "error", "fatal"}

	for i, level := range logLevels {
		if level != expectedLevels[i] {
			t.Errorf("日志级别 %d 应该是 '%s'，实际是 '%s'", i, expectedLevels[i], level)
		}
	}
}

func TestEnvironmentConstants(t *testing.T) {
	// 测试环境常量
	environments := []string{EnvDevelopment, EnvTesting, EnvStaging, EnvProduction}
	expectedEnvs := []string{"development", "testing", "staging", "production"}

	for i, env := range environments {
		if env != expectedEnvs[i] {
			t.Errorf("环境 %d 应该是 '%s'，实际是 '%s'", i, expectedEnvs[i], env)
		}
	}
}

func TestTimeDurationConversion(t *testing.T) {
	// 测试时间常量转换为 time.Duration
	accessTokenDuration := time.Duration(DefaultAccessTokenDuration) * time.Second
	expectedAccessTokenDuration := 15 * time.Minute

	if accessTokenDuration != expectedAccessTokenDuration {
		t.Errorf("Access Token 持续时间应该是 %v，实际是 %v", expectedAccessTokenDuration, accessTokenDuration)
	}

	refreshTokenDuration := time.Duration(DefaultRefreshTokenDuration) * time.Second
	expectedRefreshTokenDuration := 7 * 24 * time.Hour

	if refreshTokenDuration != expectedRefreshTokenDuration {
		t.Errorf("Refresh Token 持续时间应该是 %v，实际是 %v", expectedRefreshTokenDuration, refreshTokenDuration)
	}
}
