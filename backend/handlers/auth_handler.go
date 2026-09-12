package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"todo-go/backend/auth"
	"todo-go/backend/models"
)

type AuthHandler struct {
	DB *pgxpool.Pool
}

func NewAuthHandler(db *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{
		DB: db,
	}
}

// POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var request models.RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	request.Username = strings.TrimSpace(request.Username)

	if request.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username cannot be empty",
		})
		return
	}

	if len(request.Username) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username must be at least 3 characters",
		})
		return
	}

	if request.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password cannot be empty",
		})
		return
	}

	if len(request.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be at least 6 characters",
		})
		return
	}

	// Check whether username already exists.
	var existingID int

	err := h.DB.QueryRow(
		context.Background(),
		`
		SELECT id
		FROM users
		WHERE username = $1
		`,
		request.Username,
	).Scan(&existingID)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Username already exists",
		})
		return
	}

	// Hash password.
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to secure password",
		})
		return
	}

	// Create user.
	var user models.User

	err = h.DB.QueryRow(
		context.Background(),
		`
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, username, created_at
		`,
		request.Username,
		string(passwordHash),
	).Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Generate JWT.
	token, err := auth.GenerateToken(user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create authentication token",
		})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var request models.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	request.Username = strings.TrimSpace(request.Username)

	var user models.User
	var passwordHash string

	err := h.DB.QueryRow(
		context.Background(),
		`
		SELECT id, username, password_hash, created_at
		FROM users
		WHERE username = $1
		`,
		request.Username,
	).Scan(
		&user.ID,
		&user.Username,
		&passwordHash,
		&user.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	// Check password.
	err = bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(request.Password),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	// Generate JWT.
	token, err := auth.GenerateToken(user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create authentication token",
		})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user,
	})
}
