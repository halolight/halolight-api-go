package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ActivityLog struct {
	ID         string         `gorm:"primaryKey;type:char(26)" json:"id"`
	ActorID    string         `gorm:"index;type:char(26);not null" json:"actorId"`
	Action     string         `gorm:"size:100;not null" json:"action"`
	TargetType string         `gorm:"size:50;not null" json:"targetType"`
	TargetID   string         `gorm:"size:50;not null" json:"targetId"`
	Metadata   datatypes.JSON `json:"metadata,omitempty"`
	CreatedAt  time.Time      `json:"createdAt"`

	// Relations
	Actor User `gorm:"constraint:OnDelete:CASCADE" json:"actor,omitempty"`
}

func (ActivityLog) TableName() string {
	return "activity_logs"
}

func (al *ActivityLog) BeforeCreate(tx *gorm.DB) error {
	if al.ID == "" {
		al.ID = GenerateULID()
	}
	return nil
}
