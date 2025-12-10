package services

import (
	"github.com/halolight/halolight-api-go/internal/models"
	"gorm.io/gorm"
)

type RoleService interface {
	List() ([]models.Role, error)
	Get(id string) (*models.Role, error)
	Create(name, label, description string) (*models.Role, error)
	Update(id, name, label, description string) (*models.Role, error)
	Delete(id string) error
	AssignPermissions(roleID string, permissionIDs []string) (*models.Role, error)
}

type roleService struct {
	db *gorm.DB
}

func NewRoleService(db *gorm.DB) RoleService {
	return &roleService{db: db}
}

func (s *roleService) List() ([]models.Role, error) {
	var roles []models.Role
	err := s.db.Preload("Permissions.Permission").Find(&roles).Error
	return roles, err
}

func (s *roleService) Get(id string) (*models.Role, error) {
	var role models.Role
	err := s.db.Preload("Permissions.Permission").First(&role, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *roleService) Create(name, label, description string) (*models.Role, error) {
	role := &models.Role{
		Name:        name,
		Label:       label,
		Description: &description,
	}
	err := s.db.Create(role).Error
	return role, err
}

func (s *roleService) Update(id, name, label, description string) (*models.Role, error) {
	role, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		role.Name = name
	}
	if label != "" {
		role.Label = label
	}
	if description != "" {
		role.Description = &description
	}

	err = s.db.Save(role).Error
	return role, err
}

func (s *roleService) Delete(id string) error {
	return s.db.Delete(&models.Role{}, "id = ?", id).Error
}

func (s *roleService) AssignPermissions(roleID string, permissionIDs []string) (*models.Role, error) {
	// Remove existing permissions
	if err := s.db.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
		return nil, err
	}

	// Add new permissions
	for _, permID := range permissionIDs {
		rp := &models.RolePermission{
			RoleID:       roleID,
			PermissionID: permID,
		}
		if err := s.db.Create(rp).Error; err != nil {
			return nil, err
		}
	}

	return s.Get(roleID)
}
