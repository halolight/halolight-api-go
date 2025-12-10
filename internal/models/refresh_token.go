package models

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        string    `gorm:"primaryKey;type:char(26)" json:"id"`
	UserID    string    `gorm:"index;type:char(26);not null" json:"userId"`
	Token     string    `gorm:"uniqueIndex;size:500;not null" json:"token"`
	ExpiresAt time.Time `gorm:"index;not null" json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`

	// Relations
	User User `gorm:"constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

func (rt *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	if rt.ID == "" {
		rt.ID = GenerateULID()
	}
	return nil
}

// IsExpired checks if the refresh token has expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}
