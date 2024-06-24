package models

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// User represents the user model
type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Designation  string    `json:"designation,omitempty"`
	Bio          string    `json:"bio,omitempty"`
	ProfileImage string    `json:"profile_image,omitempty"`
	Location     string    `json:"location,omitempty"`
}

// ConnectDatabase establishes a connection to the MySQL database
func ConnectDatabase() (*sql.DB, error) {
	dbUser := os.Getenv("DB_USER_AIVEN")
	dbPassword := os.Getenv("DB_PASSWORD_AIVEN")
	dbHost := os.Getenv("DB_HOST_AIVEN")
	dbPort := os.Getenv("DB_PORT_AIVEN")
	dbName := os.Getenv("DB_NAME_AIVEN")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
