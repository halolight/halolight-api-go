package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/halolight/halolight-api-go/internal/models"
	"github.com/halolight/halolight-api-go/internal/services"
)

type DocumentHandler struct {
	svc services.DocumentService
}

func NewDocumentHandler(svc services.DocumentService) *DocumentHandler {
	return &DocumentHandler{svc: svc}
}

func (h *DocumentHandler) List(c *gin.Context) {
	userID := c.GetString("userID")
	page := getIntQuery(c, "page", 1)
	limit := getIntQuery(c, "limit", 10)
	search := c.Query("search")
	folder := c.Query("folder")
	tagsStr := c.Query("tags")
	var tags []string
	if tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
	}

	docs, total, err := h.svc.List(userID, page, limit, search, folder, tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    docs,
		"meta": gin.H{
			"total":      total,
			"page":       page,
			"limit":      limit,
			"totalPages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func (h *DocumentHandler) Get(c *gin.Context) {
	userID := c.GetString("userID")
	docID := c.Param("id")

	if !h.svc.HasAccess(docID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Access denied"})
		return
	}

	doc, err := h.svc.Get(docID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Document not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": doc})
}

func (h *DocumentHandler) Create(c *gin.Context) {
	var req struct {
		Title   string  `json:"title" binding:"required,min=1"`
		Content string  `json:"content"`
		Folder  string  `json:"folder"`
		Type    string  `json:"type"`
		TeamID  *string `json:"teamId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	doc, err := h.svc.Create(req.Title, req.Content, req.Folder, req.Type, userID, req.TeamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": doc, "message": "Document created"})
}

func (h *DocumentHandler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	docID := c.Param("id")

	if !h.svc.IsOwner(docID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only owner can update"})
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Folder  string `json:"folder"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	doc, err := h.svc.Update(docID, req.Title, req.Content, req.Folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": doc, "message": "Document updated"})
}

func (h *DocumentHandler) Rename(c *gin.Context) {
	userID := c.GetString("userID")
	docID := c.Param("id")

	if !h.svc.IsOwner(docID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only owner can rename"})
		return
	}

	var req struct {
		Title string `json:"title" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	doc, err := h.svc.Rename(docID, req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": doc, "message": "Document renamed"})
}

func (h *DocumentHandler) Move(c *gin.Context) {
	userID := c.GetString("userID")
	docID := c.Param("id")

	if !h.svc.IsOwner(docID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only owner can move"})
		return
	}

	var req struct {
		Folder string `json:"folder" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	doc, err := h.svc.Move(docID, req.Folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": doc, "message": "Document moved"})
}

func (h *DocumentHandler) UpdateTags(c *gin.Context) {
	userID := c.GetString("userID")
	docID := c.Param("id")

	if !h.svc.IsOwner(docID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only owner can update tags"})
		return
	}

	var req struct {
		Tags []string `json:"tags" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	doc, err := h.svc.UpdateTags(docID, req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": doc, "message": "Tags updated"})
}

func (h *DocumentHandler) Share(c *gin.Context) {
	userID := c.GetString("userID")
	docID := c.Param("id")

	if !h.svc.IsOwner(docID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only owner can share"})
		return
	}

	var req struct {
		UserID     string                  `json:"userId" binding:"required"`
		Permission models.SharePermission `json:"permission"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if req.Permission == "" {
		req.Permission = models.SharePermissionRead
	}

	if err := h.svc.Share(docID, req.UserID, req.Permission); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Document shared"})
}

func (h *DocumentHandler) Unshare(c *gin.Context) {
	userID := c.GetString("userID")
	docID := c.Param("id")

	if !h.svc.IsOwner(docID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only owner can unshare"})
		return
	}

	var req struct {
		UserID string `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if err := h.svc.Unshare(docID, req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Share removed"})
}

func (h *DocumentHandler) BatchDelete(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Documents deleted"})
}

func (h *DocumentHandler) Delete(c *gin.Context) {
	userID := c.GetString("userID")
	docID := c.Param("id")

	if !h.svc.IsOwner(docID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only owner can delete"})
		return
	}

	if err := h.svc.Delete(docID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Document deleted"})
}
