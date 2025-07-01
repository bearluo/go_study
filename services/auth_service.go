package services

import (
	"errors"

	"go-study/db/models"
	"go-study/db/repositories"
	"go-study/utils"

	"golang.org/x/crypto/bcrypt"
)

// IAuthService 认证服务接口
type IAuthService interface {
	Register(req *RegisterRequest) (*RegisterResponse, error)
	Login(req *LoginRequest) (*LoginResponse, error)
	RefreshToken(req *RefreshTokenRequest) (*LoginResponse, error)
	Logout(req *LogoutRequest) error
	LogoutAll(userID uint) error
	ValidateAccessToken(tokenString string) (*utils.JWTClaims, error)
	GetUserFromToken(tokenString string) (*models.User, error)
	CleanupExpiredTokens() error
	CleanupRevokedTokens() error
}

// AuthService 认证服务
type AuthService struct {
	userRepo         repositories.UserRepository
	refreshTokenRepo repositories.RefreshTokenRepositoryInterface
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo repositories.UserRepository, refreshTokenRepo repositories.RefreshTokenRepositoryInterface) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

// RegisterRequest 注册请求结构体
type RegisterRequest struct {
	Name     string `json:"name" form:"name" validate:"required,username"`
	Email    string `json:"email" form:"email" validate:"required,email,max=50"`
	Password string `json:"password" form:"password" validate:"required,password"`
}

// RegisterResponse 注册响应结构体
type RegisterResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email,max=50"`
	Password string `json:"password" form:"password" validate:"required,password"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// LogoutRequest 登出请求
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Register 用户注册
func (s *AuthService) Register(req *RegisterRequest) (*RegisterResponse, error) {
	// 检查邮箱是否已存在
	exists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("邮箱已存在")
	}

	// 检查用户名是否已存在
	exists, err = s.userRepo.ExistsByName(req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	// 保存用户到数据库
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 生成令牌对
	tokenPair, err := utils.GenerateTokenPair(user.ID, user.Name, user.Email, user.Role, s.refreshTokenRepo)
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 查找用户
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("密码错误")
	}

	// 生成令牌对
	tokenPair, err := utils.GenerateTokenPair(user.ID, user.Name, user.Email, user.Role, s.refreshTokenRepo)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

// RefreshToken 刷新访问令牌
func (s *AuthService) RefreshToken(req *RefreshTokenRequest) (*LoginResponse, error) {
	// 从数据库查找 Refresh Token
	refreshToken, err := s.refreshTokenRepo.FindByToken(req.RefreshToken)
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

	// 获取用户信息
	user, err := s.userRepo.GetByID(refreshToken.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// 刷新令牌对
	tokenPair, err := utils.RefreshAccessTokenWithUserInfo(
		req.RefreshToken,
		user.ID,
		user.Name,
		user.Email,
		user.Role,
		s.refreshTokenRepo,
	)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

// Logout 用户登出
func (s *AuthService) Logout(req *LogoutRequest) error {
	// 撤销 Refresh Token
	return utils.RevokeRefreshToken(req.RefreshToken, s.refreshTokenRepo)
}

// LogoutAll 撤销用户的所有令牌
func (s *AuthService) LogoutAll(userID uint) error {
	return utils.RevokeAllUserTokens(userID, s.refreshTokenRepo)
}

// ValidateAccessToken 验证访问令牌
func (s *AuthService) ValidateAccessToken(tokenString string) (*utils.JWTClaims, error) {
	return utils.ValidateAccessToken(tokenString)
}

// GetUserFromToken 从令牌获取用户信息
func (s *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	claims, err := utils.ValidateAccessToken(tokenString)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// CleanupExpiredTokens 清理过期的令牌
func (s *AuthService) CleanupExpiredTokens() error {
	return utils.CleanupExpiredTokens(s.refreshTokenRepo)
}

// CleanupRevokedTokens 清理已撤销的令牌
func (s *AuthService) CleanupRevokedTokens() error {
	return utils.CleanupRevokedTokens(s.refreshTokenRepo)
}
