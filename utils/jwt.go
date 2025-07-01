package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"strconv"
	"time"

	"go-study/db/models"
	"go-study/db/repositories"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims 自定义 JWT 声明结构（用于 Access Token）
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	JTI      string `json:"jti,omitempty"` // JWT ID，用于确保唯一性
	jwt.RegisteredClaims
}

// JWTConfig JWT 配置结构
type JWTConfig struct {
	AccessTokenSecret    string        // Access Token 密钥
	AccessTokenDuration  time.Duration // Access Token 有效期（短期，如15分钟）
	RefreshTokenDuration time.Duration // Refresh Token 有效期（长期，如7天）
}

// DefaultJWTConfig 默认 JWT 配置
var DefaultJWTConfig *JWTConfig

// 环境变量常量
const (
	EnvJWTSecretKey            = "JWT_SECRET_KEY"
	EnvJWTAccessTokenDuration  = "JWT_ACCESS_TOKEN_DURATION"
	EnvJWTRefreshTokenDuration = "JWT_REFRESH_TOKEN_DURATION"
)

// 默认值常量
const (
	DefaultJWTSecretKey            = "your-access-token-secret-key-here"
	DefaultJWTAccessTokenDuration  = 900    // 15分钟，单位：秒
	DefaultJWTRefreshTokenDuration = 604800 // 7天，单位：秒
)

// InitJWTConfig 初始化 JWT 配置
func InitJWTConfig() {
	DefaultJWTConfig = &JWTConfig{
		AccessTokenSecret:    getEnvOrDefault(EnvJWTSecretKey, DefaultJWTSecretKey),
		AccessTokenDuration:  time.Duration(getEnvIntOrDefault(EnvJWTAccessTokenDuration, DefaultJWTAccessTokenDuration)) * time.Second,
		RefreshTokenDuration: time.Duration(getEnvIntOrDefault(EnvJWTRefreshTokenDuration, DefaultJWTRefreshTokenDuration)) * time.Second,
	}
}

// getEnvOrDefault 从环境变量获取值，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntOrDefault 从环境变量获取整数值，如果不存在或解析失败则返回默认值
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetJWTConfig 获取当前 JWT 配置
func GetJWTConfig() *JWTConfig {
	if DefaultJWTConfig == nil {
		InitJWTConfig()
	}
	return DefaultJWTConfig
}

// TokenPair 令牌对结构
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // Access Token 过期时间（秒）
}

// generateJTI 生成唯一的JWT ID
func generateJTI() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// generateRefreshToken 生成刷新令牌
func generateRefreshToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GenerateTokenPair 生成令牌对（Access Token + Refresh Token）
func GenerateTokenPair(userID uint, username, email, role string, refreshTokenRepo repositories.RefreshTokenRepositoryInterface) (*TokenPair, error) {
	return GenerateTokenPairWithConfig(userID, username, email, role, refreshTokenRepo, GetJWTConfig())
}

// GenerateTokenPairWithConfig 使用自定义配置生成令牌对
func GenerateTokenPairWithConfig(userID uint, username, email, role string, refreshTokenRepo repositories.RefreshTokenRepositoryInterface, config *JWTConfig) (*TokenPair, error) {
	// 生成 Access Token
	accessToken, err := generateAccessToken(userID, username, email, role, config)
	if err != nil {
		return nil, err
	}

	// 生成 Refresh Token
	refreshTokenStr := generateRefreshToken()
	refreshToken := &models.RefreshToken{
		UserID:    userID,
		Token:     refreshTokenStr,
		ExpiresAt: models.Time{Time: time.Now().Add(config.RefreshTokenDuration)},
		IsRevoked: false,
	}

	// 保存 Refresh Token 到数据库
	if err := refreshTokenRepo.Create(refreshToken); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    int64(config.AccessTokenDuration.Seconds()),
	}, nil
}

// generateAccessToken 生成 Access Token
func generateAccessToken(userID uint, username, email, role string, config *JWTConfig) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Role:     role,
		JTI:      generateJTI(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-study-app",
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AccessTokenSecret))
}

// ValidateAccessToken 验证 Access Token
func ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	return ValidateAccessTokenWithConfig(tokenString, GetJWTConfig())
}

// ValidateAccessTokenWithConfig 使用自定义配置验证 Access Token
func ValidateAccessTokenWithConfig(tokenString string, config *JWTConfig) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.AccessTokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshAccessToken 使用 Refresh Token 刷新 Access Token
func RefreshAccessToken(refreshTokenStr string, refreshTokenRepo repositories.RefreshTokenRepositoryInterface) (*TokenPair, error) {
	return RefreshAccessTokenWithConfig(refreshTokenStr, refreshTokenRepo, GetJWTConfig())
}

