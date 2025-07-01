package services

import (
	"go-study/db/repositories"
)

// ServiceManager 服务管理器
type ServiceManager struct {
	UserService *UserService
	AuthService *AuthService
}

// NewServiceManager 创建服务管理器
func NewServiceManager(repoManager *repositories.RepositoryManager) *ServiceManager {
	return &ServiceManager{
		UserService: NewUserService(repoManager.User),
		AuthService: NewAuthService(repoManager.User, repoManager.RefreshToken),
	}
}

// GetUserService 获取用户服务
func (sm *ServiceManager) GetUserService() *UserService {
	return sm.UserService
}

// GetAuthService 获取认证服务
func (sm *ServiceManager) GetAuthService() *AuthService {
	return sm.AuthService
}
