package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/halolight/halolight-api-go/internal/services"
)

type RoleHandler struct {
	svc services.RoleService
}

func NewRoleHandler(svc services.RoleService) *RoleHandler {
	return &RoleHandler{svc: svc}
}

func (h *RoleHandler) List(c *gin.Context) {
	roles, err := h.svc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": roles})
}

func (h *RoleHandler) Get(c *gin.Context) {
	role, err := h.svc.Get(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Role not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": role})
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required,min=2"`
		Label       string `json:"label" binding:"required,min=2"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	role, err := h.svc.Create(req.Name, req.Label, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": role, "message": "Role created"})
}

func (h *RoleHandler) Update(c *gin.Context) {
	var req struct {
		Name        string `json:"name"`
		Label       string `json:"label"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	role, err := h.svc.Update(c.Param("id"), req.Name, req.Label, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": role, "message": "Role updated"})
}

func (h *RoleHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Role deleted"})
}

func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	var req struct {
		PermissionIDs []string `json:"permissionIds" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	role, err := h.svc.AssignPermissions(c.Param("id"), req.PermissionIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": role, "message": "Permissions assigned"})
}
