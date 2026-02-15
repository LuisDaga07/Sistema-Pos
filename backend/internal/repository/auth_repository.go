package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pos-saas/restaurant-pos/internal/errors"
	"github.com/pos-saas/restaurant-pos/internal/models"
)

type AuthRepository struct {
	pool *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{pool: pool}
}

func (r *AuthRepository) CreateRestaurant(ctx context.Context, rest *models.Restaurant) error {
	query := `
		INSERT INTO restaurants (id, name, email, phone, address, tax_id, logo_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.pool.Exec(ctx, query,
		rest.ID, rest.Name, rest.Email, rest.Phone, rest.Address, rest.TaxID, rest.LogoURL,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return errors.ErrConflict
		}
		return err
	}
	return nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, restaurant_id, email, password_hash, role, active)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query,
		user.ID, user.RestaurantID, user.Email, user.PasswordHash, user.Role, user.Active,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return errors.ErrConflict
		}
		return err
	}
	return nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, restaurantID uuid.UUID, email string) (*models.User, error) {
	query := `
		SELECT id, restaurant_id, email, password_hash, role, active, created_at, updated_at
		FROM users
		WHERE restaurant_id = $1 AND LOWER(email) = LOWER($2) AND deleted_at IS NULL
	`
	var user models.User
	err := r.pool.QueryRow(ctx, query, restaurantID, email).Scan(
		&user.ID, &user.RestaurantID, &user.Email, &user.PasswordHash,
		&user.Role, &user.Active, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if isNoRows(err) {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) GetRestaurantByEmail(ctx context.Context, email string) (*models.Restaurant, error) {
	query := `
		SELECT id, name, email, phone, address, tax_id, logo_url, created_at, updated_at
		FROM restaurants
		WHERE LOWER(email) = LOWER($1) AND deleted_at IS NULL
	`
	var rest models.Restaurant
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&rest.ID, &rest.Name, &rest.Email, &rest.Phone, &rest.Address,
		&rest.TaxID, &rest.LogoURL, &rest.CreatedAt, &rest.UpdatedAt,
	)
	if err != nil {
		if isNoRows(err) {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return &rest, nil
}

func (r *AuthRepository) GetRestaurantByID(ctx context.Context, id uuid.UUID) (*models.Restaurant, error) {
	query := `
		SELECT id, name, email, phone, address, tax_id, logo_url, created_at, updated_at
		FROM restaurants
		WHERE id = $1 AND deleted_at IS NULL
	`
	var rest models.Restaurant
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&rest.ID, &rest.Name, &rest.Email, &rest.Phone, &rest.Address,
		&rest.TaxID, &rest.LogoURL, &rest.CreatedAt, &rest.UpdatedAt,
	)
	if err != nil {
		if isNoRows(err) {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return &rest, nil
}
