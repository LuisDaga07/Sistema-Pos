package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pos-saas/restaurant-pos/internal/errors"
	"github.com/pos-saas/restaurant-pos/internal/models"
)

type SaleRepository struct {
	pool *pgxpool.Pool
}

func NewSaleRepository(pool *pgxpool.Pool) *SaleRepository {
	return &SaleRepository{pool: pool}
}

func (r *SaleRepository) Create(ctx context.Context, sale *models.Sale) error {
	query := `
		INSERT INTO sales (id, restaurant_id, user_id, total, status)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query,
		sale.ID, sale.RestaurantID, sale.UserID, sale.Total, sale.Status,
	)
	return err
}

func (r *SaleRepository) CreateItem(ctx context.Context, item *models.SaleItem) error {
	query := `INSERT INTO sale_items (id, sale_id, product_id, quantity, unit_price, subtotal, notes) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, query, item.ID, item.SaleID, item.ProductID, item.Quantity, item.UnitPrice, item.Subtotal, item.Notes)
	return err
}

func (r *SaleRepository) CreateItemTopping(ctx context.Context, topping *models.Topping) error {
	query := `INSERT INTO sale_item_toppings (id, sale_item_id, name, price, quantity) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.pool.Exec(ctx, query, topping.ID, topping.SaleItemID, topping.Name, topping.Price, topping.Quantity)
	return err
}

func (r *SaleRepository) CreatePayment(ctx context.Context, payment *models.SalePayment) error {
	query := `INSERT INTO sale_payments (id, sale_id, method, amount, reference) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.pool.Exec(ctx, query, payment.ID, payment.SaleID, payment.Method, payment.Amount, payment.Reference)
	return err
}

func (r *SaleRepository) GetByID(ctx context.Context, restaurantID, saleID uuid.UUID) (*models.Sale, error) {
	query := `
		SELECT id, restaurant_id, user_id, total, status, created_at, updated_at
		FROM sales
		WHERE id = $1 AND restaurant_id = $2
	`
	var s models.Sale
	err := r.pool.QueryRow(ctx, query, saleID, restaurantID).Scan(
		&s.ID, &s.RestaurantID, &s.UserID, &s.Total, &s.Status, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if isNoRows(err) {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *SaleRepository) GetItems(ctx context.Context, saleID uuid.UUID) ([]*models.SaleItem, error) {
	query := `
		SELECT si.id, si.sale_id, si.product_id, si.quantity, si.unit_price, si.subtotal, si.notes
		FROM sale_items si
		WHERE si.sale_id = $1
	`
	rows, err := r.pool.Query(ctx, query, saleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.SaleItem
	for rows.Next() {
		var item models.SaleItem
		if err := rows.Scan(&item.ID, &item.SaleID, &item.ProductID, &item.Quantity, &item.UnitPrice, &item.Subtotal, &item.Notes); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, rows.Err()
}

func (r *SaleRepository) GetItemToppings(ctx context.Context, saleItemID uuid.UUID) ([]*models.Topping, error) {
	query := `SELECT id, sale_item_id, name, price, quantity FROM sale_item_toppings WHERE sale_item_id = $1`
	rows, err := r.pool.Query(ctx, query, saleItemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var toppings []*models.Topping
	for rows.Next() {
		var t models.Topping
		if err := rows.Scan(&t.ID, &t.SaleItemID, &t.Name, &t.Price, &t.Quantity); err != nil {
			return nil, err
		}
		toppings = append(toppings, &t)
	}
	return toppings, rows.Err()
}

func (r *SaleRepository) GetPayments(ctx context.Context, saleID uuid.UUID) ([]*models.SalePayment, error) {
	query := `SELECT id, sale_id, method, amount, reference FROM sale_payments WHERE sale_id = $1`
	rows, err := r.pool.Query(ctx, query, saleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*models.SalePayment
	for rows.Next() {
		var p models.SalePayment
		if err := rows.Scan(&p.ID, &p.SaleID, &p.Method, &p.Amount, &p.Reference); err != nil {
			return nil, err
		}
		payments = append(payments, &p)
	}
	return payments, rows.Err()
}
