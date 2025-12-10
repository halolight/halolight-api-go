package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/halolight/halolight-api-go/internal/services"
)

type CalendarHandler struct {
	svc services.CalendarService
}

func NewCalendarHandler(svc services.CalendarService) *CalendarHandler {
	return &CalendarHandler{svc: svc}
}

func (h *CalendarHandler) List(c *gin.Context) {
	userID := c.GetString("userID")
	var startDate, endDate *time.Time

	if s := c.Query("startDate"); s != "" {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			startDate = &t
		}
	}
	if e := c.Query("endDate"); e != "" {
		if t, err := time.Parse(time.RFC3339, e); err == nil {
			endDate = &t
		}
	}

	events, err := h.svc.List(userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": events})
}

func (h *CalendarHandler) Get(c *gin.Context) {
	event, err := h.svc.Get(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Event not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": event})
}

func (h *CalendarHandler) Create(c *gin.Context) {
	var req struct {
		Title       string   `json:"title" binding:"required,min=1"`
		Description string   `json:"description"`
		StartAt     string   `json:"startAt" binding:"required"`
		EndAt       string   `json:"endAt" binding:"required"`
		Type        string   `json:"type"`
		Color       string   `json:"color"`
		AllDay      bool     `json:"allDay"`
		Location    string   `json:"location"`
		AttendeeIDs []string `json:"attendeeIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	startAt, _ := time.Parse(time.RFC3339, req.StartAt)
	endAt, _ := time.Parse(time.RFC3339, req.EndAt)
	userID := c.GetString("userID")

	event, err := h.svc.Create(req.Title, req.Description, req.Location, startAt, endAt, req.Type, req.Color, req.AllDay, userID, req.AttendeeIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": event, "message": "Event created"})
}

func (h *CalendarHandler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	eventID := c.Param("id")

	if !h.svc.IsOwner(eventID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only owner can update"})
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		StartAt     string `json:"startAt"`
		EndAt       string `json:"endAt"`
		Type        string `json:"type"`
		Color       string `json:"color"`
		AllDay      *bool  `json:"allDay"`
		Location    string `json:"location"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	var startAt, endAt *time.Time
	if req.StartAt != "" {
		t, _ := time.Parse(time.RFC3339, req.StartAt)
		startAt = &t
	}
	if req.EndAt != "" {
		t, _ := time.Parse(time.RFC3339, req.EndAt)
		endAt = &t
	}

	event, err := h.svc.Update(eventID, req.Title, req.Description, req.Location, req.Type, req.Color, req.AllDay, startAt, endAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": event, "message": "Event updated"})
}

func (h *CalendarHandler) Reschedule(c *gin.Context) {
	userID := c.GetString("userID")
	eventID := c.Param("id")

	if !h.svc.IsOwner(eventID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only owner can reschedule"})
		return
	}

	var req struct {
		StartAt string `json:"startAt" binding:"required"`
		EndAt   string `json:"endAt" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	startAt, _ := time.Parse(time.RFC3339, req.StartAt)
	endAt, _ := time.Parse(time.RFC3339, req.EndAt)

	event, err := h.svc.Reschedule(eventID, startAt, endAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": event, "message": "Event rescheduled"})
}

func (h *CalendarHandler) AddAttendee(c *gin.Context) {
	userID := c.GetString("userID")
	eventID := c.Param("id")

	if !h.svc.IsOwner(eventID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only owner can add attendees"})
		return
	}

	var req struct {
		UserID string `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	attendee, err := h.svc.AddAttendee(eventID, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": attendee, "message": "Attendee added"})
}

func (h *CalendarHandler) RemoveAttendee(c *gin.Context) {
	userID := c.GetString("userID")
	eventID := c.Param("id")

	if !h.svc.IsOwner(eventID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only owner can remove attendees"})
		return
	}

	if err := h.svc.RemoveAttendee(eventID, c.Param("attendeeId")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Attendee removed"})
}

func (h *CalendarHandler) BatchDelete(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if err := h.svc.DeleteMany(req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Events deleted"})
}

func (h *CalendarHandler) Delete(c *gin.Context) {
	userID := c.GetString("userID")
	eventID := c.Param("id")

	if !h.svc.IsOwner(eventID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only owner can delete"})
		return
	}

	if err := h.svc.Delete(eventID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Event deleted"})
}
