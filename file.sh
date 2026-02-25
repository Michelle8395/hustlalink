#!/bin/bash

# Script: create_project_structure.sh
# Purpose: Generate folder and file structure for Youth Employment PWA

# Base folder

echo "Creating project structure in ./...."

# Create main folders
mkdir -p ./frontend/public
mkdir -p ./frontend/src/components
mkdir -p ./frontend/src/pages
mkdir -p ./frontend/src/hooks
mkdir -p ./backend/db
mkdir -p ./backend/controllers
mkdir -p ./backend/routes
mkdir -p ./ussd

# Create placeholder files for frontend
touch ./frontend/public/index.html
touch ./frontend/public/manifest.json
touch ./frontend/src/db.js
touch ./frontend/src/App.jsx
touch ./frontend/src/main.jsx
touch ./frontend/package.json
touch ./frontend/vite.config.js

# Create placeholder files for backend
touch ./backend/db/schema.sql
touch ./backend/db/seed.sql
touch ./backend/controllers/.gitkeep
touch ./backend/routes/.gitkeep
touch ./backend/main.go
touch ./backend/go.mod

# Create placeholder files for USSD
touch ./ussd/main.go
touch ./ussd/routes.go

# Create project-level files
touch ./README.md
touch ./.gitignore

echo "Project structure created successfully" 
