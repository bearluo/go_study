package repositories

import (
	"errors"
	"go-study/db/models"
	"time"

	"gorm.io/gorm"
)

// RefreshTokenRepositoryInterface 刷新令牌仓库接口
type RefreshTokenRepositoryInterface interface {
	Create(refreshToken *models.RefreshToken) error
	FindByToken(token string) (*models.RefreshToken, error)
	FindByUserID(userID uint) ([]models.RefreshToken, error)
	RevokeToken(token string) error
	RevokeAllUserTokens(userID uint) error
	DeleteExpiredTokens() error
	DeleteRevokedTokens() error
	CountByUserID(userID uint) (int64, error)
}

// RefreshTokenRepository 刷新令牌仓库
type RefreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository 创建新的刷新令牌仓库
func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepositoryInterface {
	return &RefreshTokenRepository{db: db}
}

// Create 创建刷新令牌
func (r *RefreshTokenRepository) Create(refreshToken *models.RefreshToken) error {
	return r.db.Create(refreshToken).Error
}

// FindByToken 根据令牌查找刷新令牌
func (r *RefreshTokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &refreshToken, nil
}

// FindByUserID 根据用户ID查找所有刷新令牌
func (r *RefreshTokenRepository) FindByUserID(userID uint) ([]models.RefreshToken, error) {
	var refreshTokens []models.RefreshToken
	err := r.db.Where("user_id = ?", userID).Find(&refreshTokens).Error
	return refreshTokens, err
}

// RevokeToken 撤销指定的刷新令牌
func (r *RefreshTokenRepository) RevokeToken(token string) error {
	return r.db.Model(&models.RefreshToken{}).Where("token = ?", token).Update("is_revoked", true).Error
}

// RevokeAllUserTokens 撤销用户的所有刷新令牌
func (r *RefreshTokenRepository) RevokeAllUserTokens(userID uint) error {
	return r.db.Model(&models.RefreshToken{}).Where("user_id = ?", userID).Update("is_revoked", true).Error
}

// DeleteExpiredTokens 删除过期的刷新令牌
func (r *RefreshTokenRepository) DeleteExpiredTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{}).Error
}

// DeleteRevokedTokens 删除已撤销的刷新令牌
func (r *RefreshTokenRepository) DeleteRevokedTokens() error {
	return r.db.Where("is_revoked = ?", true).Delete(&models.RefreshToken{}).Error
}

// CountByUserID 统计用户的刷新令牌数量
func (r *RefreshTokenRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.RefreshToken{}).Where("user_id = ? AND is_revoked = ?", userID, false).Count(&count).Error
	return count, err
}
