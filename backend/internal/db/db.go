package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New() (*gorm.DB, error) {
	_ = godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB"))
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	return db, nil
}

func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
