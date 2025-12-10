package services

import (
	"github.com/halolight/halolight-api-go/internal/models"
	"gorm.io/gorm"
)

type FolderService interface {
	List(userID string, parentID *string) ([]models.Folder, error)
	Get(id string) (*models.Folder, error)
	GetTree(userID string) ([]FolderTreeNode, error)
	Create(name string, parentID *string, teamID *string, ownerID string) (*models.Folder, error)
	Rename(id, name string) (*models.Folder, error)
	Delete(id string) error
	IsOwner(folderID, userID string) bool
}

type FolderTreeNode struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	Path     string           `json:"path"`
	ParentID *string          `json:"parentId"`
	Children []FolderTreeNode `json:"children"`
}

type folderService struct {
	db *gorm.DB
}

func NewFolderService(db *gorm.DB) FolderService {
	return &folderService{db: db}
}

func (s *folderService) List(userID string, parentID *string) ([]models.Folder, error) {
	var folders []models.Folder
	query := s.db.Where("owner_id = ?", userID)

	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	err := query.Order("name ASC").Find(&folders).Error
	return folders, err
}

func (s *folderService) Get(id string) (*models.Folder, error) {
	var folder models.Folder
	err := s.db.Preload("Parent").Preload("Children").Preload("Files").First(&folder, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (s *folderService) GetTree(userID string) ([]FolderTreeNode, error) {
	var folders []models.Folder
	err := s.db.Where("owner_id = ?", userID).Order("path ASC").Find(&folders).Error
	if err != nil {
		return nil, err
	}

	// Build tree
	folderMap := make(map[string]*FolderTreeNode)
	var roots []FolderTreeNode

	for _, f := range folders {
		node := FolderTreeNode{
			ID:       f.ID,
			Name:     f.Name,
			Path:     f.Path,
			ParentID: f.ParentID,
			Children: []FolderTreeNode{},
		}
		folderMap[f.ID] = &node
	}

	for _, f := range folders {
		node := folderMap[f.ID]
		if f.ParentID == nil {
			roots = append(roots, *node)
		} else if parent, ok := folderMap[*f.ParentID]; ok {
			parent.Children = append(parent.Children, *node)
		}
	}

	return roots, nil
}

func (s *folderService) Create(name string, parentID *string, teamID *string, ownerID string) (*models.Folder, error) {
	path := "/" + name

	if parentID != nil {
		var parent models.Folder
		if err := s.db.First(&parent, "id = ?", *parentID).Error; err == nil {
			path = parent.Path + "/" + name
		}
	}

	folder := &models.Folder{
		Name:     name,
		Path:     path,
		ParentID: parentID,
		TeamID:   teamID,
		OwnerID:  ownerID,
	}

	err := s.db.Create(folder).Error
	return folder, err
}

func (s *folderService) Rename(id, name string) (*models.Folder, error) {
	folder, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	// Update path
	pathParts := splitPath(folder.Path)
	if len(pathParts) > 0 {
		pathParts[len(pathParts)-1] = name
	}
	newPath := "/" + joinPath(pathParts)

	err = s.db.Model(&models.Folder{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name": name,
		"path": newPath,
	}).Error

	if err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *folderService) Delete(id string) error {
	return s.db.Delete(&models.Folder{}, "id = ?", id).Error
}

func (s *folderService) IsOwner(folderID, userID string) bool {
	var folder models.Folder
	err := s.db.Select("owner_id").First(&folder, "id = ?", folderID).Error
	if err != nil {
		return false
	}
	return folder.OwnerID == userID
}

func splitPath(path string) []string {
	// Simple split
	result := []string{}
	current := ""
	for _, c := range path {
		if c == '/' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func joinPath(parts []string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += "/"
		}
		result += p
	}
	return result
}
