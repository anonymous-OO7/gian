package email

import (
	"bytes"
	"encoding/json"
	"gian/db"
	"gian/models"
	"gian/utils"
	"net/http"
	"os"
	"text/template"

	"gopkg.in/mail.v2"
)

type CandidateDetails struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	Phone  int    `json:"phone"`
	Resume string `json:"resume"` // Assuming resume is a file path or base64 encoded string
	Cover  string `json:"cover"`  // Assuming cover letter is a file path or base64 encoded string
}

// Function to handle job application submission
func JobApply(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming JSON request
	var candidate CandidateDetails
	err := json.NewDecoder(r.Body).Decode(&candidate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Invalid request payload",
			"reason":  err.Error(),
		})
		return
	}

	// Fetch user details based on role and department
	var user models.User
	err = db.DB.Where("role = ? AND department = ?", utils.SALES, utils.SALES).First(&user).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Unable to fetch user details",
			"reason":  err.Error(),
		})
		return
	}

	// Prepare the email body
	emailBodyHtml, err := renderEmailBody([]CandidateDetails{candidate})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Unable to create email body",
			"reason":  err.Error(),
		})
		return
	}

	// Send the email
	receiver := []string{user.Email} // Send to the fetched user's email
	subject := "New Job Application Received"
	bcc := []string{}
	cc := []string{}
	err, status := SendEmailMeet(receiver, subject, emailBodyHtml, cc, bcc)

	if err != nil {
		w.WriteHeader(http.StatusFailedDependency)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Unable to send email", "status": status, "reason": err.Error()})
		return
	}

	// After successfully sending the email, store the application in the database
	err = db.DB.Create(&candidate).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Failed to store job application",
			"reason":  err.Error(),
		})
		return
	}

	// Prepare the JSON response
	responseBody := map[string]interface{}{
		"message": "Application submitted successfully",
		"data":    candidate,
	}

	jsonResponse, _ := json.Marshal(responseBody)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func renderEmailBody(data []CandidateDetails) (string, error) {
	// Load the HTML template
	t, err := template.ParseFiles("email/templates/MeetMail.html")
	if err != nil {
		return "", err
	}

	// Execute the template with the data
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func SendEmailMeet(email []string, subject string, body string, cc []string, bcc []string) (error, string) {
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
	m.SetBody("text/html", body)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err, "Failed"
	}

	return nil, "Success"
}
