package handler

import (
	"auth-service/internal/dto"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"auth-service/pkg/jwt"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler обрабатывает HTTP запросы аутентификации
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler создаёт новый экземпляр обработчика аутентификации
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register обрабатывает запрос на регистрацию нового пользователя
// @Summary Регистрация пользователя
// @Description Создаёт нового пользователя и возвращает JWT токены
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Данные регистрации"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	// Парсинг и валидация запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Ошибка валидации",
			Message: err.Error(),
		})
		return
	}

	// Вызов сервиса регистрации
	response, err := h.authService.Register(&req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login обрабатывает запрос на аутентификацию
// @Summary Вход в систему
// @Description Аутентифицирует пользователя и возвращает JWT токены
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Данные для входа"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	// Парсинг и валидация запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Ошибка валидации",
			Message: err.Error(),
		})
		return
	}

	// Вызов сервиса аутентификации
	response, err := h.authService.Login(&req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProfile возвращает профиль текущего пользователя
// @Summary Получение профиля
// @Description Возвращает профиль текущего аутентифицированного пользователя
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /auth/me [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Получение ID пользователя из контекста (установлен middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Не авторизован",
			Message: "Требуется аутентификация",
		})
		return
	}

	// Преобразование ID в UUID
	id, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Внутренняя ошибка",
			Message: "Некорректный формат ID пользователя",
		})
		return
	}

	// Получение профиля
	response, err := h.authService.GetProfile(id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken обновляет JWT токены
// @Summary Обновление токена
// @Description Обновляет access токен используя refresh токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshRequest true "Refresh токен"
// @Success 200 {object} dto.TokenResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshRequest

	// Парсинг и валидация запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Ошибка валидации",
			Message: err.Error(),
		})
		return
	}

	// Вызов сервиса обновления токена
	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// handleServiceError обрабатывает ошибки сервиса и возвращает соответствующий HTTP ответ
func handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrUserAlreadyExists):
		c.JSON(http.StatusConflict, dto.ErrorResponse{
			Error:   "Конфликт",
			Message: "Пользователь с таким email уже существует",
		})
	case errors.Is(err, repository.ErrUserNotFound):
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "Не найдено",
			Message: "Пользователь не найден",
		})
	case errors.Is(err, service.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Ошибка аутентификации",
			Message: "Неверный email или пароль",
		})
	case errors.Is(err, service.ErrUserNotActive):
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error:   "Доступ запрещён",
			Message: "Учётная запись деактивирована",
		})
	case errors.Is(err, service.ErrInvalidRole):
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Ошибка валидации",
			Message: "Недопустимая роль пользователя",
		})
	case errors.Is(err, jwt.ErrInvalidToken):
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Ошибка аутентификации",
			Message: "Недействительный токен",
		})
	case errors.Is(err, jwt.ErrExpiredToken):
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Ошибка аутентификации",
			Message: "Срок действия токена истёк",
		})
	default:
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Внутренняя ошибка сервера",
			Message: "Произошла непредвиденная ошибка",
		})
	}
}
