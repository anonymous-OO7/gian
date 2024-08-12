package utils

import (
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"gopkg.in/mail.v2"
)

func EncryptAES(key, text []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(crand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], text)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func DecryptAES(key []byte, cryptoText string) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func SendEmailWithAttachemnt(email []string, subject string, body string, cc []string, bcc []string, filename string) (error, string) {
	emailSender := os.Getenv("EMAIL_SENDER")
	emailPass := os.Getenv("EMAIL_PASSWORD")
	d := mail.NewDialer("smtp.gmail.com", 587, emailSender, emailPass) // Replace with your SMTP server details
	d.TLSConfig = nil
	m := mail.NewMessage()
	m.SetHeader("From", emailSender) // Replace with sender email
	if os.Getenv("ENV") == "production" {
		m.SetHeader("To", email...)
	} else {
		m.SetHeader("To", "akshay.kumar@connectrpl.com")
	}
	if os.Getenv("ENV") == "production" {
		m.SetHeader("Cc", cc...)
		m.SetHeader("Bcc", bcc...)
	} else {
		m.SetHeader("Cc", []string{"ayush@connectrpl.com"}...)
	}
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	m.Attach(filename)
	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err, "Failed"
	}

	return nil, "Success"
}

func SendEmail(email string, subject string, body string, cc []string, bcc []string) (error, string) {
	emailSender := os.Getenv("EMAIL_SENDER")
	emailPass := os.Getenv("EMAIL_PASSWORD")
	d := mail.NewDialer("smtp.gmail.com", 587, emailSender, emailPass) // Replace with your SMTP server details
	d.TLSConfig = nil
	m := mail.NewMessage()
	m.SetHeader("From", emailSender) // Replace with sender email
	if os.Getenv("ENV") == "production" {
		m.SetHeader("To", email)
	} else {
		m.SetHeader("To", "akshay.kumar@connectrpl.com")
	}
	if os.Getenv("ENV") == "production" {
		m.SetHeader("Cc", cc...)
		m.SetHeader("Bcc", bcc...)
	} else {
		m.SetHeader("Cc", []string{"ayush@connectrpl.com"}...)
	}
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err, "Failed"
	}

	return nil, "Success"
}

func HandleError(w http.ResponseWriter, statusCode int, message string, err error) {
	w.WriteHeader(statusCode)
	errorResponse := map[string]string{
		"message": message,
	}
	if err != nil {
		errorResponse["reason"] = err.Error()
	}
	json.NewEncoder(w).Encode(errorResponse)
}

// response struct
type Response struct {
	Status  bool        `json:"status"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// common
type ResponderParams struct {
	Status    int
	IsError   bool
	Message   string
	MainError error
	Data      interface{}
}

// common response sender
func Responder(w http.ResponseWriter, params ResponderParams) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(params.Status)
	var errString string
	if params.MainError != nil {
		errString = params.MainError.Error()
	}
	response := Response{
		Status:  !params.IsError,
		Error:   errString,
		Message: params.Message,
		Data:    params.Data,
	}
	json.NewEncoder(w).Encode(response)
}

func GenerateRandom16Digit() int64 {
	rand.Seed(time.Now().UnixNano())
	// Generate a random number between 10^15 and 10^16 - 1
	return int64(rand.Intn(9e15) + 1e15)
}

func GenerateRandom4Digit() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(10000) // Generate a random number between 0 and 9999
}

func AddLeadingZero(number int) string {
	if number < 10 {
		return fmt.Sprintf("0%d", number)
	}
	return fmt.Sprintf("%d", number)
}

func ConvertDateFormat(inputDate string) (string, error) {
	// Parse the input date string
	parsedDate, err := time.Parse("02-Jan-06", inputDate)
	if err != nil {
		return "", err
	}

	// Format the date into "DD-MM-YYYY"
	formattedDate := parsedDate.Format("02-01-2006")
	return formattedDate, nil
}

func IsValidDate(date string) bool {
	// Define a regular expression for the format DD-MM-YYYY
	dateRegex := regexp.MustCompile(`^(\d{2})-(\d{2})-(\d{4})$`)

	// Check if the date matches the expected format
	if !dateRegex.MatchString(date) {
		return false
	}

	// Extract day, month, and year from the date string
	matches := dateRegex.FindStringSubmatch(date)
	day, _ := strconv.Atoi(matches[1])
	month, _ := strconv.Atoi(matches[2])
	year, _ := strconv.Atoi(matches[3])

	// Validate day, month, and year
	if year < 1000 || year > 9999 {
		return false
	}
	if month < 1 || month > 12 {
		return false
	}
	daysInMonth := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if month == 2 && (year%4 == 0 && (year%100 != 0 || year%400 == 0)) {
		daysInMonth[1] = 29 // Leap year
	}
	if day < 1 || day > daysInMonth[month-1] {
		return false
	}

	// If all checks pass, the date is valid
	return true
}

func GeneratePassword(length int) (string, error) {
	const (
		upperCase    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lowerCase    = "abcdefghijklmnopqrstuvwxyz"
		numbers      = "0123456789"
		specialChars = "!@#$%^&*()-_=+[]{}|;:',.<>?/"
		allChars     = upperCase + lowerCase + numbers + specialChars
	)

	password := make([]byte, length)
	charSets := []string{upperCase, lowerCase, numbers, specialChars}

	// Ensure password contains at least one character from each set
	for i := 0; i < 4; i++ {
		charSet := charSets[i]
		index, err := crand.Int(crand.Reader, big.NewInt(int64(len(charSet))))
		if err != nil {
			return "", err
		}
		password[i] = charSet[index.Int64()]
	}

	// Fill the rest of the password with random characters from all sets
	for i := 4; i < length; i++ {
		index, err := crand.Int(crand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			return "", err
		}
		password[i] = allChars[index.Int64()]
	}

	// Seed math/rand with a cryptographic value
	seed, err := crand.Int(crand.Reader, big.NewInt(1<<63-1))
	if err != nil {
		return "", err
	}
	rand.Seed(seed.Int64())

	// Shuffle the password to ensure randomness
	rand.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

	return string(password), nil
}
