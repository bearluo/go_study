package utils

import (
	"testing"
	"time"

	"go-study/db/models"
	"go-study/db/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRefreshTokenRepository 模拟 Refresh Token 仓库
type MockRefreshTokenRepository struct {
	mock.Mock
}

// 确保 MockRefreshTokenRepository 实现了 RefreshTokenRepositoryInterface 接口
var _ repositories.RefreshTokenRepositoryInterface = (*MockRefreshTokenRepository)(nil)

func (m *MockRefreshTokenRepository) Create(refreshToken *models.RefreshToken) error {
	args := m.Called(refreshToken)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) FindByUserID(userID uint) ([]models.RefreshToken, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) RevokeToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeAllUserTokens(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteExpiredTokens() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteRevokedTokens() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func TestGenerateTokenPair(t *testing.T) {
	mockRepo := new(MockRefreshTokenRepository)

	// 设置模拟行为
	mockRepo.On("Create", mock.AnythingOfType("*models.RefreshToken")).Return(nil)

	// 测试生成令牌对
	tokenPair, err := GenerateTokenPair(1, "testuser", "test@example.com", "user", mockRepo)

	assert.NoError(t, err)
	assert.NotNil(t, tokenPair)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
	assert.Equal(t, int64(900), tokenPair.ExpiresIn) // 15分钟 = 900秒

	mockRepo.AssertExpectations(t)
}

func TestValidateAccessToken(t *testing.T) {
	mockRepo := new(MockRefreshTokenRepository)

	// 生成令牌对
	mockRepo.On("Create", mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	tokenPair, err := GenerateTokenPair(1, "testuser", "test@example.com", "user", mockRepo)
	assert.NoError(t, err)

	// 测试验证 Access Token
	claims, err := ValidateAccessToken(tokenPair.AccessToken)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, "testuser", claims.Username)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "user", claims.Role)

	mockRepo.AssertExpectations(t)
}

func TestValidateAccessToken_InvalidToken(t *testing.T) {
	// 测试无效的 Access Token
	claims, err := ValidateAccessToken("invalid.token.here")

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestRefreshAccessTokenWithUserInfo(t *testing.T) {
	mockRepo := new(MockRefreshTokenRepository)

	// 生成初始令牌对
	mockRepo.On("Create", mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	tokenPair, err := GenerateTokenPair(1, "testuser", "test@example.com", "user", mockRepo)
	assert.NoError(t, err)

	// 设置模拟行为：查找 Refresh Token
	refreshToken := &models.RefreshToken{
		UserID:    1,
		Token:     tokenPair.RefreshToken,
		ExpiresAt: models.Time{Time: time.Now().Add(7 * 24 * time.Hour)},
		IsRevoked: false,
	}
	mockRepo.On("FindByToken", tokenPair.RefreshToken).Return(refreshToken, nil)
	mockRepo.On("RevokeToken", tokenPair.RefreshToken).Return(nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.RefreshToken")).Return(nil)

	// 测试刷新令牌
	newTokenPair, err := RefreshAccessTokenWithUserInfo(
		tokenPair.RefreshToken,
		1,
		"testuser",
		"test@example.com",
		"user",
		mockRepo,
	)

	assert.NoError(t, err)
	assert.NotNil(t, newTokenPair)
	assert.NotEmpty(t, newTokenPair.AccessToken)
	assert.NotEmpty(t, newTokenPair.RefreshToken)
	assert.NotEqual(t, tokenPair.AccessToken, newTokenPair.AccessToken)
	assert.NotEqual(t, tokenPair.RefreshToken, newTokenPair.RefreshToken)

	mockRepo.AssertExpectations(t)
}

func TestRefreshAccessTokenWithUserInfo_InvalidToken(t *testing.T) {
	mockRepo := new(MockRefreshTokenRepository)

	// 设置模拟行为：找不到 Refresh Token
	mockRepo.On("FindByToken", "invalid-token").Return(nil, nil)

	// 测试无效的 Refresh Token
	tokenPair, err := RefreshAccessTokenWithUserInfo(
		"invalid-token",
		1,
		"testuser",
		"test@example.com",
		"user",
		mockRepo,
	)

	assert.Error(t, err)
	assert.Nil(t, tokenPair)
	assert.Contains(t, err.Error(), "refresh token not found")

	mockRepo.AssertExpectations(t)
}

func TestRefreshAccessTokenWithUserInfo_ExpiredToken(t *testing.T) {
	mockRepo := new(MockRefreshTokenRepository)

	// 设置模拟行为：过期的 Refresh Token
	expiredToken := &models.RefreshToken{
		UserID:    1,
		Token:     "expired-token",
		ExpiresAt: models.Time{Time: time.Now().Add(-1 * time.Hour)}, // 1小时前过期
		IsRevoked: false,
	}
	mockRepo.On("FindByToken", "expired-token").Return(expiredToken, nil)

	// 测试过期的 Refresh Token
	tokenPair, err := RefreshAccessTokenWithUserInfo(
		"expired-token",
		1,
		"testuser",
		"test@example.com",
		"user",
		mockRepo,
	)

	assert.Error(t, err)
	assert.Nil(t, tokenPair)
	assert.Contains(t, err.Error(), "refresh token is invalid or expired")

	mockRepo.AssertExpectations(t)
}

func TestRefreshAccessTokenWithUserInfo_RevokedToken(t *testing.T) {
	mockRepo := new(MockRefreshTokenRepository)

	// 设置模拟行为：已撤销的 Refresh Token
	revokedToken := &models.RefreshToken{
		UserID:    1,
		Token:     "revoked-token",
		ExpiresAt: models.Time{Time: time.Now().Add(7 * 24 * time.Hour)},
		IsRevoked: true,
	}
	mockRepo.On("FindByToken", "revoked-token").Return(revokedToken, nil)

	// 测试已撤销的 Refresh Token
	tokenPair, err := RefreshAccessTokenWithUserInfo(
		"revoked-token",
		1,
		"testuser",
		"test@example.com",
		"user",
		mockRepo,
	)

	assert.Error(t, err)
	assert.Nil(t, tokenPair)
	assert.Contains(t, err.Error(), "refresh token is invalid or expired")

	mockRepo.AssertExpectations(t)
}

func TestRevokeRefreshToken(t *testing.T) {
	mockRepo := new(MockRefreshTokenRepository)

	// 设置模拟行为
	mockRepo.On("RevokeToken", "test-token").Return(nil)

	// 测试撤销 Refresh Token
	err := RevokeRefreshToken("test-token", mockRepo)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestRevokeAllUserTokens(t *testing.T) {
	mockRepo := new(MockRefreshTokenRepository)

	// 设置模拟行为
	mockRepo.On("RevokeAllUserTokens", uint(1)).Return(nil)

	// 测试撤销用户所有令牌
	err := RevokeAllUserTokens(1, mockRepo)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestIsAccessTokenExpired(t *testing.T) {
	mockRepo := new(MockRefreshTokenRepository)

	// 生成令牌对
	mockRepo.On("Create", mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	tokenPair, err := GenerateTokenPair(1, "testuser", "test@example.com", "user", mockRepo)
	assert.NoError(t, err)

	// 测试未过期的 Access Token
	expired, err := IsAccessTokenExpired(tokenPair.AccessToken)

	assert.NoError(t, err)
	assert.False(t, expired)

	mockRepo.AssertExpectations(t)
}

func TestIsAccessTokenExpired_InvalidToken(t *testing.T) {
	// 测试无效的 Access Token
	expired, err := IsAccessTokenExpired("invalid.token.here")

	assert.Error(t, err)
	assert.True(t, expired)
}

func TestExtractUserInfo(t *testing.T) {
	mockRepo := new(MockRefreshTokenRepository)

	// 生成令牌对
	mockRepo.On("Create", mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	tokenPair, err := GenerateTokenPair(1, "testuser", "test@example.com", "user", mockRepo)
	assert.NoError(t, err)

	// 测试提取用户ID
	userID, err := ExtractUserID(tokenPair.AccessToken)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), userID)

	// 测试提取用户名
	username, err := ExtractUsername(tokenPair.AccessToken)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", username)

	// 测试提取角色
	role, err := ExtractRole(tokenPair.AccessToken)
	assert.NoError(t, err)
	assert.Equal(t, "user", role)

	mockRepo.AssertExpectations(t)
}
