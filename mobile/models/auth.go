package models

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthUser struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type AuthResponse struct {
	Token string   `json:"token"`
	User  AuthUser `json:"user"`
}
