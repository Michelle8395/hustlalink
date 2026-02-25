package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"hustlalink/backend/sms"
)

// ListJobs returns all open jobs with optional search filters
func ListJobs(c *gin.Context) {
	keyword := c.Query("keyword")
	location := c.Query("location")

	query := `SELECT j.id, j.employer_id, j.title, j.description, j.skills, j.salary, j.location, j.status, j.created_at, u.username as employer_name
		FROM jobs j JOIN users u ON j.employer_id = u.id WHERE j.status = 'open'`
	args := []interface{}{}

	if keyword != "" {
		query += " AND (j.title LIKE ? OR j.description LIKE ? OR j.skills LIKE ?)"
		k := "%" + keyword + "%"
		args = append(args, k, k, k)
	}
	if location != "" {
		query += " AND j.location LIKE ?"
		args = append(args, "%"+location+"%")
	}

	query += " ORDER BY j.created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
		return
	}
	defer rows.Close()

	var jobs []gin.H
	for rows.Next() {
		var j struct {
			ID           int
			EmployerID   int
			Title        string
			Description  sql.NullString
			Skills       sql.NullString
			Salary       sql.NullString
			Location     sql.NullString
			Status       string
			CreatedAt    string
			EmployerName string
		}
		if err := rows.Scan(&j.ID, &j.EmployerID, &j.Title, &j.Description, &j.Skills, &j.Salary, &j.Location, &j.Status, &j.CreatedAt, &j.EmployerName); err != nil {
			continue
		}
		jobs = append(jobs, gin.H{
			"id":            j.ID,
			"employer_id":   j.EmployerID,
			"title":         j.Title,
			"description":   j.Description.String,
			"skills":        j.Skills.String,
			"salary":        j.Salary.String,
			"location":      j.Location.String,
			"status":        j.Status,
			"created_at":    j.CreatedAt,
			"employer_name": j.EmployerName,
		})
	}

	if jobs == nil {
		jobs = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{"jobs": jobs})
}

// GetJob returns a single job by ID
func GetJob(c *gin.Context) {
	id := c.Param("id")

	var j struct {
		ID           int
		EmployerID   int
		Title        string
		Description  sql.NullString
		Skills       sql.NullString
		Salary       sql.NullString
		Location     sql.NullString
		Status       string
		CreatedAt    string
		EmployerName string
	}

	err := db.QueryRow(
		`SELECT j.id, j.employer_id, j.title, j.description, j.skills, j.salary, j.location, j.status, j.created_at, u.username as employer_name
		FROM jobs j JOIN users u ON j.employer_id = u.id WHERE j.id = ?`, id,
	).Scan(&j.ID, &j.EmployerID, &j.Title, &j.Description, &j.Skills, &j.Salary, &j.Location, &j.Status, &j.CreatedAt, &j.EmployerName)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"job": gin.H{
			"id":            j.ID,
			"employer_id":   j.EmployerID,
			"title":         j.Title,
			"description":   j.Description.String,
			"skills":        j.Skills.String,
			"salary":        j.Salary.String,
			"location":      j.Location.String,
			"status":        j.Status,
			"created_at":    j.CreatedAt,
			"employer_name": j.EmployerName,
		},
	})
}

type CreateJobInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Skills      string `json:"skills"`
	Salary      string `json:"salary"`
	Location    string `json:"location"`
}

// CreateJob allows an employer to post a new job
func CreateJob(c *gin.Context) {
	claims, ok := c.MustGet("claims").(*UserClaims)
	if !ok || claims.Role != "employer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only employers can post jobs"})
		return
	}

	var input CreateJobInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(
		"INSERT INTO jobs (employer_id, title, description, skills, salary, location) VALUES (?, ?, ?, ?, ?, ?)",
		claims.UserID, input.Title, input.Description, input.Skills, input.Salary, input.Location,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}

	jobID, _ := result.LastInsertId()

	// Notify job seekers via SMS
	go notifyMatchingJobSeekers(input.Title, input.Location)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Job posted successfully",
		"job_id":  jobID,
	})
}

// notifyMatchingJobSeekers sends SMS to seekers about new opportunities
func notifyMatchingJobSeekers(jobTitle, jobLocation string) {
	rows, err := db.Query(
		`SELECT phone, username FROM users WHERE role = 'jobseeker' AND phone IS NOT NULL AND phone != ''`,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var phone, username string
		if err := rows.Scan(&phone, &username); err != nil {
			continue
		}
		msg := fmt.Sprintf("Hi %s! New vocational job on HustlaLink: %s", username, jobTitle)
		if jobLocation != "" {
			msg += fmt.Sprintf(" in %s", jobLocation)
		}
		msg += ". Check it out: hustlalink.co.ke"
		sms.SendSMS(phone, msg)
	}
}

// UpdateJob allows an employer to update their job
func UpdateJob(c *gin.Context) {
	claims := c.MustGet("claims").(*UserClaims)
	jobID := c.Param("id")

	// Verify ownership
	var ownerID int
	err := db.QueryRow("SELECT employer_id FROM jobs WHERE id = ?", jobID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	if ownerID != claims.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only edit your own jobs"})
		return
	}

	var input CreateJobInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = db.Exec(
		"UPDATE jobs SET title=?, description=?, skills=?, salary=?, location=?, updated_at=CURRENT_TIMESTAMP WHERE id=?",
		input.Title, input.Description, input.Skills, input.Salary, input.Location, jobID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update job"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job updated successfully"})
}

// DeleteJob allows an employer to delete their job
func DeleteJob(c *gin.Context) {
	claims := c.MustGet("claims").(*UserClaims)
	jobID := c.Param("id")

	var ownerID int
	err := db.QueryRow("SELECT employer_id FROM jobs WHERE id = ?", jobID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	if ownerID != claims.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own jobs"})
		return
	}

	_, err = db.Exec("DELETE FROM jobs WHERE id = ?", jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete job"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job deleted successfully"})
}

// GetEmployerJobs returns all jobs posted by the authenticated employer
func GetEmployerJobs(c *gin.Context) {
	claims := c.MustGet("claims").(*UserClaims)

	rows, err := db.Query(
		`SELECT j.id, j.title, j.location, j.salary, j.status, j.created_at,
		(SELECT COUNT(*) FROM applications WHERE job_id = j.id) as applicant_count
		FROM jobs j WHERE j.employer_id = ? ORDER BY j.created_at DESC`,
		claims.UserID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
		return
	}
	defer rows.Close()

	var jobs []gin.H
	for rows.Next() {
		var id int
		var title, status, createdAt string
		var locationN, salaryN sql.NullString
		var applicantCount int
		if err := rows.Scan(&id, &title, &locationN, &salaryN, &status, &createdAt, &applicantCount); err != nil {
			continue
		}
		jobs = append(jobs, gin.H{
			"id":              id,
			"title":           title,
			"location":        locationN.String,
			"salary":          salaryN.String,
			"status":          status,
			"created_at":      createdAt,
			"applicant_count": applicantCount,
		})
	}

	if jobs == nil {
		jobs = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{"jobs": jobs})
}
