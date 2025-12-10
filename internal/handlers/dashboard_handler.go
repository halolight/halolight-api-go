package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/halolight/halolight-api-go/internal/services"
)

type DashboardHandler struct {
	svc services.DashboardService
}

func NewDashboardHandler(svc services.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	stats, err := h.svc.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}

func (h *DashboardHandler) GetVisits(c *gin.Context) {
	visits := h.svc.GetVisits()
	c.JSON(http.StatusOK, gin.H{"success": true, "data": visits})
}

func (h *DashboardHandler) GetSales(c *gin.Context) {
	sales := h.svc.GetSales()
	c.JSON(http.StatusOK, gin.H{"success": true, "data": sales})
}

func (h *DashboardHandler) GetProducts(c *gin.Context) {
	products := h.svc.GetProducts()
	c.JSON(http.StatusOK, gin.H{"success": true, "data": products})
}

func (h *DashboardHandler) GetOrders(c *gin.Context) {
	orders := h.svc.GetOrders()
	c.JSON(http.StatusOK, gin.H{"success": true, "data": orders})
}

func (h *DashboardHandler) GetActivities(c *gin.Context) {
	activities, err := h.svc.GetActivities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": activities})
}

func (h *DashboardHandler) GetPieData(c *gin.Context) {
	pieData := h.svc.GetPieData()
	c.JSON(http.StatusOK, gin.H{"success": true, "data": pieData})
}

func (h *DashboardHandler) GetTasks(c *gin.Context) {
	tasks := h.svc.GetTasks()
	c.JSON(http.StatusOK, gin.H{"success": true, "data": tasks})
}

func (h *DashboardHandler) GetOverview(c *gin.Context) {
	overview, err := h.svc.GetOverview()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": overview})
}
