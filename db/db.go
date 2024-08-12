package db

import (
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	DBMu sync.Mutex
)

func SetupDB() error {
	// Directly specifying credentials (not recommended for production)
	dbUser := "postgres"
	dbPassword := "root"
	dbName := "postgres"
	dbPort := "5432" // Default PostgreSQL port
	dbHost := "localhost"

	// dbUser := "postgres"
    // dbPassword := "12345678"
    // dbName := "postgres"
    // dbPort := "5432"
    // dbHost := "192.168.0.118"

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPassword, dbName, dbHost, dbPort)
	var err error
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	return nil
}
