package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Notification struct {
	ID        string         `gorm:"primaryKey;type:char(26)" json:"id"`
	UserID    string         `gorm:"index;type:char(26);not null" json:"userId"`
	Type      string         `gorm:"size:50;not null" json:"type"` // system, user, message, task, alert
	Title     string         `gorm:"size:255;not null" json:"title"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Link      *string        `gorm:"size:255" json:"link,omitempty"`
	Payload   datatypes.JSON `json:"payload,omitempty"`
	Read      bool           `gorm:"default:false" json:"read"`
	ReadAt    *time.Time     `json:"readAt,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User User `gorm:"constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

func (Notification) TableName() string {
	return "notifications"
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = GenerateULID()
	}
	return nil
}
