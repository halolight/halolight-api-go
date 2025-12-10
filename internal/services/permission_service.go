package services

import (
	"github.com/halolight/halolight-api-go/internal/models"
	"gorm.io/gorm"
)

type PermissionService interface {
	List() ([]models.Permission, error)
	Get(id string) (*models.Permission, error)
	Create(action, resource, description string) (*models.Permission, error)
	Delete(id string) error
}

type permissionService struct {
	db *gorm.DB
}

func NewPermissionService(db *gorm.DB) PermissionService {
	return &permissionService{db: db}
}

func (s *permissionService) List() ([]models.Permission, error) {
	var permissions []models.Permission
	err := s.db.Order("resource, action").Find(&permissions).Error
	return permissions, err
}

func (s *permissionService) Get(id string) (*models.Permission, error) {
	var permission models.Permission
	err := s.db.First(&permission, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (s *permissionService) Create(action, resource, description string) (*models.Permission, error) {
	permission := &models.Permission{
		Action:      action,
		Resource:    resource,
		Description: &description,
	}
	err := s.db.Create(permission).Error
	return permission, err
}

func (s *permissionService) Delete(id string) error {
	return s.db.Delete(&models.Permission{}, "id = ?", id).Error
}
