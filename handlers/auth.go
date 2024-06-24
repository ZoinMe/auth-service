package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/ZoinMe/auth-service/models"
	"github.com/ZoinMe/auth-service/utils"

	"github.com/gin-gonic/gin"
)

// SignupHandler handles user registration
func SignupHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Hash the password before storing it
		// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		// 	return
		// }

		//user.Password = string(hashedPassword)
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		query := `INSERT INTO users (name, email, password, created_at, updated_at, designation, bio, profile_image, location) 
                  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, err := db.Exec(query, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt, user.Designation, user.Bio, user.ProfileImage, user.Location)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	}
}

// LoginHandler handles user login
func LoginHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		query := "SELECT id, password FROM users WHERE email = ?"
		err := db.QueryRow(query, credentials.Email).Scan(&user.ID, &user.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying database"})
			}
			return
		}

		// Assume password check is done here (commented out bcrypt check)

		if user.Password != credentials.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token, err := utils.GenerateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}

		// Return the token and user ID
		c.JSON(http.StatusOK, gin.H{
			"token":   token,
			"user_id": user.ID, // Include the user ID in the response
		})
	}
}
