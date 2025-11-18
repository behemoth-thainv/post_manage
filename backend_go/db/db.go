package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB is the global database handle for convenience in this small project.
// For larger apps, prefer passing *gorm.DB explicitly to constructors.
var DB *gorm.DB

// getenv returns value or default when empty
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// Init initializes the GORM DB using environment variables.
// Supported env vars: DATABASE_USER, DATABASE_PASSWORD, DATABASE_HOST, DATABASE_PORT, DATABASE_NAME
func Init() {
	user := getenv("DATABASE_USER", "root")
	pass := getenv("DATABASE_PASSWORD", "")
	host := getenv("DATABASE_HOST", "127.0.0.1")
	port := getenv("DATABASE_PORT", "3306")
	name := getenv("DATABASE_NAME", "sample_crud")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
}
