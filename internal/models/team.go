package models

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	ID          string         `gorm:"primaryKey;type:char(26)" json:"id"`
	Name        string         `gorm:"size:191;not null" json:"name"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	Avatar      *string        `gorm:"size:255" json:"avatar,omitempty"`
	OwnerID     string         `gorm:"index;type:char(26);not null" json:"ownerId"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Owner     User            `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE" json:"owner,omitempty"`
	Members   []TeamMember    `gorm:"foreignKey:TeamID" json:"members,omitempty"`
	Documents []Document      `gorm:"foreignKey:TeamID" json:"documents,omitempty"`
	Files     []File          `gorm:"foreignKey:TeamID" json:"files,omitempty"`
	Folders   []Folder        `gorm:"foreignKey:TeamID" json:"folders,omitempty"`
	Shares    []DocumentShare `gorm:"foreignKey:TeamID" json:"shares,omitempty"`
}

func (Team) TableName() string {
	return "teams"
}

func (t *Team) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = GenerateULID()
	}
	return nil
}

type TeamMember struct {
	TeamID   string     `gorm:"primaryKey;type:char(26)" json:"teamId"`
	UserID   string     `gorm:"primaryKey;type:char(26)" json:"userId"`
	RoleID   *string    `gorm:"type:char(26)" json:"roleId,omitempty"`
	JoinedAt time.Time  `gorm:"autoCreateTime" json:"joinedAt"`
	Team     Team       `gorm:"constraint:OnDelete:CASCADE" json:"team,omitempty"`
	User     User       `gorm:"constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Role     *Role      `json:"role,omitempty"`
}

func (TeamMember) TableName() string {
	return "team_members"
}
