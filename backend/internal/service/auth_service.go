package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pos-saas/restaurant-pos/internal/errors"
	"github.com/pos-saas/restaurant-pos/internal/middleware"
	"github.com/pos-saas/restaurant-pos/internal/models"
	"github.com/pos-saas/restaurant-pos/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo         *repository.AuthRepository
	jwtSecret    string
	jwtExpHours  int
}

func NewAuthService(repo *repository.AuthRepository, jwtSecret string, jwtExpHours int) *AuthService {
	return &AuthService{
		repo:        repo,
		jwtSecret:   jwtSecret,
		jwtExpHours: jwtExpHours,
	}
}

type RegisterInput struct {
	RestaurantName string `json:"restaurant_name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=6"`
	Phone          string `json:"phone"`
	Address        string `json:"address"`
	TaxID          string `json:"tax_id"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token        string           `json:"token"`
	ExpiresAt    time.Time        `json:"expires_at"`
	User         UserResponse     `json:"user"`
	Restaurant   RestaurantResponse `json:"restaurant"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type RestaurantResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
	// Verificar si el email del restaurante ya existe
	_, err := s.repo.GetRestaurantByEmail(ctx, input.Email)
	if err == nil {
		return nil, NewAppError(errors.ErrConflict, 409, "el restaurante ya está registrado con ese email")
	}
	if !errors.Is(err, errors.ErrNotFound) {
		return nil, err
	}

	// Hash de la contraseña
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	restID := uuid.New()
	userID := uuid.New()

	restaurant := &models.Restaurant{
		ID:      restID,
		Name:    input.RestaurantName,
		Email:   input.Email,
		Phone:   input.Phone,
		Address: input.Address,
		TaxID:   input.TaxID,
	}

	user := &models.User{
		ID:           userID,
		RestaurantID: restID,
		Email:        input.Email,
		PasswordHash: string(hash),
		Role:         "admin",
		Active:       true,
	}

	if err := s.repo.CreateRestaurant(ctx, restaurant); err != nil {
		return nil, err
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	token, exp, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:     token,
		ExpiresAt: exp,
		User: UserResponse{
			ID:    user.ID.String(),
			Email: user.Email,
			Role:  user.Role,
		},
		Restaurant: RestaurantResponse{
			ID:    restaurant.ID.String(),
			Name:  restaurant.Name,
			Email: restaurant.Email,
		},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	// Buscar restaurante por email (el login usa email del restaurante como identificador)
	restaurant, err := s.repo.GetRestaurantByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, NewAppError(errors.ErrInvalidCredentials, 401, "credenciales inválidas")
		}
		return nil, err
	}

	user, err := s.repo.GetUserByEmail(ctx, restaurant.ID, input.Email)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, NewAppError(errors.ErrInvalidCredentials, 401, "credenciales inválidas")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, NewAppError(errors.ErrInvalidCredentials, 401, "credenciales inválidas")
	}

	if !user.Active {
		return nil, NewAppError(errors.ErrForbidden, 403, "usuario inactivo")
	}

	token, exp, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:     token,
		ExpiresAt: exp,
		User: UserResponse{
			ID:    user.ID.String(),
			Email: user.Email,
			Role:  user.Role,
		},
		Restaurant: RestaurantResponse{
			ID:    restaurant.ID.String(),
			Name:  restaurant.Name,
			Email: restaurant.Email,
		},
	}, nil
}

func (s *AuthService) generateToken(user *models.User) (string, time.Time, error) {
	exp := time.Now().Add(time.Duration(s.jwtExpHours) * time.Hour)
	claims := &middleware.Claims{
		UserID:       user.ID.String(),
		RestaurantID: user.RestaurantID.String(),
		Email:        user.Email,
		Role:         user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, exp, nil
}
