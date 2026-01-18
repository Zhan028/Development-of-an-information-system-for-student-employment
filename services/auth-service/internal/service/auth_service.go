package service

import (
	"auth-service/internal/dto"
	"auth-service/internal/models"
	"auth-service/internal/repository"
	"auth-service/pkg/jwt"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Ошибки сервиса аутентификации
var (
	ErrInvalidCredentials = errors.New("неверный email или пароль")
	ErrUserNotActive      = errors.New("учётная запись деактивирована")
	ErrInvalidRole        = errors.New("недопустимая роль пользователя")
)

// AuthService определяет интерфейс сервиса аутентификации
type AuthService interface {
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
	GetProfile(userID uuid.UUID) (*dto.UserResponse, error)
	RefreshToken(refreshToken string) (*dto.TokenResponse, error)
}

// authService реализует AuthService
type authService struct {
	userRepo   repository.UserRepository
	jwtManager *jwt.JWTManager
}

// NewAuthService создаёт новый экземпляр сервиса аутентификации
func NewAuthService(userRepo repository.UserRepository, jwtManager *jwt.JWTManager) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Register регистрирует нового пользователя и возвращает JWT токены
func (s *authService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Валидация роли
	if !req.Role.IsValid() {
		return nil, ErrInvalidRole
	}

	// Хеширование пароля с использованием bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создание нового пользователя
	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		IsActive:     true,
	}

	// Сохранение в базе данных
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Генерация JWT токенов
	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(
		user.ID,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	// Формирование ответа
	return &dto.AuthResponse{
		User: dto.ToUserResponse(user),
		Tokens: dto.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    s.jwtManager.GetAccessDuration(),
		},
	}, nil
}

// Login аутентифицирует пользователя и возвращает JWT токены
func (s *authService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Поиск пользователя по email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Проверка активности учётной записи
	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Генерация JWT токенов
	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(
		user.ID,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	// Формирование ответа
	return &dto.AuthResponse{
		User: dto.ToUserResponse(user),
		Tokens: dto.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    s.jwtManager.GetAccessDuration(),
		},
	}, nil
}

// GetProfile возвращает профиль пользователя по ID
func (s *authService) GetProfile(userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	response := dto.ToUserResponse(user)
	return &response, nil
}

// RefreshToken обновляет access токен используя refresh токен
func (s *authService) RefreshToken(refreshToken string) (*dto.TokenResponse, error) {
	// Валидация refresh токена
	claims, err := s.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Проверка существования пользователя
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	// Проверка активности учётной записи
	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Генерация новой пары токенов
	newAccessToken, newRefreshToken, err := s.jwtManager.GenerateTokenPair(
		user.ID,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.jwtManager.GetAccessDuration(),
	}, nil
}
