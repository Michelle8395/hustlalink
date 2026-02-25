package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"

	"hustlalink/backend/controllers"
	"hustlalink/backend/routes"
	"hustlalink/backend/sms"
)

var DB *sql.DB

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// SQLite database
	dbPath := getEnv("DB_PATH", "./hustlalink.db")

	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer DB.Close()

	// Enable WAL mode and foreign keys for SQLite
	DB.Exec("PRAGMA journal_mode=WAL")
	DB.Exec("PRAGMA foreign_keys=ON")

	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Auto-create tables from schema
	initDB(DB)

	log.Println("Connected to SQLite database:", dbPath)

	// Initialize controllers with DB
	controllers.InitDB(DB)

	// Initialize SMS service
	sms.Init()

	// Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Register routes
	routes.RegisterRoutes(r)

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func initDB(db *sql.DB) {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		phone TEXT,
		password TEXT NOT NULL,
		skills TEXT,
		role TEXT NOT NULL CHECK(role IN ('jobseeker', 'employer')),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		employer_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		skills TEXT,
		salary TEXT,
		location TEXT,
		status TEXT DEFAULT 'open' CHECK(status IN ('open', 'closed')),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (employer_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS applications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		job_id INTEGER NOT NULL,
		jobseeker_id INTEGER NOT NULL,
		status TEXT DEFAULT 'pending' CHECK(status IN ('pending', 'accepted', 'rejected')),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE,
		FOREIGN KEY (jobseeker_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS ratings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		employer_id INTEGER NOT NULL,
		seeker_id INTEGER NOT NULL,
		score INTEGER NOT NULL CHECK(score BETWEEN 1 AND 5),
		comment TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (employer_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (seeker_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`
	_, err := db.Exec(schema)
	if err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}
}
