package models

import (
	"time"

	"gorm.io/gorm"
)

type Conversation struct {
	ID        string         `gorm:"primaryKey;type:char(26)" json:"id"`
	Name      *string        `gorm:"size:191" json:"name,omitempty"`
	IsGroup   bool           `gorm:"default:false" json:"isGroup"`
	Avatar    *string        `gorm:"size:255" json:"avatar,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Participants []ConversationParticipant `gorm:"foreignKey:ConversationID" json:"participants,omitempty"`
	Messages     []Message                 `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
}

func (Conversation) TableName() string {
	return "conversations"
}

func (c *Conversation) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = GenerateULID()
	}
	return nil
}

type ConversationParticipant struct {
	ConversationID string     `gorm:"primaryKey;type:char(26)" json:"conversationId"`
	UserID         string     `gorm:"primaryKey;type:char(26)" json:"userId"`
	Role           *string    `gorm:"size:20;default:'member'" json:"role,omitempty"`
	UnreadCount    int        `gorm:"default:0" json:"unreadCount"`
	LastReadAt     *time.Time `json:"lastReadAt,omitempty"`
	JoinedAt       time.Time  `gorm:"autoCreateTime" json:"joinedAt"`

	// Relations
	Conversation Conversation `gorm:"constraint:OnDelete:CASCADE" json:"conversation,omitempty"`
	User         User         `gorm:"constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

func (ConversationParticipant) TableName() string {
	return "conversation_participants"
}

type Message struct {
	ID             string         `gorm:"primaryKey;type:char(26)" json:"id"`
	ConversationID string         `gorm:"index;type:char(26);not null" json:"conversationId"`
	SenderID       string         `gorm:"index;type:char(26);not null" json:"senderId"`
	Type           string         `gorm:"size:20;default:'text'" json:"type"`
	Content        string         `gorm:"type:text;not null" json:"content"`
	CreatedAt      time.Time      `json:"createdAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Conversation Conversation `gorm:"constraint:OnDelete:CASCADE" json:"conversation,omitempty"`
	Sender       User         `gorm:"constraint:OnDelete:CASCADE" json:"sender,omitempty"`
}

func (Message) TableName() string {
	return "messages"
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = GenerateULID()
	}
	return nil
}
