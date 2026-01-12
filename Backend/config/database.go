package config

import (
	models2 "backend-test-mekari/internal/models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	_ = godotenv.Load()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	timezone := os.Getenv("DB_TIMEZONE")

	if timezone == "" {
		timezone = "Asia/Jakarta"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, timezone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("Database connection established")

	// Auto Migration
	err = db.AutoMigrate(
		&models2.User{},
		&models2.Expense{},
		&models2.Approval{},
	)
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	SeedData(db)

	log.Println("Database migration & seeding completed")

	return db
}

func SeedData(db *gorm.DB) {
	var userCount int64
	db.Model(&models2.User{}).Count(&userCount)

	if userCount == 0 {
		// 1. Inisialisasi Users
		users := []models2.User{
			{
				ID:        1,
				Name:      "Budi Employee",
				Email:     "employee@mekari.com",
				Password:  "password123", // Mock password
				Role:      "employee",
				CreatedAt: time.Now(),
			},
			{
				ID:        2,
				Name:      "Siti Manager",
				Email:     "manager@mekari.com",
				Password:  "password123",
				Role:      "manager",
				CreatedAt: time.Now(),
			},
		}

		for _, u := range users {
			db.Create(&u)
		}
		log.Println("Seeding: Users successfully created")

		// 2. Inisialisasi Expenses (Hanya jika belum ada data)
		var expenseCount int64
		db.Model(&models2.Expense{}).Count(&expenseCount)

		if expenseCount == 0 {
			expenses := []models2.Expense{
				{
					UserID:      1,
					Amount:      850000,
					Description: "Membeli tinta printer dan kertas A4",
					Status:      "auto_approved",
					ExternalID:  "EXP-AUTO-001",
					SubmittedAt: time.Now().Add(-72 * time.Hour),
				},
				{
					UserID:      1,
					Amount:      2500000,
					Description: "Makan malam dengan klien",
					Status:      "pending",
					ExternalID:  "EXP-PEND-002",
					SubmittedAt: time.Now().Add(-24 * time.Hour),
				},
				{
					UserID:      1,
					Amount:      12000000,
					Description: "Pembelian laptop operasional baru",
					Status:      "rejected",
					ExternalID:  "EXP-REJ-003",
					SubmittedAt: time.Now().Add(-120 * time.Hour),
				},
				{
					UserID:      1,
					Amount:      150000,
					Description: "Biaya parkir dan tol operasional",
					Status:      "completed",
					ExternalID:  "EXP-COMP-004",
					SubmittedAt: time.Now().Add(-48 * time.Hour),
					ProcessedAt: func() *time.Time { t := time.Now().Add(-45 * time.Hour); return &t }(),
				},
			}

			for _, e := range expenses {
				db.Create(&e)
			}
			log.Println("Seeding: Sample expenses successfully created")
		}
	}
}
