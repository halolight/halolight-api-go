package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/halolight/halolight-api-go/internal/services"
)

type FileHandler struct {
	svc       services.FileService
	folderSvc services.FolderService
}

func NewFileHandler(svc services.FileService, folderSvc services.FolderService) *FileHandler {
	return &FileHandler{svc: svc, folderSvc: folderSvc}
}

func (h *FileHandler) List(c *gin.Context) {
	userID := c.GetString("userID")
	page := getIntQuery(c, "page", 1)
	limit := getIntQuery(c, "limit", 20)
	path := c.Query("path")
	fileType := c.Query("type")
	search := c.Query("search")
	folderID := c.Query("folderId")

	files, total, err := h.svc.List(userID, page, limit, path, fileType, search, folderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    files,
		"meta": gin.H{
			"total":      total,
			"page":       page,
			"limit":      limit,
			"totalPages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func (h *FileHandler) Get(c *gin.Context) {
	file, err := h.svc.Get(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "File not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": file})
}

func (h *FileHandler) GetStorage(c *gin.Context) {
	userID := c.GetString("userID")
	used, total, available, usedPercent := h.svc.GetStorageInfo(userID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"used":        used,
			"total":       total,
			"available":   available,
			"usedPercent": usedPercent,
		},
	})
}

func (h *FileHandler) GetDownloadURL(c *gin.Context) {
	userID := c.GetString("userID")
	fileID := c.Param("id")

	if !h.svc.IsOwner(fileID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Access denied"})
		return
	}

	url := h.svc.GetDownloadURL(fileID)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"url": url}})
}

func (h *FileHandler) Upload(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required,min=1"`
		Path     string `json:"path" binding:"required"`
		MimeType string `json:"mimeType" binding:"required"`
		Size     int64  `json:"size" binding:"required"`
		FolderID string `json:"folderId"`
		TeamID   string `json:"teamId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	file, err := h.svc.Create(req.Name, req.Path, req.MimeType, req.Size, req.FolderID, req.TeamID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": file, "message": "File uploaded"})
}

func (h *FileHandler) CreateFolder(c *gin.Context) {
	var req struct {
		Name     string  `json:"name" binding:"required,min=1"`
		ParentID *string `json:"parentId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	folder, err := h.folderSvc.Create(req.Name, req.ParentID, nil, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": folder, "message": "Folder created"})
}

func (h *FileHandler) Rename(c *gin.Context) {
	userID := c.GetString("userID")
	fileID := c.Param("id")

	if !h.svc.IsOwner(fileID, userID) {
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

	file, err := h.svc.Rename(fileID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": file, "message": "File renamed"})
}

func (h *FileHandler) Move(c *gin.Context) {
	userID := c.GetString("userID")
	fileID := c.Param("id")

	if !h.svc.IsOwner(fileID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only owner can move"})
		return
	}

	var req struct {
		FolderID string `json:"folderId"`
		Path     string `json:"path" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	file, err := h.svc.Move(fileID, req.FolderID, req.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": file, "message": "File moved"})
}

func (h *FileHandler) Copy(c *gin.Context) {
	userID := c.GetString("userID")
	file, err := h.svc.Copy(c.Param("id"), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": file, "message": "File copied"})
}

func (h *FileHandler) ToggleFavorite(c *gin.Context) {
	userID := c.GetString("userID")
	fileID := c.Param("id")

	if !h.svc.IsOwner(fileID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only owner can favorite"})
		return
	}

	file, err := h.svc.ToggleFavorite(fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": file, "message": "Favorite toggled"})
}

func (h *FileHandler) Share(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"shared": true}, "message": "File shared"})
}

func (h *FileHandler) BatchDelete(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	if err := h.svc.DeleteMany(req.IDs, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Files deleted"})
}

func (h *FileHandler) Delete(c *gin.Context) {
	userID := c.GetString("userID")
	fileID := c.Param("id")

	if !h.svc.IsOwner(fileID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only owner can delete"})
		return
	}

	if err := h.svc.Delete(fileID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "File deleted"})
}
