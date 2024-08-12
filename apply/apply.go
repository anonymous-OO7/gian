package apply

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gian/db"
	"gian/models"
	"gian/utils"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"gopkg.in/mail.v2"
	"gorm.io/gorm"
)

type CandidateDetails struct {
	Name   string
	Email  string
	Phone  string
	Resume string
	Cover  string
}

func renderEmailBody(data CandidateDetails) (string, error) {
	t, err := template.ParseFiles("email/templates/MeetMail.html")
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func SendEmailMeet(email []string, subject string, body string, attachments []string, cc []string, bcc []string) (error, string) {
	emailSender := "gauravyadav00729@gmail.com"
	emailPass := "aiif rfjh jwdd iptp"
	d := mail.NewDialer("smtp.gmail.com", 587, emailSender, emailPass)
	d.TLSConfig = nil
	m := mail.NewMessage()
	m.SetHeader("From", emailSender)
	m.SetHeader("To", email...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	for _, attachment := range attachments {
		m.Attach(attachment)
	}

	if err := d.DialAndSend(m); err != nil {
		return err, "Failed"
	}

	return nil, "Success"
}

func SaveApplication(w http.ResponseWriter, r *http.Request) {

	applicantID := r.Header.Get("user_id")
	name := r.Header.Get("name")
	email := r.Header.Get("email")

	jobUUID := r.FormValue("job_id")
	cover := r.FormValue("cover")
	ownersUUID := r.FormValue("owners_uuid")

	fmt.Println(applicantID, name, email, jobUUID, cover, ownersUUID)

	if applicantID == "" || jobUUID == "" || cover == "" || name == "" || email == "" || ownersUUID == "" {
		utils.HandleError(w, http.StatusBadRequest, "Missing required fields", nil)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Error parsing form", err)
		return
	}

	file, handler, err := r.FormFile("resume")
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Error retrieving resume file", err)
		return
	}
	defer file.Close()

	tempFile, err := os.CreateTemp("", handler.Filename)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Error creating temporary file", err)
		return
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Error saving resume file", err)
		return
	}

	// Ensure the temporary file is closed before proceeding
	tempFile.Close()

	candidate := CandidateDetails{
		Email:  email,
		Name:   name,
		Phone:  "0", // Replace with the actual phone number if available
		Resume: tempFile.Name(),
		Cover:  cover,
	}

	emailBody, err := renderEmailBody(candidate)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to render email body", err)
		return
	}

	attachments := []string{tempFile.Name()}

	var job models.Jobs
	if err := db.DB.Where("uuid = ?", jobUUID).First(&job).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.HandleError(w, http.StatusNotFound, "Job not found", err)
		} else {
			utils.HandleError(w, http.StatusInternalServerError, "Failed to retrieve job", err)
		}
		return
	}

	var user models.User
	if err := db.DB.Where("uuid = ?", ownersUUID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.HandleError(w, http.StatusNotFound, "User not found", err)
		} else {
			utils.HandleError(w, http.StatusInternalServerError, "Failed to retrieve user", err)
		}
		return
	}

	if err, status := SendEmailMeet([]string{user.Email}, "Job Application: "+job.Position, emailBody, attachments, nil, nil); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to send email", err)
		return
	} else if status != "Success" {
		utils.HandleError(w, http.StatusInternalServerError, "Email sending failed", nil)
		return
	}

	// Save the application in the database
	var application models.Applications
	err = db.DB.Where("applicant_id = ?", applicantID).First(&application).Error
	if err == gorm.ErrRecordNotFound {
		// Applicant is not found, create a new record
		jobIDs := []string{jobUUID}
		jobIDsJSON, _ := json.Marshal(jobIDs)
		application = models.Applications{
			ApplicantID: applicantID,
			JobIDs:      string(jobIDsJSON),
			Uuid:        uuid.New().String(),
		}
		if err := db.DB.Create(&application).Error; err != nil {
			utils.HandleError(w, http.StatusInternalServerError, "Failed to save application", err)
			return
		}
	} else if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to find applicant", err)
		return
	} else {
		// Applicant found, update the JobIDs
		var jobIDs []string
		if err := json.Unmarshal([]byte(application.JobIDs), &jobIDs); err != nil {
			utils.HandleError(w, http.StatusInternalServerError, "Failed to parse job IDs", err)
			return
		}

		// Check if jobUUID is already saved
		for _, id := range jobIDs {
			if id == jobUUID {
				responseBody := map[string]interface{}{
					"message": "Job already applied",
					"status":  true,
				}
				jsonResponse, _ := json.Marshal(responseBody)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(jsonResponse)
				return
			}
		}

		// Append new jobUUID and encode back to JSON
		jobIDs = append(jobIDs, jobUUID)
		jobIDsJSON, _ := json.Marshal(jobIDs)
		application.JobIDs = string(jobIDsJSON)
		if err := db.DB.Save(&application).Error; err != nil {
			utils.HandleError(w, http.StatusInternalServerError, "Failed to save application", err)
			return
		}
	}

	responseBody := map[string]interface{}{
		"message": "Job application saved successfully",
		"status":  true,
	}

	jsonResponse, _ := json.Marshal(responseBody)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)

}

func GetApplications(w http.ResponseWriter, r *http.Request) {
	applicantID := r.Header.Get("user_id")

	if applicantID == "" {
		utils.HandleError(w, http.StatusBadRequest, "Missing applicant_id", nil)
		return
	}

	var application models.Applications

	if err := db.DB.Where("applicant_id = ?", applicantID).First(&application).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.HandleError(w, http.StatusNotFound, "No applications found for the applicant", err)
		} else {
			utils.HandleError(w, http.StatusInternalServerError, "Failed to find applications", err)
		}
		return
	}

	var jobIDs []string
	if err := json.Unmarshal([]byte(application.JobIDs), &jobIDs); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to parse job IDs", err)
		return
	}

	var jobsList []models.Jobs
	if len(jobIDs) > 0 {
		if err := db.DB.Where("uuid IN ?", jobIDs).Find(&jobsList).Error; err != nil {
			utils.HandleError(w, http.StatusInternalServerError, "Failed to retrieve jobs", err)
			return
		}
	}

	responseBody := map[string]interface{}{
		"message": true,
		"data":    jobsList,
	}

	jsonResponse, _ := json.Marshal(responseBody)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
