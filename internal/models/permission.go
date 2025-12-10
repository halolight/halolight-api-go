package models

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID          string         `gorm:"primaryKey;type:char(26)" json:"id"`
	Action      string         `gorm:"uniqueIndex;size:150;not null" json:"action"` // e.g. users:view, documents:edit
	Resource    string         `gorm:"size:100;not null" json:"resource"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Roles []RolePermission `gorm:"foreignKey:PermissionID" json:"roles,omitempty"`
}

func (Permission) TableName() string {
	return "permissions"
}

func (p *Permission) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = GenerateULID()
	}
	return nil
}

// RolePermission is the join table between Role and Permission
type RolePermission struct {
	RoleID       string     `gorm:"primaryKey;type:char(26)" json:"roleId"`
	PermissionID string     `gorm:"primaryKey;type:char(26)" json:"permissionId"`
	CreatedAt    time.Time  `json:"createdAt"`
	Role         Role       `gorm:"constraint:OnDelete:CASCADE" json:"role,omitempty"`
	Permission   Permission `gorm:"constraint:OnDelete:CASCADE" json:"permission,omitempty"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

// UserRole is the join table between User and Role
type UserRole struct {
	UserID    string    `gorm:"primaryKey;type:char(26)" json:"userId"`
	RoleID    string    `gorm:"primaryKey;type:char(26)" json:"roleId"`
	CreatedAt time.Time `json:"createdAt"`
	User      User      `gorm:"constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Role      Role      `gorm:"constraint:OnDelete:CASCADE" json:"role,omitempty"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
