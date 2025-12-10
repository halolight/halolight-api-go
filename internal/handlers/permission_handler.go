package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/halolight/halolight-api-go/internal/services"
)

type PermissionHandler struct {
	svc services.PermissionService
}

func NewPermissionHandler(svc services.PermissionService) *PermissionHandler {
	return &PermissionHandler{svc: svc}
}

func (h *PermissionHandler) List(c *gin.Context) {
	permissions, err := h.svc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": permissions})
}

func (h *PermissionHandler) Get(c *gin.Context) {
	permission, err := h.svc.Get(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Permission not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": permission})
}

func (h *PermissionHandler) Create(c *gin.Context) {
	var req struct {
		Action      string `json:"action" binding:"required,min=1"`
		Resource    string `json:"resource" binding:"required,min=1"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	permission, err := h.svc.Create(req.Action, req.Resource, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": permission, "message": "Permission created"})
}

func (h *PermissionHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Permission deleted"})
}
