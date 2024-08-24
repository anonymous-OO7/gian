package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"gian/db"
	"gian/models"
	"gian/utils"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mail.v2"
	"gorm.io/gorm"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))
	name := strings.TrimSpace(r.FormValue("name"))
	role := strings.TrimSpace(r.FormValue("role"))
	username := strings.TrimSpace(r.FormValue("username"))
	phone := strings.TrimSpace(r.FormValue("phone"))
	gender := strings.TrimSpace(r.FormValue("gender"))
	organisation := strings.TrimSpace(r.FormValue("organisation"))
	title := strings.TrimSpace(r.FormValue("title"))
	country := strings.TrimSpace(r.FormValue("country"))

	// Basic validation
	if email == "" || name == "" || role == "" || password == "" || username == "" || phone == "" || gender == "" || organisation == "" || title == "" || country == "" || !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email) {
		utils.HandleError(w, http.StatusBadRequest, "Invalid name, email, or role", nil)
		return
	}

	if password == "" || len(password) < 6 {
		utils.HandleError(w, http.StatusBadRequest, "Password should be more than 6 characters", nil)
		return
	}

	// Hash the password using bcrypt for security
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	// Create a new user instance
	newUser := models.User{
		Name:         name,
		Email:        email,
		Password:     string(hashedPassword),
		Role:         role,
		Username:     username,
		Phone:        phone,
		Gender:       gender,
		Organisation: organisation,
		Title:        title,
		Country:      country,
		Uuid:         uuid.New().String(),
	}

	// Save the user to the database
	if err := db.DB.Create(&newUser).Error; err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Unable to create user", err)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func isValidCredentials(email string, password string) (models.User, error) {
	// Query the database for the user with the given email
	var user models.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}

	// Compare the provided password with the hashed password from the database
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, err
	}
	return user, nil // Replace with your validation logic
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Parse the request to get the email and password
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := isValidCredentials(email, password)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid Credentials", "reason": err.Error()})
		return
	} else {
		token := jwt.New(jwt.SigningMethodHS256)
		expTime := time.Now().Add(time.Hour * 24).Unix()
		claims := token.Claims.(jwt.MapClaims)
		claims["email"] = email
		claims["name"] = user.Name
		claims["role"] = user.Role
		claims["exp"] = expTime
		tokenString, err := token.SignedString([]byte("your-secret-key"))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Unable to generate token", "reason": err.Error()})
			return
		}

		responseBody := map[string]interface{}{
			"message": "success",
			"token":   tokenString,
			"email":   user.Email,
			"name":    user.Name,
			"role":    user.Role,
		}

		jsonResponse, _ := json.Marshal(responseBody)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.FormValue("email"))
	oldPassword := strings.TrimSpace(r.FormValue("old_password"))
	newPassword := strings.TrimSpace(r.FormValue("new_password"))

	// Validate email
	if email == "" || !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email) {
		utils.HandleError(w, http.StatusBadRequest, "Invalid email", nil)
		return
	}

	// Validate old password
	if oldPassword == "" || len(newPassword) < 6 {
		utils.HandleError(w, http.StatusBadRequest, "Old password is required", nil)
		return
	}

	// Validate new password
	if newPassword == "" || len(newPassword) < 6 {
		utils.HandleError(w, http.StatusBadRequest, "New password should be more than 6 characters", nil)
		return
	}

	// Find user by email
	var user models.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		utils.HandleError(w, http.StatusNotFound, "User not found", err)
		return
	}

	// Check old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Old password is incorrect", nil)
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Unknown error", err)
		return
	}

	// Update password
	user.Password = string(hashedPassword)
	if err := db.DB.Save(&user).Error; err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to update password", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully"})
}

func MasterUpdatePassword(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.FormValue("email"))
	newPassword := strings.TrimSpace(r.FormValue("new_password"))

	// Validate email
	if email == "" || !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email) {
		utils.HandleError(w, http.StatusBadRequest, "Invalid email", nil)
		return
	}

	// Validate new password
	if newPassword == "" || len(newPassword) < 6 {
		utils.HandleError(w, http.StatusBadRequest, "New password should be more than 6 characters", nil)
		return
	}

	// Find user by email
	var user models.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		utils.HandleError(w, http.StatusNotFound, "User not found", err)
		return
	}
	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Unknown error", err)
		return
	}

	// Update password
	user.Password = string(hashedPassword)
	if err := db.DB.Save(&user).Error; err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to update password", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully"})
}

// GetOtp handles the OTP generation and sending process
func GetOtp(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.FormValue("email"))

	var user models.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	user.Otp = otp
	if err := db.DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to save OTP", http.StatusInternalServerError)
		return
	}

	subject := "Your OTP Code"
	body := fmt.Sprintf("Your OTP code is: %s", otp)
	if err, msg := SendEmail([]string{user.Email}, subject, body); err != nil {
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "OTP sent successfully"})
}

func VerifyOtp(w http.ResponseWriter, r *http.Request) {
	// Parse and validate email and otp from the request
	email := strings.TrimSpace(r.FormValue("email"))
	otp := strings.TrimSpace(r.FormValue("otp"))

	if email == "" || otp == "" {
		http.Error(w, "Email and OTP are required", http.StatusBadRequest)
		return
	}

	// Fetch the user based on the email
	var user models.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Verify the OTP
	if user.Otp != otp {
		http.Error(w, "Invalid OTP", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	expTime := time.Now().Add(time.Hour * 24).Unix()
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["name"] = user.Name
	claims["role"] = user.Role
	claims["exp"] = expTime
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unable to generate token", "reason": err.Error()})
		return
	}

	// Respond with the token and user details
	responseBody := map[string]interface{}{
		"message": "success",
		"token":   tokenString,
		"email":   user.Email,
		"name":    user.Name,
		"role":    user.Role,
		"id":      user.ID,
		"uuid":    user.Uuid,

		"existingUser": true,
		"status":       true,
	}

	jsonResponse, _ := json.Marshal(responseBody)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// SendEmail sends an email with the specified subject and body
func SendEmail(to []string, subject string, body string) (error, string) {
	emailSender := "gauravyadav00729@gmail.com"
	emailPass := "aiif rfjh jwdd iptp"
	d := mail.NewDialer("smtp.gmail.com", 587, emailSender, emailPass)
	d.TLSConfig = nil

	m := mail.NewMessage()
	m.SetHeader("From", emailSender)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if err := d.DialAndSend(m); err != nil {
		return err, "Failed to send email"
	}

	return nil, "Email sent successfully"
}
