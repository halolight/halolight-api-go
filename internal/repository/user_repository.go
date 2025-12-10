package repository

import (
	"errors"

	"github.com/halolight/halolight-api-go/internal/models"
	"gorm.io/gorm"
)

var (
	ErrNotFound      = errors.New("record not found")
	ErrDuplicateKey  = errors.New("duplicate key violation")
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByStringID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	List(offset, limit int) ([]models.User, int64, error)
	Update(user *models.User) error
	Delete(id uint) error
	BatchDelete(ids []string) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *models.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		// Check for unique constraint violation
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrDuplicateKey
		}
		return err
	}
	return nil
}

func (r *userRepo) GetByID(id uint) (*models.User, error) {
	var u models.User
	if err := r.db.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetByEmail(email string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetByUsername(username string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("username = ?", username).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) List(offset, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Get total count
	if err := r.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.Offset(offset).Limit(limit).Order("id asc").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepo) Update(user *models.User) error {
	err := r.db.Save(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrDuplicateKey
		}
		return err
	}
	return nil
}

func (r *userRepo) Delete(id uint) error {
	result := r.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *userRepo) GetByStringID(id string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("id = ?", id).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) BatchDelete(ids []string) error {
	return r.db.Where("id IN ?", ids).Delete(&models.User{}).Error
}
