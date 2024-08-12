package saver

import (
	"encoding/json"
	"gian/db"
	"gian/models"
	"gian/utils"
	"net/http"

	"gorm.io/gorm"
)

func SaveJob(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("user_id")
	jobUUID := r.FormValue("job_id")

	// Check if necessary values are present
	if userID == "" || jobUUID == "" {
		utils.HandleError(w, http.StatusBadRequest, "Missing user_id or job_id", nil)
		return
	}

	var saved models.Saved

	// Try to find the existing user in the Saved table
	err := db.DB.Where("user_id = ?", userID).First(&saved).Error
	if err == gorm.ErrRecordNotFound {
		// User is not found, so create a new record
		jobIDs := []string{jobUUID}
		jobIDsJSON, _ := json.Marshal(jobIDs)
		saved = models.Saved{
			UserID: userID,
			JobIDs: string(jobIDsJSON),
		}
		if err := db.DB.Create(&saved).Error; err != nil {
			utils.HandleError(w, http.StatusInternalServerError, "Failed to save job", err)
			return
		}
	} else if err != nil {
		// Any other error
		utils.HandleError(w, http.StatusInternalServerError, "Failed to find user", err)
		return
	} else {
		// User found, update the JobIDs
		var jobIDs []string
		if err := json.Unmarshal([]byte(saved.JobIDs), &jobIDs); err != nil {
			utils.HandleError(w, http.StatusInternalServerError, "Failed to parse job IDs", err)
			return
		}

		// Check if jobUUID is already saved
		for _, id := range jobIDs {
			if id == jobUUID {
				responseBody := map[string]interface{}{
					"message": "Job already saved",
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
		saved.JobIDs = string(jobIDsJSON)
		if err := db.DB.Save(&saved).Error; err != nil {
			utils.HandleError(w, http.StatusInternalServerError, "Failed to save job", err)
			return
		}
	}

	responseBody := map[string]interface{}{
		"message": "Job saved successfully",
		"status":  true,
	}

	jsonResponse, _ := json.Marshal(responseBody)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}

// get list of saved jobs
func GetSavedJobs(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("user_id")

	// Check if the user_uuid is provided
	if userID == "" {
		utils.HandleError(w, http.StatusBadRequest, "Missing user_uuid", nil)
		return
	}

	var saved models.Saved

	// Find the saved jobs for the user
	if err := db.DB.Where("user_id = ?", userID).First(&saved).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.HandleError(w, http.StatusNotFound, "No saved jobs found for the user", err)
		} else {
			utils.HandleError(w, http.StatusInternalServerError, "Failed to find saved jobs", err)
		}
		return
	}

	// Decode the JobIDs from JSON
	var jobIDs []string
	if err := json.Unmarshal([]byte(saved.JobIDs), &jobIDs); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to parse job IDs", err)
		return
	}

	// Fetch the job details for each job UUID in the JobIDs array
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
