package services

import (
	"github.com/halolight/halolight-api-go/internal/models"
	"gorm.io/gorm"
)

type DocumentService interface {
	List(userID string, page, limit int, search, folder string, tags []string) ([]models.Document, int64, error)
	Get(id string) (*models.Document, error)
	Create(title, content, folder, docType, ownerID string, teamID *string) (*models.Document, error)
	Update(id, title, content, folder string) (*models.Document, error)
	Rename(id, title string) (*models.Document, error)
	Move(id, folder string) (*models.Document, error)
	UpdateTags(id string, tagNames []string) (*models.Document, error)
	Share(docID, userID string, permission models.SharePermission) error
	Unshare(docID, userID string) error
	Delete(id string) error
	DeleteMany(ids []string) error
	IsOwner(docID, userID string) bool
	HasAccess(docID, userID string) bool
}

type documentService struct {
	db *gorm.DB
}

func NewDocumentService(db *gorm.DB) DocumentService {
	return &documentService{db: db}
}

func (s *documentService) List(userID string, page, limit int, search, folder string, tags []string) ([]models.Document, int64, error) {
	var docs []models.Document
	var total int64

	query := s.db.Model(&models.Document{}).
		Where("owner_id = ? OR id IN (SELECT document_id FROM document_shares WHERE shared_with_id = ?)", userID, userID)

	if search != "" {
		query = query.Where("title ILIKE ? OR content ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if folder != "" {
		query = query.Where("folder = ?", folder)
	}
	if len(tags) > 0 {
		query = query.Where("id IN (SELECT document_id FROM document_tags dt JOIN tags t ON dt.tag_id = t.id WHERE t.name IN ?)", tags)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Preload("Owner").Preload("Tags.Tag").Offset(offset).Limit(limit).Order("updated_at DESC").Find(&docs).Error

	return docs, total, err
}

func (s *documentService) Get(id string) (*models.Document, error) {
	var doc models.Document
	// Increment views
	s.db.Model(&models.Document{}).Where("id = ?", id).Update("views", gorm.Expr("views + 1"))

	err := s.db.Preload("Owner").Preload("Team").Preload("Tags.Tag").Preload("Shares.SharedWith").First(&doc, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (s *documentService) Create(title, content, folder, docType, ownerID string, teamID *string) (*models.Document, error) {
	doc := &models.Document{
		Title:   title,
		Content: content,
		Folder:  &folder,
		Type:    docType,
		OwnerID: ownerID,
		TeamID:  teamID,
	}
	err := s.db.Create(doc).Error
	return doc, err
}

func (s *documentService) Update(id, title, content, folder string) (*models.Document, error) {
	doc, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	if title != "" {
		doc.Title = title
	}
	if content != "" {
		doc.Content = content
	}
	if folder != "" {
		doc.Folder = &folder
	}

	err = s.db.Save(doc).Error
	return doc, err
}

func (s *documentService) Rename(id, title string) (*models.Document, error) {
	return s.Update(id, title, "", "")
}

func (s *documentService) Move(id, folder string) (*models.Document, error) {
	return s.Update(id, "", "", folder)
}

func (s *documentService) UpdateTags(id string, tagNames []string) (*models.Document, error) {
	// Remove existing tags
	s.db.Where("document_id = ?", id).Delete(&models.DocumentTag{})

	// Add new tags
	for _, name := range tagNames {
		var tag models.Tag
		s.db.FirstOrCreate(&tag, models.Tag{Name: name})
		s.db.Create(&models.DocumentTag{DocumentID: id, TagID: tag.ID})
	}

	return s.Get(id)
}

func (s *documentService) Share(docID, userID string, permission models.SharePermission) error {
	share := &models.DocumentShare{
		DocumentID:   docID,
		SharedWithID: &userID,
		Permission:   permission,
	}
	return s.db.Save(share).Error
}

func (s *documentService) Unshare(docID, userID string) error {
	return s.db.Where("document_id = ? AND shared_with_id = ?", docID, userID).Delete(&models.DocumentShare{}).Error
}

func (s *documentService) Delete(id string) error {
	return s.db.Delete(&models.Document{}, "id = ?", id).Error
}

func (s *documentService) DeleteMany(ids []string) error {
	return s.db.Delete(&models.Document{}, "id IN ?", ids).Error
}

func (s *documentService) IsOwner(docID, userID string) bool {
	var doc models.Document
	err := s.db.Select("owner_id").First(&doc, "id = ?", docID).Error
	if err != nil {
		return false
	}
	return doc.OwnerID == userID
}

func (s *documentService) HasAccess(docID, userID string) bool {
	var count int64
	s.db.Model(&models.Document{}).
		Where("id = ? AND (owner_id = ? OR id IN (SELECT document_id FROM document_shares WHERE shared_with_id = ?))", docID, userID, userID).
		Count(&count)
	return count > 0
}
