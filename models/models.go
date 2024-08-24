package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string `gorm:"not null"`                 // User's full name
	Email        string `gorm:"not null;unique;index"`    // Unique email address
	Password     string `gorm:"not null"`                 // Hashed password
	Role         string `gorm:"not null;default:'admin'"` // User role (e.g., admin, user)
	Username     string `gorm:"not null;unique"`          // Unique username
	Phone        string `gorm:"size:20"`                  // Contact phone number
	Gender       string `gorm:"size:10"`                  // Gender of the user
	Organisation string `gorm:"size:100"`                 // Organisation of the user
	Title        string `gorm:"size:100"`                 // Job title or designation
	Country      string `gorm:"size:50"`                  // User's country
	Otp          string
	Uuid         string `gorm:"not null;unique;index"` // UUID for additional identifier, if needed
}

type Jobs struct {
	gorm.Model
	UserID      int    `gorm:"not null"`                  // Foreign key to Users table (referred user)
	Status      string `gorm:"type:varchar(50);not null"` // Status of the referral (e.g., "active", "inactive")
	Uuid        string `gorm:"not null;unique;index"`     // UUID for additional identifier, if needed
	CompanyName string `gorm:"not null"`                  // Name of the company offering the job
	Position    string `gorm:"not null"`                  // Title of the job
	Location    string `gorm:"not null"`                  // Location of the job
	Type        string `gorm:"not null"`                  // Type of the job (e.g., "Full-time", "Part-time")
	Description string `gorm:"type:text;not null"`        // Description of the job
	Field       string `gorm:"not null"`                  // Type or category of the job (e.g., "Engineering", "Marketing")
	Owner       string `gorm:"not null"`                  // uuid of owner
	MinPay      int    `gorm:"not null"`                  // Minimum pay for the job
	MaxPay      int    `gorm:"not null"`                  // Maximum pay for the job
	Price       int    `gorm:"not null"`                  // Maximum pay for the job
	TotalEmp    int    `gorm:"not null"`                  // Maximum pay for the job
	LogoUrl     string `gorm:"not null"`                  // Title of the job
}
type Saved struct {
	gorm.Model
	UserID string `gorm:"type:varchar(50);not null"`
	JobIDs string `gorm:"type:text"`
}

type Applications struct {
	gorm.Model
	ApplicantID string `gorm:"not null"`
	JobIDs      string `gorm:"type:text"`
	Uuid        string `gorm:"not null"`
}
