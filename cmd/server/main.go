package main

import (
	"log"
	"os"

	"github.com/halolight/halolight-api-go/internal/routes"
	"github.com/halolight/halolight-api-go/pkg/config"
	"github.com/halolight/halolight-api-go/pkg/database"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (ignore error in production where env vars are set directly)
	_ = godotenv.Load()

	// Load configuration
	cfg := config.Load()

	log.Println("🚀 Starting HaloLight API Server...")
	log.Printf("📝 Environment: %s", cfg.AppEnv)
	log.Printf("🔌 Port: %s", cfg.AppPort)

	// Initialize database
	db, err := database.Init(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}

	// Setup router
	r := routes.SetupRouter(cfg, db)

	// Start server
	addr := ":" + cfg.AppPort
	log.Printf("✅ Server running on http://localhost%s", addr)
	log.Println("🏠 Homepage: http://localhost" + addr + "/")
	log.Println("📚 API Documentation: http://localhost" + addr + "/docs")
	log.Println("❤️  Health Check: http://localhost" + addr + "/health")

	if err := r.Run(addr); err != nil {
		log.Printf("❌ Server stopped: %v", err)
		os.Exit(1)
	}
}
