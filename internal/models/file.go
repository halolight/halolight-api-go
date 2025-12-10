package models

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID         string         `gorm:"primaryKey;type:char(26)" json:"id"`
	Name       string         `gorm:"size:255;not null" json:"name"`
	Path       string         `gorm:"size:500;not null" json:"path"`
	MimeType   string         `gorm:"size:100;not null" json:"mimeType"`
	Size       int64          `gorm:"not null" json:"size"`
	Thumbnail  *string        `gorm:"size:255" json:"thumbnail,omitempty"`
	FolderID   *string        `gorm:"index;type:char(26)" json:"folderId,omitempty"`
	OwnerID    string         `gorm:"index;type:char(26);not null" json:"ownerId"`
	TeamID     *string        `gorm:"index;type:char(26)" json:"teamId,omitempty"`
	IsFavorite bool           `gorm:"default:false" json:"isFavorite"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Owner  User    `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE" json:"owner,omitempty"`
	Team   *Team   `json:"team,omitempty"`
	Folder *Folder `json:"folder,omitempty"`
}

func (File) TableName() string {
	return "files"
}

func (f *File) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = GenerateULID()
	}
	return nil
}

type Folder struct {
	ID        string         `gorm:"primaryKey;type:char(26)" json:"id"`
	Name      string         `gorm:"size:191;not null" json:"name"`
	Path      string         `gorm:"uniqueIndex;size:500;not null" json:"path"`
	ParentID  *string        `gorm:"type:char(26);index" json:"parentId,omitempty"`
	OwnerID   string         `gorm:"index;type:char(26);not null" json:"ownerId"`
	TeamID    *string        `gorm:"type:char(26)" json:"teamId,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Parent   *Folder  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Folder `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Owner    User     `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE" json:"owner,omitempty"`
	Team     *Team    `json:"team,omitempty"`
	Files    []File   `gorm:"foreignKey:FolderID" json:"files,omitempty"`
}

func (Folder) TableName() string {
	return "folders"
}

func (f *Folder) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = GenerateULID()
	}
	return nil
}
