package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pos-saas/restaurant-pos/internal/errors"
	"github.com/pos-saas/restaurant-pos/internal/models"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

func (r *ProductRepository) Create(ctx context.Context, p *models.Product) error {
	query := `
		INSERT INTO products (id, restaurant_id, category_id, name, description, price, image_url, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.pool.Exec(ctx, query,
		p.ID, p.RestaurantID, p.CategoryID, p.Name, p.Description,
		p.Price, p.ImageURL, p.Active,
	)
	return err
}

func (r *ProductRepository) GetByID(ctx context.Context, restaurantID, productID uuid.UUID) (*models.Product, error) {
	query := `
		SELECT id, restaurant_id, category_id, name, description, price, image_url, active, created_at, updated_at
		FROM products
		WHERE id = $1 AND restaurant_id = $2
	`
	var p models.Product
	err := r.pool.QueryRow(ctx, query, productID, restaurantID).Scan(
		&p.ID, &p.RestaurantID, &p.CategoryID, &p.Name, &p.Description,
		&p.Price, &p.ImageURL, &p.Active, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if isNoRows(err) {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) List(ctx context.Context, restaurantID uuid.UUID, categoryID *uuid.UUID, activeOnly bool) ([]*models.Product, error) {
	query := `
		SELECT id, restaurant_id, category_id, name, description, price, image_url, active, created_at, updated_at
		FROM products
		WHERE restaurant_id = $1
	`
	args := []interface{}{restaurantID}
	argNum := 2

	if categoryID != nil {
		query += fmt.Sprintf(" AND category_id = $%d", argNum)
		args = append(args, *categoryID)
		argNum++
	}
	if activeOnly {
		query += " AND active = true"
	}
	query += " ORDER BY name"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.RestaurantID, &p.CategoryID, &p.Name, &p.Description,
			&p.Price, &p.ImageURL, &p.Active, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, rows.Err()
}

func (r *ProductRepository) Update(ctx context.Context, p *models.Product) error {
	query := `
		UPDATE products
		SET category_id = $2, name = $3, description = $4, price = $5, image_url = $6, active = $7
		WHERE id = $1 AND restaurant_id = $8
	`
	result, err := r.pool.Exec(ctx, query,
		p.ID, p.CategoryID, p.Name, p.Description, p.Price, p.ImageURL, p.Active, p.RestaurantID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, restaurantID, productID uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1 AND restaurant_id = $2`
	result, err := r.pool.Exec(ctx, query, productID, restaurantID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.ErrNotFound
	}
	return nil
}
