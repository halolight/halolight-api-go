package services

import (
	"github.com/halolight/halolight-api-go/internal/models"
	"gorm.io/gorm"
)

type FileService interface {
	List(userID string, page, limit int, path, fileType, search, folderID string) ([]models.File, int64, error)
	Get(id string) (*models.File, error)
	Create(name, path, mimeType string, size int64, folderID, teamID, ownerID string) (*models.File, error)
	Rename(id, name string) (*models.File, error)
	Move(id, folderID, newPath string) (*models.File, error)
	Copy(id, ownerID string) (*models.File, error)
	ToggleFavorite(id string) (*models.File, error)
	Delete(id string) error
	DeleteMany(ids []string, userID string) error
	GetStorageInfo(userID string) (used, total, available int64, usedPercent int)
	GetDownloadURL(id string) string
	IsOwner(fileID, userID string) bool
}

type fileService struct {
	db *gorm.DB
}

func NewFileService(db *gorm.DB) FileService {
	return &fileService{db: db}
}

func (s *fileService) List(userID string, page, limit int, path, fileType, search, folderID string) ([]models.File, int64, error) {
	var files []models.File
	var total int64

	query := s.db.Model(&models.File{}).Where("owner_id = ?", userID)

	if folderID != "" {
		query = query.Where("folder_id = ?", folderID)
	}
	if path != "" {
		query = query.Where("path LIKE ?", path+"%")
	}
	if fileType != "" {
		query = query.Where("mime_type LIKE ?", fileType+"%")
	}
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Preload("Folder").Offset(offset).Limit(limit).Order("created_at DESC").Find(&files).Error

	return files, total, err
}

func (s *fileService) Get(id string) (*models.File, error) {
	var file models.File
	err := s.db.Preload("Owner").Preload("Folder").First(&file, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (s *fileService) Create(name, path, mimeType string, size int64, folderID, teamID, ownerID string) (*models.File, error) {
	file := &models.File{
		Name:     name,
		Path:     path,
		MimeType: mimeType,
		Size:     size,
		OwnerID:  ownerID,
	}
	if folderID != "" {
		file.FolderID = &folderID
	}
	if teamID != "" {
		file.TeamID = &teamID
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(file).Error; err != nil {
			return err
		}
		// Update user quota
		return tx.Model(&models.User{}).Where("id = ?", ownerID).Update("quota_used", gorm.Expr("quota_used + ?", size)).Error
	})

	return file, err
}

func (s *fileService) Rename(id, name string) (*models.File, error) {
	err := s.db.Model(&models.File{}).Where("id = ?", id).Update("name", name).Error
	if err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *fileService) Move(id, folderID, newPath string) (*models.File, error) {
	updates := map[string]interface{}{"path": newPath}
	if folderID != "" {
		updates["folder_id"] = folderID
	} else {
		updates["folder_id"] = nil
	}
	err := s.db.Model(&models.File{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *fileService) Copy(id, ownerID string) (*models.File, error) {
	original, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	copy := &models.File{
		Name:     original.Name + " (copy)",
		Path:     original.Path,
		MimeType: original.MimeType,
		Size:     original.Size,
		FolderID: original.FolderID,
		TeamID:   original.TeamID,
		OwnerID:  ownerID,
	}

	err = s.db.Create(copy).Error
	return copy, err
}

func (s *fileService) ToggleFavorite(id string) (*models.File, error) {
	file, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	err = s.db.Model(&models.File{}).Where("id = ?", id).Update("is_favorite", !file.IsFavorite).Error
	if err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *fileService) Delete(id string) error {
	file, err := s.Get(id)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Update user quota
		if err := tx.Model(&models.User{}).Where("id = ?", file.OwnerID).Update("quota_used", gorm.Expr("quota_used - ?", file.Size)).Error; err != nil {
			return err
		}
		return tx.Delete(&models.File{}, "id = ?", id).Error
	})
}

func (s *fileService) DeleteMany(ids []string, userID string) error {
	var totalSize int64
	s.db.Model(&models.File{}).Where("id IN ? AND owner_id = ?", ids, userID).Select("COALESCE(SUM(size), 0)").Scan(&totalSize)

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.User{}).Where("id = ?", userID).Update("quota_used", gorm.Expr("quota_used - ?", totalSize)).Error; err != nil {
			return err
		}
		return tx.Delete(&models.File{}, "id IN ?", ids).Error
	})
}

func (s *fileService) GetStorageInfo(userID string) (used, total, available int64, usedPercent int) {
	var user models.User
	s.db.Select("quota_used").First(&user, "id = ?", userID)

	used = user.QuotaUsed
	total = 10 * 1024 * 1024 * 1024 // 10GB
	available = total - used
	if total > 0 {
		usedPercent = int(float64(used) / float64(total) * 100)
	}
	return
}

func (s *fileService) GetDownloadURL(id string) string {
	return "/api/files/" + id + "/download"
}

func (s *fileService) IsOwner(fileID, userID string) bool {
	var file models.File
	err := s.db.Select("owner_id").First(&file, "id = ?", fileID).Error
	if err != nil {
		return false
	}
	return file.OwnerID == userID
}