// RefreshAccessTokenWithConfig 使用自定义配置刷新 Access Token
func RefreshAccessTokenWithConfig(refreshTokenStr string, refreshTokenRepo repositories.RefreshTokenRepositoryInterface, config *JWTConfig) (*TokenPair, error) {
	// 从数据库查找 Refresh Token
	refreshToken, err := refreshTokenRepo.FindByToken(refreshTokenStr)
	if err != nil {
		return nil, err
	}

	if refreshToken == nil {
		return nil, errors.New("refresh token not found")
	}

	// 检查 Refresh Token 是否有效
	if !refreshToken.IsValid() {
		return nil, errors.New("refresh token is invalid or expired")
	}

	// 撤销旧的 Refresh Token
	if err := refreshTokenRepo.RevokeToken(refreshTokenStr); err != nil {
		return nil, err
	}

	// 生成新的令牌对
	// 注意：这里需要从用户表获取用户信息，因为 Refresh Token 只存储了 UserID
	// 在实际使用中，你可能需要注入 UserRepository 或者传递用户信息
	// 这里简化处理，假设用户信息可以从其他地方获取
	return nil, errors.New("需要用户信息来生成新的令牌对，请使用 RefreshAccessTokenWithUserInfo 方法")
}

// RefreshAccessTokenWithUserInfo 使用 Refresh Token 和用户信息刷新 Access Token
func RefreshAccessTokenWithUserInfo(refreshTokenStr string, userID uint, username, email, role string, refreshTokenRepo repositories.RefreshTokenRepositoryInterface) (*TokenPair, error) {
	return RefreshAccessTokenWithUserInfoAndConfig(refreshTokenStr, userID, username, email, role, refreshTokenRepo, GetJWTConfig())
}

// RefreshAccessTokenWithUserInfoAndConfig 使用自定义配置和用户信息刷新 Access Token
func RefreshAccessTokenWithUserInfoAndConfig(refreshTokenStr string, userID uint, username, email, role string, refreshTokenRepo repositories.RefreshTokenRepositoryInterface, config *JWTConfig) (*TokenPair, error) {
	// 从数据库查找 Refresh Token
	refreshToken, err := refreshTokenRepo.FindByToken(refreshTokenStr)
	if err != nil {
		return nil, err
	}

	if refreshToken == nil {
		return nil, errors.New("refresh token not found")
	}

	// 检查 Refresh Token 是否有效
	if !refreshToken.IsValid() {
		return nil, errors.New("refresh token is invalid or expired")
	}

	// 验证用户ID是否匹配
	if refreshToken.UserID != userID {
		return nil, errors.New("refresh token user mismatch")
	}

	// 撤销旧的 Refresh Token
	if err := refreshTokenRepo.RevokeToken(refreshTokenStr); err != nil {
		return nil, err
	}

	// 生成新的令牌对
	return GenerateTokenPairWithConfig(userID, username, email, role, refreshTokenRepo, config)
}

// RevokeRefreshToken 撤销 Refresh Token
func RevokeRefreshToken(refreshTokenStr string, refreshTokenRepo repositories.RefreshTokenRepositoryInterface) error {
	return refreshTokenRepo.RevokeToken(refreshTokenStr)
}

// RevokeAllUserTokens 撤销用户的所有 Refresh Token
func RevokeAllUserTokens(userID uint, refreshTokenRepo repositories.RefreshTokenRepositoryInterface) error {
	return refreshTokenRepo.RevokeAllUserTokens(userID)
}

// ExtractUserID 从 Access Token 中提取用户 ID
func ExtractUserID(tokenString string) (uint, error) {
	claims, err := ValidateAccessToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// ExtractUsername 从 Access Token 中提取用户名
func ExtractUsername(tokenString string) (string, error) {
	claims, err := ValidateAccessToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Username, nil
}

// ExtractRole 从 Access Token 中提取用户角色
func ExtractRole(tokenString string) (string, error) {
	claims, err := ValidateAccessToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Role, nil
}

// IsAccessTokenExpired 检查 Access Token 是否过期
func IsAccessTokenExpired(tokenString string) (bool, error) {
	claims, err := ValidateAccessToken(tokenString)
	if err != nil {
		return true, err
	}

	// 检查是否过期
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return true, nil
	}

	return false, nil
}

// IsAccessTokenExpiredWithConfig 使用自定义配置检查 Access Token 是否过期
func IsAccessTokenExpiredWithConfig(tokenString string, config *JWTConfig) (bool, error) {
	claims, err := ValidateAccessTokenWithConfig(tokenString, config)
	if err != nil {
		return true, err
	}

	// 检查是否过期
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return true, nil
	}

	return false, nil
}

// GetAccessTokenExpirationTime 获取 Access Token 过期时间
func GetAccessTokenExpirationTime(tokenString string) (*time.Time, error) {
	claims, err := ValidateAccessToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.ExpiresAt == nil {
		return nil, errors.New("token has no expiration time")
	}

	return &claims.ExpiresAt.Time, nil
}

// IsAccessTokenExpiredWithClaims 使用传入的 claims 检查 Access Token 是否过期
func IsAccessTokenExpiredWithClaims(claims *JWTClaims) bool {
	if claims == nil {
		return true
	}

	// 检查是否过期
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return true
	}

	return false
}

// CleanupExpiredTokens 清理过期的 Refresh Token
func CleanupExpiredTokens(refreshTokenRepo repositories.RefreshTokenRepositoryInterface) error {
	return refreshTokenRepo.DeleteExpiredTokens()
}

// CleanupRevokedTokens 清理已撤销的 Refresh Token
func CleanupRevokedTokens(refreshTokenRepo repositories.RefreshTokenRepositoryInterface) error {
	return refreshTokenRepo.DeleteRevokedTokens()
}
