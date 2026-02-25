package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"hustlalink/backend/sms"
)

// ApplyForJob allows a job seeker to apply for a job
func ApplyForJob(c *gin.Context) {
	claims, ok := c.MustGet("claims").(*UserClaims)
	if !ok || claims.Role != "jobseeker" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only job seekers can apply for jobs"})
		return
	}

	var input struct {
		JobID int `json:"job_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if job exists and is open
	var jobStatus, jobTitle, employerPhone string
	var employerPhoneN sql.NullString
	err := db.QueryRow(
		`SELECT j.status, j.title, u.phone FROM jobs j JOIN users u ON j.employer_id = u.id WHERE j.id = ?`,
		input.JobID,
	).Scan(&jobStatus, &jobTitle, &employerPhoneN)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	if jobStatus != "open" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This job is no longer accepting applications"})
		return
	}
	if employerPhoneN.Valid {
		employerPhone = employerPhoneN.String
	}

	// Check for duplicate application
	var exists int
	db.QueryRow("SELECT COUNT(*) FROM applications WHERE job_id = ? AND jobseeker_id = ?", input.JobID, claims.UserID).Scan(&exists)
	if exists > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "You have already applied for this job"})
		return
	}

	_, err = db.Exec(
		"INSERT INTO applications (job_id, jobseeker_id) VALUES (?, ?)",
		input.JobID, claims.UserID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit application"})
		return
	}

	// Notify employer via SMS
	if employerPhone != "" {
		go sms.SendSMS(employerPhone, fmt.Sprintf("New worker interested in '%s' on HustlaLink. Check apps now.", jobTitle))
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Application submitted successfully"})
}

// GetMyApplications returns the authenticated job seeker's applications
func GetMyApplications(c *gin.Context) {
	claims := c.MustGet("claims").(*UserClaims)

	rows, err := db.Query(
		`SELECT a.id, a.status, a.created_at, j.id, j.title, j.location, j.salary, u.username as employer_name
		FROM applications a
		JOIN jobs j ON a.job_id = j.id
		JOIN users u ON j.employer_id = u.id
		WHERE a.jobseeker_id = ?
		ORDER BY a.created_at DESC`,
		claims.UserID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applications"})
		return
	}
	defer rows.Close()

	var applications []gin.H
	for rows.Next() {
		var appID, jobID int
		var appStatus, appCreatedAt, jobTitle, employerName string
		var jobLocation, jobSalary sql.NullString
		if err := rows.Scan(&appID, &appStatus, &appCreatedAt, &jobID, &jobTitle, &jobLocation, &jobSalary, &employerName); err != nil {
			continue
		}
		applications = append(applications, gin.H{
			"id":            appID,
			"status":        appStatus,
			"created_at":    appCreatedAt,
			"job_id":        jobID,
			"job_title":     jobTitle,
			"job_location":  jobLocation.String,
			"job_salary":    jobSalary.String,
			"employer_name": employerName,
		})
	}

	if applications == nil {
		applications = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{"applications": applications})
}

// GetJobApplications returns applications for a specific job (employer only)
func GetJobApplications(c *gin.Context) {
	claims := c.MustGet("claims").(*UserClaims)
	jobID := c.Param("id")

	// Verify the employer owns this job
	var ownerID int
	err := db.QueryRow("SELECT employer_id FROM jobs WHERE id = ?", jobID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	if ownerID != claims.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	rows, err := db.Query(
		`SELECT a.id, a.status, a.created_at, u.id, u.username, u.email, u.phone, u.skills
		FROM applications a
		JOIN users u ON a.jobseeker_id = u.id
		WHERE a.job_id = ?
		ORDER BY a.created_at DESC`,
		jobID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applications"})
		return
	}
	defer rows.Close()

	var applications []gin.H
	for rows.Next() {
		var appID, userID int
		var appStatus, appCreatedAt, username, email string
		var phone, skills sql.NullString
		if err := rows.Scan(&appID, &appStatus, &appCreatedAt, &userID, &username, &email, &phone, &skills); err != nil {
			continue
		}
		applications = append(applications, gin.H{
			"id":         appID,
			"status":     appStatus,
			"created_at": appCreatedAt,
			"applicant": gin.H{
				"id":       userID,
				"username": username,
				"email":    email,
				"phone":    phone.String,
				"skills":   skills.String,
			},
		})
	}

	if applications == nil {
		applications = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{"applications": applications})
}

// UpdateApplicationStatus allows an employer to accept/reject an application
func UpdateApplicationStatus(c *gin.Context) {
	claims := c.MustGet("claims").(*UserClaims)
	appID := c.Param("id")

	var input struct {
		Status string `json:"status" binding:"required,oneof=accepted rejected"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify the employer owns the job
	var ownerID int
	var jobTitle, seekerPhone, seekerName string
	var seekerPhoneN sql.NullString
	err := db.QueryRow(
		`SELECT j.employer_id, j.title, u.phone, u.username
		FROM applications a
		JOIN jobs j ON a.job_id = j.id
		JOIN users u ON a.jobseeker_id = u.id
		WHERE a.id = ?`,
		appID,
	).Scan(&ownerID, &jobTitle, &seekerPhoneN, &seekerName)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}
	if ownerID != claims.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	if seekerPhoneN.Valid {
		seekerPhone = seekerPhoneN.String
	}

	_, err = db.Exec("UPDATE applications SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", input.Status, appID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update application"})
		return
	}

	// Notify seeker via SMS
	if seekerPhone != "" {
		statusMsg := "ACCEPTED ✅"
		if input.Status == "rejected" {
			statusMsg = "REVIEWED. Check status in app."
		}
		go sms.SendSMS(seekerPhone, fmt.Sprintf("Hi %s, your interest in '%s' was %s. Good luck!", seekerName, jobTitle, statusMsg))
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application status updated"})
}
