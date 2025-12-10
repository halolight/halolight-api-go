package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/halolight/halolight-api-go/internal/services"
)

type MessageHandler struct {
	svc services.MessageService
}

func NewMessageHandler(svc services.MessageService) *MessageHandler {
	return &MessageHandler{svc: svc}
}

func (h *MessageHandler) GetConversations(c *gin.Context) {
	userID := c.GetString("userID")
	conversations, err := h.svc.GetConversations(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": conversations})
}

func (h *MessageHandler) GetConversation(c *gin.Context) {
	userID := c.GetString("userID")
	convID := c.Param("id")

	if !h.svc.IsParticipant(convID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Access denied"})
		return
	}

	conversation, err := h.svc.GetConversation(convID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Conversation not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": conversation})
}

func (h *MessageHandler) CreateConversation(c *gin.Context) {
	var req struct {
		Name           string   `json:"name"`
		IsGroup        bool     `json:"isGroup"`
		ParticipantIDs []string `json:"participantIds" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	conversation, err := h.svc.CreateConversation(req.Name, req.IsGroup, req.ParticipantIDs, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": conversation, "message": "Conversation created"})
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	var req struct {
		ConversationID string `json:"conversationId" binding:"required"`
		Content        string `json:"content" binding:"required,min=1"`
		Type           string `json:"type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	if !h.svc.IsParticipant(req.ConversationID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Access denied"})
		return
	}

	message, err := h.svc.SendMessage(req.ConversationID, userID, req.Content, req.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": message, "message": "Message sent"})
}

func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetString("userID")
	convID := c.Param("id")

	if err := h.svc.MarkAsRead(convID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Marked as read"})
}

func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	userID := c.GetString("userID")
	msgID := c.Param("id")

	if !h.svc.IsMessageOwner(msgID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only sender can delete"})
		return
	}

	if err := h.svc.DeleteMessage(msgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Message deleted"})
}
