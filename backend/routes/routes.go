package routes

import (
	"hustlalink/backend/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all API route groups
func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "hustlalink-api"})
	})

	// Auth routes (public)
	auth := api.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// Job routes
	jobs := api.Group("/jobs")
	{
		jobs.GET("", controllers.ListJobs)   // Public: browse/search jobs
		jobs.GET("/:id", controllers.GetJob) // Public: job detail
	}

	// Protected job routes (employer)
	protectedJobs := api.Group("/jobs")
	protectedJobs.Use(controllers.AuthMiddleware())
	{
		protectedJobs.POST("", controllers.CreateJob)
		protectedJobs.PUT("/:id", controllers.UpdateJob)
		protectedJobs.DELETE("/:id", controllers.DeleteJob)
	}

	// Employer-specific routes
	employer := api.Group("/employer")
	employer.Use(controllers.AuthMiddleware())
	{
		employer.GET("/jobs", controllers.GetEmployerJobs)
		employer.GET("/applications/:id", controllers.GetJobApplications)
		employer.PUT("/applications/:id/status", controllers.UpdateApplicationStatus)
	}

	// Application routes (protected)
	applications := api.Group("/applications")
	applications.Use(controllers.AuthMiddleware())
	{
		applications.POST("", controllers.ApplyForJob)
		applications.GET("/mine", controllers.GetMyApplications)
	}

	// Rating routes
	ratings := api.Group("/ratings")
	{
		ratings.GET("/seeker/:id", controllers.GetSeekerRatings) // Public: view a worker's ratings
	}

	protectedRatings := api.Group("/ratings")
	protectedRatings.Use(controllers.AuthMiddleware())
	{
		protectedRatings.POST("", controllers.AddRating) // Employer: rate a worker
	}
}
