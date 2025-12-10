package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/halolight/halolight-api-go/internal/handlers"
	"github.com/halolight/halolight-api-go/internal/middleware"
	"github.com/halolight/halolight-api-go/internal/repository"
	"github.com/halolight/halolight-api-go/internal/services"
	"github.com/halolight/halolight-api-go/pkg/config"
	"gorm.io/gorm"
)

func SetupRouter(cfg config.Config, db *gorm.DB) *gin.Engine {
	// Set Gin mode
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Apply CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Initialize home handler
	homeHandler := handlers.NewHomeHandler(cfg)

	// Root endpoint - Beautiful HTML homepage
	r.GET("/", homeHandler.Home)

	// Health check endpoint
	r.GET("/health", homeHandler.Health)

	// Swagger documentation (redirect to static swagger-ui)
	r.GET("/docs", func(c *gin.Context) {
		c.Redirect(302, "/docs/index.html")
	})

	// Serve swagger-ui static files
	r.Static("/docs", "./docs/swagger-ui")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)

	// Initialize services
	authSvc := services.NewAuthService(cfg, userRepo)
	userSvc := services.NewUserService(userRepo)
	roleSvc := services.NewRoleService(db)
	permissionSvc := services.NewPermissionService(db)
	teamSvc := services.NewTeamService(db)
	documentSvc := services.NewDocumentService(db)
	fileSvc := services.NewFileService(db)
	folderSvc := services.NewFolderService(db)
	calendarSvc := services.NewCalendarService(db)
	notificationSvc := services.NewNotificationService(db)
	messageSvc := services.NewMessageService(db)
	dashboardSvc := services.NewDashboardService(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authSvc, refreshTokenRepo, cfg)
	userHandler := handlers.NewUserHandler(userSvc)
	roleHandler := handlers.NewRoleHandler(roleSvc)
	permissionHandler := handlers.NewPermissionHandler(permissionSvc)
	teamHandler := handlers.NewTeamHandler(teamSvc)
	documentHandler := handlers.NewDocumentHandler(documentSvc)
	fileHandler := handlers.NewFileHandler(fileSvc, folderSvc)
	folderHandler := handlers.NewFolderHandler(folderSvc)
	calendarHandler := handlers.NewCalendarHandler(calendarSvc)
	notificationHandler := handlers.NewNotificationHandler(notificationSvc)
	messageHandler := handlers.NewMessageHandler(messageSvc)
	dashboardHandler := handlers.NewDashboardHandler(dashboardSvc)

	// API routes
	api := r.Group("/api")
	{
		// ==================== Auth Routes (Public) ====================
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
		}

		// Auth routes requiring authentication
		authProtected := api.Group("/auth")
		authProtected.Use(middleware.AuthMiddleware(cfg))
		{
			authProtected.GET("/me", authHandler.Me)
			authProtected.POST("/logout", authHandler.Logout)
		}

		// ==================== Users Routes ====================
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(cfg))
		{
			users.GET("", userHandler.List)
			users.GET("/:id", userHandler.Get)
			users.POST("", userHandler.Create)
			users.PATCH("/:id", userHandler.Update)
			users.PATCH("/:id/status", userHandler.UpdateStatus)
			users.POST("/batch-delete", userHandler.BatchDelete)
			users.DELETE("/:id", userHandler.Delete)
		}

		// ==================== Roles Routes ====================
		roles := api.Group("/roles")
		roles.Use(middleware.AuthMiddleware(cfg))
		{
			roles.GET("", roleHandler.List)
			roles.GET("/:id", roleHandler.Get)
			roles.POST("", roleHandler.Create)
			roles.PATCH("/:id", roleHandler.Update)
			roles.DELETE("/:id", roleHandler.Delete)
			roles.POST("/:id/permissions", roleHandler.AssignPermissions)
		}

		// ==================== Permissions Routes ====================
		permissions := api.Group("/permissions")
		permissions.Use(middleware.AuthMiddleware(cfg))
		{
			permissions.GET("", permissionHandler.List)
			permissions.GET("/:id", permissionHandler.Get)
			permissions.POST("", permissionHandler.Create)
			permissions.DELETE("/:id", permissionHandler.Delete)
		}

		// ==================== Teams Routes ====================
		teams := api.Group("/teams")
		teams.Use(middleware.AuthMiddleware(cfg))
		{
			teams.GET("", teamHandler.List)
			teams.GET("/:id", teamHandler.Get)
			teams.POST("", teamHandler.Create)
			teams.PATCH("/:id", teamHandler.Update)
			teams.DELETE("/:id", teamHandler.Delete)
			teams.POST("/:id/members", teamHandler.AddMember)
			teams.DELETE("/:id/members/:userId", teamHandler.RemoveMember)
		}

		// ==================== Documents Routes ====================
		documents := api.Group("/documents")
		documents.Use(middleware.AuthMiddleware(cfg))
		{
			documents.GET("", documentHandler.List)
			documents.GET("/:id", documentHandler.Get)
			documents.POST("", documentHandler.Create)
			documents.PUT("/:id", documentHandler.Update)
			documents.PATCH("/:id/rename", documentHandler.Rename)
			documents.POST("/:id/move", documentHandler.Move)
			documents.POST("/:id/tags", documentHandler.UpdateTags)
			documents.POST("/:id/share", documentHandler.Share)
			documents.POST("/:id/unshare", documentHandler.Unshare)
			documents.POST("/batch-delete", documentHandler.BatchDelete)
			documents.DELETE("/:id", documentHandler.Delete)
		}

		// ==================== Files Routes ====================
		files := api.Group("/files")
		files.Use(middleware.AuthMiddleware(cfg))
		{
			files.POST("/upload", fileHandler.Upload)
			files.POST("/folder", fileHandler.CreateFolder)
			files.GET("", fileHandler.List)
			files.GET("/storage", fileHandler.GetStorage)
			files.GET("/storage-info", fileHandler.GetStorage)
			files.GET("/:id", fileHandler.Get)
			files.GET("/:id/download-url", fileHandler.GetDownloadURL)
			files.PATCH("/:id/rename", fileHandler.Rename)
			files.POST("/:id/move", fileHandler.Move)
			files.POST("/:id/copy", fileHandler.Copy)
			files.PATCH("/:id/favorite", fileHandler.ToggleFavorite)
			files.POST("/:id/share", fileHandler.Share)
			files.POST("/batch-delete", fileHandler.BatchDelete)
			files.DELETE("/:id", fileHandler.Delete)
		}

		// ==================== Folders Routes ====================
		folders := api.Group("/folders")
		folders.Use(middleware.AuthMiddleware(cfg))
		{
			folders.GET("", folderHandler.List)
			folders.GET("/tree", folderHandler.GetTree)
			folders.GET("/:id", folderHandler.Get)
			folders.POST("", folderHandler.Create)
			folders.DELETE("/:id", folderHandler.Delete)
		}

		// ==================== Calendar Routes ====================
		calendar := api.Group("/calendar")
		calendar.Use(middleware.AuthMiddleware(cfg))
		{
			events := calendar.Group("/events")
			{
				events.GET("", calendarHandler.List)
				events.GET("/:id", calendarHandler.Get)
				events.POST("", calendarHandler.Create)
				events.PUT("/:id", calendarHandler.Update)
				events.PATCH("/:id/reschedule", calendarHandler.Reschedule)
				events.POST("/:id/attendees", calendarHandler.AddAttendee)
				events.DELETE("/:id/attendees/:attendeeId", calendarHandler.RemoveAttendee)
				events.POST("/batch-delete", calendarHandler.BatchDelete)
				events.DELETE("/:id", calendarHandler.Delete)
			}
		}

		// ==================== Notifications Routes ====================
		notifications := api.Group("/notifications")
		notifications.Use(middleware.AuthMiddleware(cfg))
		{
			notifications.GET("", notificationHandler.List)
			notifications.GET("/unread-count", notificationHandler.GetUnreadCount)
			notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
			notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
			notifications.DELETE("/:id", notificationHandler.Delete)
		}

		// ==================== Messages Routes ====================
		messages := api.Group("/messages")
		messages.Use(middleware.AuthMiddleware(cfg))
		{
			messages.GET("/conversations", messageHandler.GetConversations)
			messages.GET("/conversations/:id", messageHandler.GetConversation)
			messages.POST("", messageHandler.SendMessage)
			messages.PUT("/:id/read", messageHandler.MarkAsRead)
			messages.DELETE("/:id", messageHandler.DeleteMessage)
		}

		// ==================== Dashboard Routes ====================
		dashboard := api.Group("/dashboard")
		dashboard.Use(middleware.AuthMiddleware(cfg))
		{
			dashboard.GET("/stats", dashboardHandler.GetStats)
			dashboard.GET("/visits", dashboardHandler.GetVisits)
			dashboard.GET("/sales", dashboardHandler.GetSales)
			dashboard.GET("/products", dashboardHandler.GetProducts)
			dashboard.GET("/orders", dashboardHandler.GetOrders)
			dashboard.GET("/activities", dashboardHandler.GetActivities)
			dashboard.GET("/pie", dashboardHandler.GetPieData)
			dashboard.GET("/tasks", dashboardHandler.GetTasks)
			dashboard.GET("/overview", dashboardHandler.GetOverview)
		}
	}

	return r
}
