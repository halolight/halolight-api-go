package models

import (
	"time"

	"gorm.io/gorm"
)

type AttendeeStatus string

const (
	AttendeePending  AttendeeStatus = "PENDING"
	AttendeeAccepted AttendeeStatus = "ACCEPTED"
	AttendeeDeclined AttendeeStatus = "DECLINED"
)

type CalendarEvent struct {
	ID          string         `gorm:"primaryKey;type:char(26)" json:"id"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	StartAt     time.Time      `gorm:"index;not null" json:"startAt"`
	EndAt       time.Time      `gorm:"index;not null" json:"endAt"`
	Type        string         `gorm:"size:50;not null" json:"type"`
	Color       *string        `gorm:"size:50" json:"color,omitempty"`
	AllDay      bool           `gorm:"default:false" json:"allDay"`
	Location    *string        `gorm:"size:255" json:"location,omitempty"`
	OwnerID     string         `gorm:"index;type:char(26);not null" json:"ownerId"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Owner     User            `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE" json:"owner,omitempty"`
	Attendees []EventAttendee `gorm:"foreignKey:EventID" json:"attendees,omitempty"`
	Reminders []EventReminder `gorm:"foreignKey:EventID" json:"reminders,omitempty"`
}

func (CalendarEvent) TableName() string {
	return "calendar_events"
}

func (ce *CalendarEvent) BeforeCreate(tx *gorm.DB) error {
	if ce.ID == "" {
		ce.ID = GenerateULID()
	}
	return nil
}

type EventAttendee struct {
	EventID   string         `gorm:"primaryKey;type:char(26)" json:"eventId"`
	UserID    string         `gorm:"primaryKey;type:char(26)" json:"userId"`
	Status    AttendeeStatus `gorm:"type:varchar(20);default:PENDING" json:"status"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`

	// Relations
	Event CalendarEvent `gorm:"constraint:OnDelete:CASCADE" json:"event,omitempty"`
	User  User          `gorm:"constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

func (EventAttendee) TableName() string {
	return "event_attendees"
}

type EventReminder struct {
	ID        string        `gorm:"primaryKey;type:char(26)" json:"id"`
	EventID   string        `gorm:"index;type:char(26);not null" json:"eventId"`
	RemindAt  time.Time     `gorm:"not null" json:"remindAt"`
	CreatedAt time.Time     `json:"createdAt"`

	// Relations
	Event CalendarEvent `gorm:"constraint:OnDelete:CASCADE" json:"event,omitempty"`
}

func (EventReminder) TableName() string {
	return "event_reminders"
}

func (er *EventReminder) BeforeCreate(tx *gorm.DB) error {
	if er.ID == "" {
		er.ID = GenerateULID()
	}
	return nil
}
