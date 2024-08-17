package jobs

import (
	"encoding/json"
	"gian/db"
	"gian/models"
	"gian/utils"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateJob(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("user_id")
	owner := r.Header.Get("uuid")

	// Extract form values
	status := r.FormValue("status") // active, hired, closed etc
	companyName := r.FormValue("company_name")
	position := r.FormValue("position")
	location := r.FormValue("location")
	jobType := r.FormValue("job_type") // fulltime or parttime
	description := r.FormValue("description")
	field := r.FormValue("field")
	minPay := r.FormValue("min_pay")
	maxPay := r.FormValue("max_pay")
	price := r.FormValue("price")
	totalEmp := r.FormValue("total_emp")
	logoUrl := r.FormValue("logo_url")

	// Check if necessary form values are present
	if userID == "" || status == "" || companyName == "" || position == "" || location == "" || description == "" || jobType == "" || field == "" || owner == "" || minPay == "" || maxPay == "" || price == "" || totalEmp == "" || logoUrl == "" {
		utils.HandleError(w, http.StatusBadRequest, "Missing necessary details", nil)
		return
	}

	// Convert minPay, maxPay, price, and totalEmp to int
	minPayInt, err := strconv.Atoi(minPay)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid value for min_pay", err)
		return
	}
	maxPayInt, err := strconv.Atoi(maxPay)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid value for max_pay", err)
		return
	}
	priceInt, err := strconv.Atoi(price)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid value for price", err)
		return
	}
	totalEmpInt, err := strconv.Atoi(totalEmp)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid value for total_emp", err)
		return
	}

	var job models.Jobs

	// Assign values to job fields
	job.UserID, _ = strconv.Atoi(userID)
	job.Status = status
	job.Uuid = uuid.New().String()
	job.CompanyName = companyName
	job.Position = position
	job.Location = location
	job.Type = jobType
	job.Description = description
	job.Field = field
	job.Owner = owner
	job.MinPay = minPayInt
	job.MaxPay = maxPayInt
	job.Price = priceInt
	job.TotalEmp = totalEmpInt
	job.LogoUrl = logoUrl

	// Save the job to the database
	if err := db.DB.Create(&job).Error; err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to create job", err)
		return
	}

	responseBody := map[string]interface{}{
		"message": "Job created successfully",
		"job":     job,
		"status":  true,
	}

	jsonResponse, _ := json.Marshal(responseBody)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}

func GetJobsList(w http.ResponseWriter, r *http.Request) {

	var jobsList []models.Jobs

	err := db.DB.Find(&jobsList).Error

	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Unable to process request", err)
		return
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

func GetJobs(w http.ResponseWriter, r *http.Request) {
	// Retrieve role, id (employee ID), and company code from the headers
	owneruuid := r.Header.Get("uuid")

	// Check if the necessary headers are present

	var jobs []models.Jobs
	var err error

	// Fetch orders based on role

	if owneruuid == "" {
		utils.HandleError(w, http.StatusBadRequest, "Missing employee ID in headers", nil)
		return
	}

	err = db.DB.Where("owner = ?", owneruuid).Find(&jobs).Error

	// Handle errors during fetching
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Error fetching jobs", err)
		return
	}

	responseBody := map[string]interface{}{
		"message": "success",
		"orders":  jobs,
	}

	jsonResponse, _ := json.Marshal(responseBody)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func UpdateJobStatus(w http.ResponseWriter, r *http.Request) {

	// Get the status and ID from the form data
	status := r.FormValue("status")
	id := r.FormValue("id")

	if id == "" || status == "" {
		responseBody := map[string]interface{}{
			"message":   "Missing status or ID",
			"status":    false,
			"id":        id,
			"newStatus": status,
		}
		jsonResponse, _ := json.Marshal(responseBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return
	}

	// Find the existing job by ID
	var job models.Jobs
	if err := db.DB.First(&job, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "job not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch job", http.StatusInternalServerError)
		}
		return
	}

	// Update the status
	job.Status = status
	if err := db.DB.Save(&job).Error; err != nil {
		http.Error(w, "Failed to update job status", http.StatusInternalServerError)
		return
	}

	responseBody := map[string]interface{}{
		"message": "Status updated successfully",
		"status":  true,
	}

	jsonResponse, _ := json.Marshal(responseBody)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}
