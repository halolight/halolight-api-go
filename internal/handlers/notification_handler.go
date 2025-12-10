package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/halolight/halolight-api-go/internal/services"
)

type NotificationHandler struct {
	svc services.NotificationService
}

func NewNotificationHandler(svc services.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) List(c *gin.Context) {
	userID := c.GetString("userID")
	page := getIntQuery(c, "page", 1)
	limit := getIntQuery(c, "limit", 20)
	unreadOnly := getBoolQuery(c, "unreadOnly", false)

	notifications, total, unreadCount, err := h.svc.List(userID, page, limit, unreadOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    notifications,
		"meta": gin.H{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"totalPages":  (total + int64(limit) - 1) / int64(limit),
			"unreadCount": unreadCount,
		},
	})
}

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetString("userID")
	count, err := h.svc.GetUnreadCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"count": count}})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetString("userID")
	notifID := c.Param("id")

	if !h.svc.IsOwner(notifID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Access denied"})
		return
	}

	notification, err := h.svc.MarkAsRead(notifID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": notification, "message": "Marked as read"})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetString("userID")
	if err := h.svc.MarkAllAsRead(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "All marked as read"})
}

func (h *NotificationHandler) Delete(c *gin.Context) {
	userID := c.GetString("userID")
	notifID := c.Param("id")

	if !h.svc.IsOwner(notifID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Access denied"})
		return
	}

	if err := h.svc.Delete(notifID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Notification deleted"})
}
