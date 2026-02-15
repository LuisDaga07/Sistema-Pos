package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pos-saas/restaurant-pos/internal/errors"
	"github.com/pos-saas/restaurant-pos/internal/models"
)

type CategoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{pool: pool}
}

func (r *CategoryRepository) Create(ctx context.Context, c *models.Category) error {
	query := `INSERT INTO categories (id, restaurant_id, name, description, sort_order) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.pool.Exec(ctx, query, c.ID, c.RestaurantID, c.Name, c.Description, c.SortOrder)
	return err
}

func (r *CategoryRepository) List(ctx context.Context, restaurantID uuid.UUID) ([]*models.Category, error) {
	query := `
		SELECT id, restaurant_id, name, description, sort_order
		FROM categories
		WHERE restaurant_id = $1
		ORDER BY sort_order, name
	`
	rows, err := r.pool.Query(ctx, query, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.RestaurantID, &cat.Name, &cat.Description, &cat.SortOrder); err != nil {
			return nil, err
		}
		categories = append(categories, &cat)
	}
	return categories, rows.Err()
}

func (r *CategoryRepository) GetByID(ctx context.Context, restaurantID, categoryID uuid.UUID) (*models.Category, error) {
	query := `SELECT id, restaurant_id, name, description, sort_order FROM categories WHERE id = $1 AND restaurant_id = $2`
	var cat models.Category
	err := r.pool.QueryRow(ctx, query, categoryID, restaurantID).Scan(
		&cat.ID, &cat.RestaurantID, &cat.Name, &cat.Description, &cat.SortOrder,
	)
	if err != nil {
		if isNoRows(err) {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return &cat, nil
}
