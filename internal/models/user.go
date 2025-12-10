package models

import (
	"time"

	"gorm.io/gorm"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "ACTIVE"
	UserStatusInactive  UserStatus = "INACTIVE"
	UserStatusSuspended UserStatus = "SUSPENDED"
)

type User struct {
	ID          string         `gorm:"primaryKey;type:char(26)" json:"id"`
	Email       string         `gorm:"uniqueIndex;size:191;not null" json:"email"`
	Phone       *string        `gorm:"uniqueIndex;size:50" json:"phone,omitempty"`
	Username    string         `gorm:"uniqueIndex;size:100;not null" json:"username"`
	Password    string         `gorm:"size:255;not null" json:"-"`
	Name        string         `gorm:"size:191;not null" json:"name"`
	Avatar      *string        `gorm:"size:255" json:"avatar,omitempty"`
	Status      UserStatus     `gorm:"type:varchar(20);default:ACTIVE" json:"status"`
	Department  *string        `gorm:"size:191" json:"department,omitempty"`
	Position    *string        `gorm:"size:191" json:"position,omitempty"`
	Bio         *string        `gorm:"type:text" json:"bio,omitempty"`
	QuotaUsed   int64          `gorm:"default:0" json:"quotaUsed"`
	LastLoginAt *time.Time     `json:"lastLoginAt,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Roles         []UserRole                 `gorm:"foreignKey:UserID" json:"roles,omitempty"`
	Teams         []TeamMember               `gorm:"foreignKey:UserID" json:"teams,omitempty"`
	OwnedTeams    []Team                     `gorm:"foreignKey:OwnerID" json:"ownedTeams,omitempty"`
	Documents     []Document                 `gorm:"foreignKey:OwnerID" json:"documents,omitempty"`
	SharedDocs    []DocumentShare            `gorm:"foreignKey:SharedWithID" json:"sharedDocs,omitempty"`
	Files         []File                     `gorm:"foreignKey:OwnerID" json:"files,omitempty"`
	Folders       []Folder                   `gorm:"foreignKey:OwnerID" json:"folders,omitempty"`
	OwnedEvents   []CalendarEvent            `gorm:"foreignKey:OwnerID" json:"ownedEvents,omitempty"`
	EventAttend   []EventAttendee            `gorm:"foreignKey:UserID" json:"eventAttend,omitempty"`
	Notifications []Notification             `gorm:"foreignKey:UserID" json:"notifications,omitempty"`
	Conversations []ConversationParticipant  `gorm:"foreignKey:UserID" json:"conversations,omitempty"`
	Messages      []Message                  `gorm:"foreignKey:SenderID" json:"messages,omitempty"`
	Activities    []ActivityLog              `gorm:"foreignKey:ActorID" json:"activities,omitempty"`
	RefreshTokens []RefreshToken             `gorm:"foreignKey:UserID" json:"-"`
}

func (User) TableName() string {
	return "users"
}

// BeforeCreate hook to generate ID using ulid
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = GenerateULID()
	}
	return nil
}
