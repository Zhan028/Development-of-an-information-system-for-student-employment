package dto

import "auth-service/internal/models"

// RegisterRequest представляет запрос на регистрацию нового пользователя
type RegisterRequest struct {
	Email    string          `json:"email" binding:"required,email" example:"user@example.com"`
	Password string          `json:"password" binding:"required,min=8" example:"password123"`
	Role     models.UserRole `json:"role" binding:"required,oneof=student employer university admin" example:"student"`
}

// LoginRequest представляет запрос на аутентификацию
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// RefreshRequest представляет запрос на обновление токена
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
