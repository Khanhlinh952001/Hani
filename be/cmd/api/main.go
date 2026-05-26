package main

import (
	"be/internal/config"
	"be/internal/middleware"
	"be/internal/routes"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	// "be/internal/config"
	// "be/internal/middleware"
	// "be/internal/routes"
)

func loadEnv() {
	// air chạy ./tmp/main từ project root — ưu tiên .env cạnh go.mod
	candidates := []string{".env"}
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(wd, ".env"))
	}
	for _, p := range candidates {
		if err := godotenv.Load(p); err == nil {
			return
		}
	}
	log.Println("no .env file, using system env")
}

func main() {
	loadEnv()

	config.ConnectDB()

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	routes.SetupRoutes(r)

	r.Run(":8080")
}
