package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RatingInput struct {
	SeekerID int    `json:"seeker_id" binding:"required"`
	Score    int    `json:"score" binding:"required,min=1,max=5"`
	Comment  string `json:"comment"`
}

// AddRating allows an employer to rate a jobseeker
func AddRating(c *gin.Context) {
	claims, ok := c.MustGet("claims").(*UserClaims)
	if !ok || claims.Role != "employer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only employers can provide ratings"})
		return
	}

	var input RatingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Optional: Check if the seeker exists
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ? AND role = 'jobseeker'", input.SeekerID).Scan(&exists)
	if err != nil || exists == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jobseeker not found"})
		return
	}

	// Insert rating
	_, err = db.Exec(
		"INSERT INTO ratings (employer_id, seeker_id, score, comment) VALUES (?, ?, ?, ?)",
		claims.UserID, input.SeekerID, input.Score, input.Comment,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit rating"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Rating submitted successfully"})
}

// GetSeekerRatings returns all ratings for a specific seeker and their average score
func GetSeekerRatings(c *gin.Context) {
	seekerID := c.Param("id")

	// Get average score
	var avgScore sql.NullFloat64
	var count int
	err := db.QueryRow("SELECT AVG(score), COUNT(*) FROM ratings WHERE seeker_id = ?", seekerID).Scan(&avgScore, &count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Get individual ratings
	rows, err := db.Query(
		`SELECT r.score, r.comment, r.created_at, u.username as employer_name
		FROM ratings r
		JOIN users u ON r.employer_id = u.id
		WHERE r.seeker_id = ?
		ORDER BY r.created_at DESC`,
		seekerID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ratings"})
		return
	}
	defer rows.Close()

	var ratings []gin.H
	for rows.Next() {
		var score int
		var comment sql.NullString
		var createdAt, employerName string
		if err := rows.Scan(&score, &comment, &createdAt, &employerName); err != nil {
			continue
		}
		ratings = append(ratings, gin.H{
			"score":         score,
			"comment":       comment.String,
			"created_at":    createdAt,
			"employer_name": employerName,
		})
	}

	if ratings == nil {
		ratings = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{
		"seeker_id":     seekerID,
		"average_score": avgScore.Float64,
		"total_ratings": count,
		"ratings":       ratings,
	})
}
