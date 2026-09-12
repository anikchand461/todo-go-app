package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"todo-go/backend/auth"
	"todo-go/backend/models"
)

type TodoHandler struct {
	DB *pgxpool.Pool
}

func NewTodoHandler(db *pgxpool.Pool) *TodoHandler {
	return &TodoHandler{
		DB: db,
	}
}

// GET /api/todos
func (h *TodoHandler) GetTodos(c *gin.Context) {
	userID, err := auth.GetUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	rows, err := h.DB.Query(
		context.Background(),
		`
		SELECT id, user_id, title, completed, created_at, updated_at
		FROM todos
		WHERE user_id = $1
		ORDER BY created_at DESC
		`,
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch todos",
		})
		return
	}

	defer rows.Close()

	todos := make([]models.Todo, 0)

	for rows.Next() {
		var todo models.Todo

		err := rows.Scan(
			&todo.ID,
			&todo.UserID,
			&todo.Title,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to read todo",
			})
			return
		}

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed while reading todos",
		})
		return
	}

	c.JSON(http.StatusOK, todos)
}

// GET /api/todos/:id
func (h *TodoHandler) GetTodo(c *gin.Context) {
	userID, err := auth.GetUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid todo ID",
		})
		return
	}

	var todo models.Todo

	err = h.DB.QueryRow(
		context.Background(),
		`
		SELECT id, user_id, title, completed, created_at, updated_at
		FROM todos
		WHERE id = $1 AND user_id = $2
		`,
		id,
		userID,
	).Scan(
		&todo.ID,
		&todo.UserID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Todo not found",
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// POST /api/todos
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	userID, err := auth.GetUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var request models.CreateTodoRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	request.Title = strings.TrimSpace(request.Title)

	if request.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Title cannot be empty",
		})
		return
	}

	var todo models.Todo

	err = h.DB.QueryRow(
		context.Background(),
		`
		INSERT INTO todos (user_id, title)
		VALUES ($1, $2)
		RETURNING id, user_id, title, completed, created_at, updated_at
		`,
		userID,
		request.Title,
	).Scan(
		&todo.ID,
		&todo.UserID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create todo",
		})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

// PUT /api/todos/:id
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	userID, err := auth.GetUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid todo ID",
		})
		return
	}

	var request models.UpdateTodoRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if request.Title == nil && request.Completed == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nothing to update",
		})
		return
	}

	if request.Title != nil {
		*request.Title = strings.TrimSpace(*request.Title)

		if *request.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Title cannot be empty",
			})
			return
		}
	}

	var todo models.Todo

	err = h.DB.QueryRow(
		context.Background(),
		`
		UPDATE todos
		SET
			title = COALESCE($1, title),
			completed = COALESCE($2, completed),
			updated_at = NOW()
		WHERE id = $3 AND user_id = $4
		RETURNING id, user_id, title, completed, created_at, updated_at
		`,
		request.Title,
		request.Completed,
		id,
		userID,
	).Scan(
		&todo.ID,
		&todo.UserID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Todo not found",
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// DELETE /api/todos/:id
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	userID, err := auth.GetUserID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid todo ID",
		})
		return
	}

	commandTag, err := h.DB.Exec(
		context.Background(),
		`
		DELETE FROM todos
		WHERE id = $1 AND user_id = $2
		`,
		id,
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete todo",
		})
		return
	}

	if commandTag.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Todo not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Todo deleted successfully",
	})
}
