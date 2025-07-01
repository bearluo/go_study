package services

import (
	"errors"
	"go-study/db/models"
	"go-study/db/repositories"

	"golang.org/x/crypto/bcrypt"
)

// IUserService 用户服务接口
// 只包含handlers用到的方法
// （如需mockgen可用此注释）
//
//go:generate mockgen -source=user_service.go -destination=mock_user_service.go -package=services
type IUserService interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
}

// UserService 用户服务层
type UserService struct {
	userRepo repositories.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// Create 创建用户（保持向后兼容）
func (u *UserService) Create(user *models.User) error {
	// 检查邮箱是否已存在
	exists, err := u.userRepo.ExistsByEmail(user.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("邮箱已存在")
	}

	// 检查用户名是否已存在
	exists, err = u.userRepo.ExistsByName(user.Name)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("用户名已存在")
	}

	// 如果密码未加密，则进行加密
	if !u.isPasswordHashed(user.Password) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	return u.userRepo.Create(user)
}

// isPasswordHashed 检查密码是否已经加密
func (u *UserService) isPasswordHashed(password string) bool {
	// bcrypt 加密的密码通常以 $2a$、$2b$ 或 $2y$ 开头
	return len(password) >= 60 && (password[:4] == "$2a$" || password[:4] == "$2b$" || password[:4] == "$2y$")
}

// GetByID 根据ID获取用户
func (u *UserService) GetByID(id uint) (*models.User, error) {
	return u.userRepo.GetByID(id)
}

// GetByEmail 根据邮箱获取用户
func (u *UserService) GetByEmail(email string) (*models.User, error) {
	return u.userRepo.GetByEmail(email)
}

// GetAll 获取所有用户
func (u *UserService) GetAll() ([]models.User, error) {
	return u.userRepo.GetAll()
}

// Update 更新用户
func (u *UserService) Update(user *models.User) error {
	// 检查用户是否存在
	existingUser, err := u.userRepo.GetByID(user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("用户不存在")
	}

	// 如果邮箱有变化，检查新邮箱是否已存在
	if existingUser.Email != user.Email {
		exists, err := u.userRepo.ExistsByEmail(user.Email)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("邮箱已存在")
		}
	}

	// 如果用户名有变化，检查新用户名是否已存在
	if existingUser.Name != user.Name {
		exists, err := u.userRepo.ExistsByName(user.Name)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("用户名已存在")
		}
	}

	return u.userRepo.Update(user)
}

// Delete 删除用户
func (u *UserService) Delete(id uint) error {
	// 检查用户是否存在
	existingUser, err := u.userRepo.GetByID(id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("用户不存在")
	}

	return u.userRepo.Delete(id)
}
