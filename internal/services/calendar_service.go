package services

import (
	"time"

	"github.com/halolight/halolight-api-go/internal/models"
	"gorm.io/gorm"
)

type CalendarService interface {
	List(userID string, startDate, endDate *time.Time) ([]models.CalendarEvent, error)
	Get(id string) (*models.CalendarEvent, error)
	Create(title, description, location string, startAt, endAt time.Time, eventType, color string, allDay bool, ownerID string, attendeeIDs []string) (*models.CalendarEvent, error)
	Update(id, title, description, location, eventType, color string, allDay *bool, startAt, endAt *time.Time) (*models.CalendarEvent, error)
	Reschedule(id string, startAt, endAt time.Time) (*models.CalendarEvent, error)
	AddAttendee(eventID, userID string) (*models.EventAttendee, error)
	RemoveAttendee(eventID, userID string) error
	Delete(id string) error
	DeleteMany(ids []string) error
	IsOwner(eventID, userID string) bool
}

type calendarService struct {
	db *gorm.DB
}

func NewCalendarService(db *gorm.DB) CalendarService {
	return &calendarService{db: db}
}

func (s *calendarService) List(userID string, startDate, endDate *time.Time) ([]models.CalendarEvent, error) {
	var events []models.CalendarEvent

	query := s.db.Model(&models.CalendarEvent{}).
		Where("owner_id = ? OR id IN (SELECT event_id FROM event_attendees WHERE user_id = ?)", userID, userID)

	if startDate != nil {
		query = query.Where("start_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("end_at <= ?", *endDate)
	}

	err := query.Preload("Owner").Preload("Attendees.User").Order("start_at ASC").Find(&events).Error
	return events, err
}

func (s *calendarService) Get(id string) (*models.CalendarEvent, error) {
	var event models.CalendarEvent
	err := s.db.Preload("Owner").Preload("Attendees.User").Preload("Reminders").First(&event, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (s *calendarService) Create(title, description, location string, startAt, endAt time.Time, eventType, color string, allDay bool, ownerID string, attendeeIDs []string) (*models.CalendarEvent, error) {
	event := &models.CalendarEvent{
		Title:       title,
		Description: &description,
		Location:    &location,
		StartAt:     startAt,
		EndAt:       endAt,
		Type:        eventType,
		Color:       &color,
		AllDay:      allDay,
		OwnerID:     ownerID,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(event).Error; err != nil {
			return err
		}

		for _, userID := range attendeeIDs {
			attendee := &models.EventAttendee{
				EventID: event.ID,
				UserID:  userID,
			}
			if err := tx.Create(attendee).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return s.Get(event.ID)
}

func (s *calendarService) Update(id, title, description, location, eventType, color string, allDay *bool, startAt, endAt *time.Time) (*models.CalendarEvent, error) {
	updates := make(map[string]interface{})

	if title != "" {
		updates["title"] = title
	}
	if description != "" {
		updates["description"] = description
	}
	if location != "" {
		updates["location"] = location
	}
	if eventType != "" {
		updates["type"] = eventType
	}
	if color != "" {
		updates["color"] = color
	}
	if allDay != nil {
		updates["all_day"] = *allDay
	}
	if startAt != nil {
		updates["start_at"] = *startAt
	}
	if endAt != nil {
		updates["end_at"] = *endAt
	}

	if len(updates) > 0 {
		if err := s.db.Model(&models.CalendarEvent{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return s.Get(id)
}

func (s *calendarService) Reschedule(id string, startAt, endAt time.Time) (*models.CalendarEvent, error) {
	err := s.db.Model(&models.CalendarEvent{}).Where("id = ?", id).Updates(map[string]interface{}{
		"start_at": startAt,
		"end_at":   endAt,
	}).Error
	if err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *calendarService) AddAttendee(eventID, userID string) (*models.EventAttendee, error) {
	attendee := &models.EventAttendee{
		EventID: eventID,
		UserID:  userID,
	}
	err := s.db.Create(attendee).Error
	if err != nil {
		return nil, err
	}
	s.db.Preload("User").First(attendee, "event_id = ? AND user_id = ?", eventID, userID)
	return attendee, nil
}

func (s *calendarService) RemoveAttendee(eventID, userID string) error {
	return s.db.Where("event_id = ? AND user_id = ?", eventID, userID).Delete(&models.EventAttendee{}).Error
}

func (s *calendarService) Delete(id string) error {
	return s.db.Delete(&models.CalendarEvent{}, "id = ?", id).Error
}

func (s *calendarService) DeleteMany(ids []string) error {
	return s.db.Delete(&models.CalendarEvent{}, "id IN ?", ids).Error
}

func (s *calendarService) IsOwner(eventID, userID string) bool {
	var event models.CalendarEvent
	err := s.db.Select("owner_id").First(&event, "id = ?", eventID).Error
	if err != nil {
		return false
	}
	return event.OwnerID == userID
}
