package repository

import (
	"github.com/halolight/halolight-api-go/internal/models"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(token *models.RefreshToken) error
	FindByToken(token string) (*models.RefreshToken, error)
	FindByUserID(userID string) ([]models.RefreshToken, error)
	Delete(id string) error
	DeleteByToken(token string) error
	DeleteByUserID(userID string) error
	DeleteExpired() error
}

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *refreshTokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.Where("token = ?", token).Preload("User").First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *refreshTokenRepository) FindByUserID(userID string) ([]models.RefreshToken, error) {
	var tokens []models.RefreshToken
	err := r.db.Where("user_id = ?", userID).Find(&tokens).Error
	return tokens, err
}

func (r *refreshTokenRepository) Delete(id string) error {
	return r.db.Delete(&models.RefreshToken{}, "id = ?", id).Error
}

func (r *refreshTokenRepository) DeleteByToken(token string) error {
	return r.db.Delete(&models.RefreshToken{}, "token = ?", token).Error
}

func (r *refreshTokenRepository) DeleteByUserID(userID string) error {
	return r.db.Delete(&models.RefreshToken{}, "user_id = ?", userID).Error
}

func (r *refreshTokenRepository) DeleteExpired() error {
	return r.db.Where("expires_at < NOW()").Delete(&models.RefreshToken{}).Error
}
