package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/halolight/halolight-api-go/internal/services"
)

type FolderHandler struct {
	svc services.FolderService
}

func NewFolderHandler(svc services.FolderService) *FolderHandler {
	return &FolderHandler{svc: svc}
}

func (h *FolderHandler) List(c *gin.Context) {
	userID := c.GetString("userID")
	var parentID *string
	if p := c.Query("parentId"); p != "" {
		parentID = &p
	}

	folders, err := h.svc.List(userID, parentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": folders})
}

func (h *FolderHandler) GetTree(c *gin.Context) {
	userID := c.GetString("userID")
	tree, err := h.svc.GetTree(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": tree})
}

func (h *FolderHandler) Get(c *gin.Context) {
	folder, err := h.svc.Get(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Folder not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": folder})
}

func (h *FolderHandler) Create(c *gin.Context) {
	var req struct {
		Name     string  `json:"name" binding:"required,min=1"`
		ParentID *string `json:"parentId"`
		TeamID   *string `json:"teamId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	folder, err := h.svc.Create(req.Name, req.ParentID, req.TeamID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": folder, "message": "Folder created"})
}

func (h *FolderHandler) Rename(c *gin.Context) {
	userID := c.GetString("userID")
	folderID := c.Param("id")

	if !h.svc.IsOwner(folderID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only owner can rename"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	folder, err := h.svc.Rename(folderID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": folder, "message": "Folder renamed"})
}

func (h *FolderHandler) Delete(c *gin.Context) {
	userID := c.GetString("userID")
	folderID := c.Param("id")

	if !h.svc.IsOwner(folderID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only owner can delete"})
		return
	}

	if err := h.svc.Delete(folderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Folder deleted"})
}
