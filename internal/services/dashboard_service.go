package services

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/halolight/halolight-api-go/internal/models"
	"gorm.io/gorm"
)

type DashboardService interface {
	GetStats() (*DashboardStats, error)
	GetVisits() []VisitData
	GetSales() []SalesData
	GetProducts() []ProductData
	GetOrders() []OrderData
	GetActivities() ([]models.ActivityLog, error)
	GetPieData() []PieData
	GetTasks() *TaskData
	GetOverview() (*OverviewData, error)
}

type DashboardStats struct {
	Users      UserStats      `json:"users"`
	Documents  DocumentStats  `json:"documents"`
	Files      FileStats      `json:"files"`
	Teams      TeamStats      `json:"teams"`
	Activities ActivityStats  `json:"activities"`
}

type UserStats struct {
	Total  int64 `json:"total"`
	Active int64 `json:"active"`
}

type DocumentStats struct {
	Total int64 `json:"total"`
}

type FileStats struct {
	Total int64 `json:"total"`
}

type TeamStats struct {
	Total int64 `json:"total"`
}

type ActivityStats struct {
	Recent int64 `json:"recent"`
}

type VisitData struct {
	Day            string `json:"day"`
	Visits         int    `json:"visits"`
	UniqueVisitors int    `json:"uniqueVisitors"`
}

type SalesData struct {
	Month   string `json:"month"`
	Revenue int    `json:"revenue"`
	Orders  int    `json:"orders"`
}

type ProductData struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Sales   int    `json:"sales"`
	Revenue int    `json:"revenue"`
}

type OrderData struct {
	ID       string    `json:"id"`
	Customer string    `json:"customer"`
	Amount   int       `json:"amount"`
	Status   string    `json:"status"`
	Date     time.Time `json:"date"`
}

type PieData struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type TaskData struct {
	Total     int        `json:"total"`
	Completed int        `json:"completed"`
	Pending   int        `json:"pending"`
	Tasks     []TaskItem `json:"tasks"`
}

type TaskItem struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
}

type OverviewData struct {
	System   SystemInfo   `json:"system"`
	Database DatabaseInfo `json:"database"`
}

type SystemInfo struct {
	Uptime      int64  `json:"uptime"`
	GoVersion   string `json:"goVersion"`
	Platform    string `json:"platform"`
	NumGoroutine int   `json:"numGoroutine"`
	MemAlloc    uint64 `json:"memAlloc"`
}

type DatabaseInfo struct {
	Users     int64 `json:"users"`
	Documents int64 `json:"documents"`
	Files     int64 `json:"files"`
	Teams     int64 `json:"teams"`
}

type dashboardService struct {
	db        *gorm.DB
	startTime time.Time
}

func NewDashboardService(db *gorm.DB) DashboardService {
	return &dashboardService{db: db, startTime: time.Now()}
}

func (s *dashboardService) GetStats() (*DashboardStats, error) {
	var stats DashboardStats

	s.db.Model(&models.User{}).Count(&stats.Users.Total)
	s.db.Model(&models.User{}).Where("status = ?", "ACTIVE").Count(&stats.Users.Active)
	s.db.Model(&models.Document{}).Count(&stats.Documents.Total)
	s.db.Model(&models.File{}).Count(&stats.Files.Total)
	s.db.Model(&models.Team{}).Count(&stats.Teams.Total)

	weekAgo := time.Now().AddDate(0, 0, -7)
	s.db.Model(&models.ActivityLog{}).Where("created_at >= ?", weekAgo).Count(&stats.Activities.Recent)

	return &stats, nil
}

func (s *dashboardService) GetVisits() []VisitData {
	days := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	result := make([]VisitData, len(days))
	for i, day := range days {
		result[i] = VisitData{
			Day:            day,
			Visits:         rand.Intn(1000) + 500,         // #nosec G404 - Mock data only
			UniqueVisitors: rand.Intn(500) + 200,          // #nosec G404 - Mock data only
		}
	}
	return result
}

func (s *dashboardService) GetSales() []SalesData {
	months := []string{"Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	result := make([]SalesData, len(months))
	for i, month := range months {
		result[i] = SalesData{
			Month:   month,
			Revenue: rand.Intn(50000) + 10000, // #nosec G404 - Mock data only
			Orders:  rand.Intn(500) + 100,     // #nosec G404 - Mock data only
		}
	}
	return result
}

func (s *dashboardService) GetProducts() []ProductData {
	return []ProductData{
		{ID: "1", Name: "Product A", Sales: 1234, Revenue: 12340},
		{ID: "2", Name: "Product B", Sales: 987, Revenue: 9870},
		{ID: "3", Name: "Product C", Sales: 765, Revenue: 7650},
		{ID: "4", Name: "Product D", Sales: 543, Revenue: 5430},
		{ID: "5", Name: "Product E", Sales: 321, Revenue: 3210},
	}
}

func (s *dashboardService) GetOrders() []OrderData {
	now := time.Now()
	return []OrderData{
		{ID: "1", Customer: "John Doe", Amount: 299, Status: "completed", Date: now},
		{ID: "2", Customer: "Jane Smith", Amount: 199, Status: "pending", Date: now},
		{ID: "3", Customer: "Bob Wilson", Amount: 499, Status: "processing", Date: now},
		{ID: "4", Customer: "Alice Brown", Amount: 149, Status: "completed", Date: now},
		{ID: "5", Customer: "Charlie Davis", Amount: 399, Status: "shipped", Date: now},
	}
}

func (s *dashboardService) GetActivities() ([]models.ActivityLog, error) {
	var activities []models.ActivityLog
	err := s.db.Preload("Actor").Order("created_at DESC").Limit(10).Find(&activities).Error
	return activities, err
}

func (s *dashboardService) GetPieData() []PieData {
	return []PieData{
		{Name: "Documents", Value: 35},
		{Name: "Images", Value: 25},
		{Name: "Videos", Value: 20},
		{Name: "Audio", Value: 10},
		{Name: "Others", Value: 10},
	}
}

func (s *dashboardService) GetTasks() *TaskData {
	return &TaskData{
		Total:     24,
		Completed: 18,
		Pending:   6,
		Tasks: []TaskItem{
			{ID: "1", Title: "Review documents", Status: "completed", Priority: "high"},
			{ID: "2", Title: "Update user permissions", Status: "pending", Priority: "medium"},
			{ID: "3", Title: "Deploy new features", Status: "in_progress", Priority: "high"},
			{ID: "4", Title: "Fix reported bugs", Status: "pending", Priority: "low"},
		},
	}
}

func (s *dashboardService) GetOverview() (*OverviewData, error) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	var dbInfo DatabaseInfo
	s.db.Model(&models.User{}).Count(&dbInfo.Users)
	s.db.Model(&models.Document{}).Count(&dbInfo.Documents)
	s.db.Model(&models.File{}).Count(&dbInfo.Files)
	s.db.Model(&models.Team{}).Count(&dbInfo.Teams)

	return &OverviewData{
		System: SystemInfo{
			Uptime:       int64(time.Since(s.startTime).Seconds()),
			GoVersion:    runtime.Version(),
			Platform:     runtime.GOOS,
			NumGoroutine: runtime.NumGoroutine(),
			MemAlloc:     memStats.Alloc,
		},
		Database: dbInfo,
	}, nil
}
