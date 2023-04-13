package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDBClient() (*gorm.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		dbPass = "postgres"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "local"
	}

	dbSSL := os.Getenv("DB_SSL")
	if dbSSL == "" {
		dbSSL = "disable"
	}
	connStr := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v", dbHost, dbPort, dbUser, dbPass, dbName, dbSSL)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{TranslateError: true})
	return db, err
}
