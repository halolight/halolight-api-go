package models

import (
	"time"

	"gorm.io/gorm"
)

type SharePermission string

const (
	SharePermissionRead    SharePermission = "READ"
	SharePermissionEdit    SharePermission = "EDIT"
	SharePermissionComment SharePermission = "COMMENT"
)

type Document struct {
	ID        string         `gorm:"primaryKey;type:char(26)" json:"id"`
	Title     string         `gorm:"size:255;not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	Folder    *string        `gorm:"index;size:255" json:"folder,omitempty"`
	Type      string         `gorm:"size:50;not null" json:"type"`
	Size      int64          `gorm:"default:0" json:"size"`
	Views     int            `gorm:"default:0" json:"views"`
	OwnerID   string         `gorm:"index;type:char(26);not null" json:"ownerId"`
	TeamID    *string        `gorm:"index;type:char(26)" json:"teamId,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Owner  User            `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE" json:"owner,omitempty"`
	Team   *Team           `json:"team,omitempty"`
	Shares []DocumentShare `gorm:"foreignKey:DocumentID" json:"shares,omitempty"`
	Tags   []DocumentTag   `gorm:"foreignKey:DocumentID" json:"tags,omitempty"`
}

func (Document) TableName() string {
	return "documents"
}

func (d *Document) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = GenerateULID()
	}
	return nil
}

type DocumentShare struct {
	ID           string          `gorm:"primaryKey;type:char(26)" json:"id"`
	DocumentID   string          `gorm:"index;type:char(26);not null" json:"documentId"`
	SharedWithID *string         `gorm:"index;type:char(26)" json:"sharedWithId,omitempty"`
	TeamID       *string         `gorm:"index;type:char(26)" json:"teamId,omitempty"`
	Permission   SharePermission `gorm:"type:varchar(10);default:READ" json:"permission"`
	ExpiresAt    *time.Time      `json:"expiresAt,omitempty"`
	CreatedAt    time.Time       `json:"createdAt"`

	// Relations
	Document   Document `gorm:"constraint:OnDelete:CASCADE" json:"document,omitempty"`
	SharedWith *User    `json:"sharedWith,omitempty"`
	Team       *Team    `json:"team,omitempty"`
}

func (DocumentShare) TableName() string {
	return "document_shares"
}

func (ds *DocumentShare) BeforeCreate(tx *gorm.DB) error {
	if ds.ID == "" {
		ds.ID = GenerateULID()
	}
	return nil
}

type Tag struct {
	ID        string         `gorm:"primaryKey;type:char(26)" json:"id"`
	Name      string         `gorm:"uniqueIndex;size:100;not null" json:"name"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Documents []DocumentTag `gorm:"foreignKey:TagID" json:"documents,omitempty"`
}

func (Tag) TableName() string {
	return "tags"
}

func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = GenerateULID()
	}
	return nil
}

type DocumentTag struct {
	DocumentID string    `gorm:"primaryKey;type:char(26)" json:"documentId"`
	TagID      string    `gorm:"primaryKey;type:char(26)" json:"tagId"`
	CreatedAt  time.Time `json:"createdAt"`
	Document   Document  `gorm:"constraint:OnDelete:CASCADE" json:"document,omitempty"`
	Tag        Tag       `gorm:"constraint:OnDelete:CASCADE" json:"tag,omitempty"`
}

func (DocumentTag) TableName() string {
	return "document_tags"
}
