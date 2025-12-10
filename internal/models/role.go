package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          string         `gorm:"primaryKey;type:char(26)" json:"id"`
	Name        string         `gorm:"uniqueIndex;size:100;not null" json:"name"`
	Label       string         `gorm:"size:191;not null" json:"label"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Permissions []RolePermission `gorm:"foreignKey:RoleID" json:"permissions,omitempty"`
	Users       []UserRole       `gorm:"foreignKey:RoleID" json:"users,omitempty"`
	TeamMembers []TeamMember     `gorm:"foreignKey:RoleID" json:"teamMembers,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = GenerateULID()
	}
	return nil
}
