package services

import (
	"time"

	"github.com/halolight/halolight-api-go/internal/models"
	"gorm.io/gorm"
)

type NotificationService interface {
	List(userID string, page, limit int, unreadOnly bool) ([]models.Notification, int64, int64, error)
	GetUnreadCount(userID string) (int64, error)
	MarkAsRead(id string) (*models.Notification, error)
	MarkAllAsRead(userID string) error
	Create(userID, notifType, title, content string, link *string, payload interface{}) (*models.Notification, error)
	Delete(id string) error
	IsOwner(notifID, userID string) bool
}

type notificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) NotificationService {
	return &notificationService{db: db}
}

func (s *notificationService) List(userID string, page, limit int, unreadOnly bool) ([]models.Notification, int64, int64, error) {
	var notifications []models.Notification
	var total, unreadCount int64

	query := s.db.Model(&models.Notification{}).Where("user_id = ?", userID)

	if unreadOnly {
		query = query.Where("read = ?", false)
	}

	query.Count(&total)
	s.db.Model(&models.Notification{}).Where("user_id = ? AND read = ?", userID, false).Count(&unreadCount)

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&notifications).Error

	return notifications, total, unreadCount, err
}

func (s *notificationService) GetUnreadCount(userID string) (int64, error) {
	var count int64
	err := s.db.Model(&models.Notification{}).Where("user_id = ? AND read = ?", userID, false).Count(&count).Error
	return count, err
}

func (s *notificationService) MarkAsRead(id string) (*models.Notification, error) {
	now := time.Now()
	err := s.db.Model(&models.Notification{}).Where("id = ?", id).Updates(map[string]interface{}{
		"read":    true,
		"read_at": now,
	}).Error
	if err != nil {
		return nil, err
	}

	var notification models.Notification
	s.db.First(&notification, "id = ?", id)
	return &notification, nil
}

func (s *notificationService) MarkAllAsRead(userID string) error {
	now := time.Now()
	return s.db.Model(&models.Notification{}).Where("user_id = ? AND read = ?", userID, false).Updates(map[string]interface{}{
		"read":    true,
		"read_at": now,
	}).Error
}

func (s *notificationService) Create(userID, notifType, title, content string, link *string, payload interface{}) (*models.Notification, error) {
	notification := &models.Notification{
		UserID:  userID,
		Type:    notifType,
		Title:   title,
		Content: content,
		Link:    link,
	}
	err := s.db.Create(notification).Error
	return notification, err
}

func (s *notificationService) Delete(id string) error {
	return s.db.Delete(&models.Notification{}, "id = ?", id).Error
}

func (s *notificationService) IsOwner(notifID, userID string) bool {
	var notification models.Notification
	err := s.db.Select("user_id").First(&notification, "id = ?", notifID).Error
	if err != nil {
		return false
	}
	return notification.UserID == userID
}
