package config

import (
	"fmt"
	"log"
	"os"

	"be/internal/db"
	"be/internal/modules/memories"
	"be/internal/modules/messages"
	"be/internal/modules/sessions"
	"be/internal/modules/users"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() {
	loadEnv()

	dsn := buildDSN()

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	db.DB = conn

	enableVectorExtension(conn)

	autoMigrate()
	seedDemoUser()
	seedAdminUser()

	log.Println("Database connected")
}

func seedDemoUser() {
	var count int64
	db.DB.Model(&users.User{}).Where("id = ?", 1).Count(&count)
	if count > 0 {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("demo123456"), bcrypt.DefaultCost)
	if err != nil {
		log.Println("seed user:", err)
		return
	}

	demo := users.User{
		Name:     "Bạn",
		Email:    "demo@hani.app",
		Password: string(hash),
		Status:   1,
		Level:    3,
	}
	if err := db.DB.Create(&demo).Error; err != nil {
		log.Println("seed user:", err)
		return
	}
	log.Printf("seeded demo user id=%d email=%s", demo.ID, demo.Email)
}

func seedAdminUser() {
	adminEmail := getEnv("ADMIN_EMAIL", "admin@hani.app")
	adminPass := getEnv("ADMIN_PASSWORD", "admin123456")

	var existing users.User
	err := db.DB.Where("email = ?", adminEmail).First(&existing).Error
	if err != nil {
		hash, herr := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)
		if herr != nil {
			log.Println("seed admin:", herr)
			return
		}
		admin := users.User{
			Name:     "Admin",
			Email:    adminEmail,
			Password: string(hash),
			Role:     users.RoleAdmin,
			Status:   1,
		}
		if err := db.DB.Create(&admin).Error; err != nil {
			log.Println("seed admin:", err)
			return
		}
		log.Printf("created admin %s (password from ADMIN_PASSWORD)", adminEmail)
		return
	}

	if existing.Role != users.RoleAdmin {
		db.DB.Model(&existing).Update("role", users.RoleAdmin)
		log.Printf("promoted admin: %s", adminEmail)
	}
}

func enableVectorExtension(conn *gorm.DB) {
	if err := conn.Exec("CREATE EXTENSION IF NOT EXISTS vector").Error; err != nil {
		log.Println("vector extension:", err)
	}
}

func loadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Println(".env file not found")
	}
}

func buildDSN() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "hani")
	password := getEnv("POSTGRES_PASSWORD", "")
	dbname := getEnv("POSTGRES_DB", "hani_db")

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host,
		user,
		password,
		dbname,
		port,
	)
}

func autoMigrate() {
	err := db.DB.AutoMigrate(
		&users.User{},
		&sessions.Session{},
		&messages.Message{},
		&memories.Memory{},
	)

	if err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
