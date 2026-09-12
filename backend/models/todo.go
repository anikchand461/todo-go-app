package models

import "time"

type Todo struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

type UpdateTodoRequest struct {
	Title     *string `json:"title,omitempty"`
	Completed *bool   `json:"completed,omitempty"`
}
