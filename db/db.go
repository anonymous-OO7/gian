package db

import (
	"fmt"
	"gian/email"
	"os"
	"os/exec"
	"path/filepath"
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

func Backup() error {
	// Define the backup directory and file name
	backupDir := filepath.Join("desktop", "dbbackup")
	backupFile := filepath.Join(backupDir, "latest_backup.dump") // Fixed file name

	// Ensure the backup directory exists
	if err := os.MkdirAll(backupDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}

	// Define PostgreSQL connection parameters
	dbUser := "postgres"
	dbPassword := "root"
	dbName := "postgres"
	dbHost := "localhost"
	dbPort := "5432" // Default PostgreSQL port

	// Set the PGPASSWORD environment variable for pg_dump
	os.Setenv("PGPASSWORD", dbPassword)

	// Prepare the pg_dump command
	cmd := exec.Command("pg_dump", "-U", dbUser, "-h", dbHost, "-p", dbPort, "-F", "c", "-b", "-v", "-f", backupFile, dbName)

	// Run the command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pg_dump failed: %s\n%s", err, output)
	}

	fmt.Printf("Backup completed successfully. File: %s\n", backupFile)

	// Send the backup file via email
	if err := email.SendDatabaseEmail(backupFile); err != nil {
		return fmt.Errorf("failed to send backup email: %v", err)
	}

	return nil
}
