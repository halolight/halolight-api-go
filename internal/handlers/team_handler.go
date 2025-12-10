package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/halolight/halolight-api-go/internal/services"
)

type TeamHandler struct {
	svc services.TeamService
}

func NewTeamHandler(svc services.TeamService) *TeamHandler {
	return &TeamHandler{svc: svc}
}

func (h *TeamHandler) List(c *gin.Context) {
	userID := c.GetString("userID")
	page := getIntQuery(c, "page", 1)
	limit := getIntQuery(c, "limit", 10)
	search := c.Query("search")

	teams, total, err := h.svc.List(userID, page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    teams,
		"meta": gin.H{
			"total":      total,
			"page":       page,
			"limit":      limit,
			"totalPages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func (h *TeamHandler) Get(c *gin.Context) {
	team, err := h.svc.Get(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Team not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": team})
}

func (h *TeamHandler) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required,min=2"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	team, err := h.svc.Create(req.Name, req.Description, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": team, "message": "Team created"})
}

func (h *TeamHandler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	teamID := c.Param("id")

	if !h.svc.IsOwner(teamID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only team owner can update"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	team, err := h.svc.Update(teamID, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": team, "message": "Team updated"})
}

func (h *TeamHandler) Delete(c *gin.Context) {
	userID := c.GetString("userID")
	teamID := c.Param("id")

	if !h.svc.IsOwner(teamID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only team owner can delete"})
		return
	}

	if err := h.svc.Delete(teamID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Team deleted"})
}

func (h *TeamHandler) AddMember(c *gin.Context) {
	userID := c.GetString("userID")
	teamID := c.Param("id")

	if !h.svc.IsOwner(teamID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only team owner can add members"})
		return
	}

	var req struct {
		UserID string  `json:"userId" binding:"required"`
		RoleID *string `json:"roleId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	member, err := h.svc.AddMember(teamID, req.UserID, req.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": member, "message": "Member added"})
}

func (h *TeamHandler) RemoveMember(c *gin.Context) {
	userID := c.GetString("userID")
	teamID := c.Param("id")

	if !h.svc.IsOwner(teamID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Only team owner can remove members"})
		return
	}

	if err := h.svc.RemoveMember(teamID, c.Param("userId")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Member removed"})
}
