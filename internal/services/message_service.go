package services

import (
	"time"

	"github.com/halolight/halolight-api-go/internal/models"
	"gorm.io/gorm"
)

type MessageService interface {
	GetConversations(userID string) ([]ConversationWithLastMessage, error)
	GetConversation(id, userID string) (*models.Conversation, error)
	CreateConversation(name string, isGroup bool, participantIDs []string, creatorID string) (*models.Conversation, error)
	SendMessage(conversationID, senderID, content, msgType string) (*models.Message, error)
	MarkAsRead(conversationID, userID string) error
	DeleteMessage(id string) error
	IsParticipant(conversationID, userID string) bool
	IsMessageOwner(messageID, userID string) bool
}

type ConversationWithLastMessage struct {
	models.Conversation
	LastMessage *models.Message `json:"lastMessage"`
	UnreadCount int             `json:"unreadCount"`
}

type messageService struct {
	db *gorm.DB
}

func NewMessageService(db *gorm.DB) MessageService {
	return &messageService{db: db}
}

func (s *messageService) GetConversations(userID string) ([]ConversationWithLastMessage, error) {
	var conversations []models.Conversation
	err := s.db.
		Joins("JOIN conversation_participants ON conversations.id = conversation_participants.conversation_id").
		Where("conversation_participants.user_id = ?", userID).
		Preload("Participants.User").
		Order("updated_at DESC").
		Find(&conversations).Error

	if err != nil {
		return nil, err
	}

	result := make([]ConversationWithLastMessage, len(conversations))
	for i, conv := range conversations {
		result[i].Conversation = conv

		// Get last message
		var lastMsg models.Message
		if err := s.db.Where("conversation_id = ?", conv.ID).Order("created_at DESC").First(&lastMsg).Error; err == nil {
			s.db.Preload("Sender").First(&lastMsg, "id = ?", lastMsg.ID)
			result[i].LastMessage = &lastMsg
		}

		// Get unread count
		var participant models.ConversationParticipant
		if err := s.db.Where("conversation_id = ? AND user_id = ?", conv.ID, userID).First(&participant).Error; err == nil {
			result[i].UnreadCount = participant.UnreadCount
		}
	}

	return result, nil
}

func (s *messageService) GetConversation(id, userID string) (*models.Conversation, error) {
	var conversation models.Conversation
	err := s.db.Preload("Participants.User").Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC").Limit(50)
	}).Preload("Messages.Sender").First(&conversation, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	// Mark as read
	now := time.Now()
	s.db.Model(&models.ConversationParticipant{}).
		Where("conversation_id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{"unread_count": 0, "last_read_at": now})

	// Reverse messages to chronological order
	for i, j := 0, len(conversation.Messages)-1; i < j; i, j = i+1, j-1 {
		conversation.Messages[i], conversation.Messages[j] = conversation.Messages[j], conversation.Messages[i]
	}

	return &conversation, nil
}

func (s *messageService) CreateConversation(name string, isGroup bool, participantIDs []string, creatorID string) (*models.Conversation, error) {
	// Ensure creator is in participants
	allParticipants := make(map[string]bool)
	allParticipants[creatorID] = true
	for _, id := range participantIDs {
		allParticipants[id] = true
	}

	conversation := &models.Conversation{
		IsGroup: isGroup || len(allParticipants) > 2,
	}
	if name != "" {
		conversation.Name = &name
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(conversation).Error; err != nil {
			return err
		}

		for userID := range allParticipants {
			role := "member"
			if userID == creatorID {
				role = "owner"
			}
			participant := &models.ConversationParticipant{
				ConversationID: conversation.ID,
				UserID:         userID,
				Role:           &role,
			}
			if err := tx.Create(participant).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	s.db.Preload("Participants.User").First(conversation, "id = ?", conversation.ID)
	return conversation, nil
}

func (s *messageService) SendMessage(conversationID, senderID, content, msgType string) (*models.Message, error) {
	if msgType == "" {
		msgType = "text"
	}

	message := &models.Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
		Type:           msgType,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(message).Error; err != nil {
			return err
		}

		// Update conversation timestamp
		if err := tx.Model(&models.Conversation{}).Where("id = ?", conversationID).Update("updated_at", time.Now()).Error; err != nil {
			return err
		}

		// Increment unread count for other participants
		return tx.Model(&models.ConversationParticipant{}).
			Where("conversation_id = ? AND user_id != ?", conversationID, senderID).
			Update("unread_count", gorm.Expr("unread_count + 1")).Error
	})

	if err != nil {
		return nil, err
	}

	s.db.Preload("Sender").First(message, "id = ?", message.ID)
	return message, nil
}

func (s *messageService) MarkAsRead(conversationID, userID string) error {
	now := time.Now()
	return s.db.Model(&models.ConversationParticipant{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Updates(map[string]interface{}{"unread_count": 0, "last_read_at": now}).Error
}

func (s *messageService) DeleteMessage(id string) error {
	return s.db.Delete(&models.Message{}, "id = ?", id).Error
}

func (s *messageService) IsParticipant(conversationID, userID string) bool {
	var count int64
	s.db.Model(&models.ConversationParticipant{}).Where("conversation_id = ? AND user_id = ?", conversationID, userID).Count(&count)
	return count > 0
}

func (s *messageService) IsMessageOwner(messageID, userID string) bool {
	var message models.Message
	err := s.db.Select("sender_id").First(&message, "id = ?", messageID).Error
	if err != nil {
		return false
	}
	return message.SenderID == userID
}
