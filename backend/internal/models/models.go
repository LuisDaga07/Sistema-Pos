package models

import (
	"time"

	"github.com/google/uuid"
)

// Restaurant representa un restaurante (tenant)
type Restaurant struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone,omitempty"`
	Address   string     `json:"address,omitempty"`
	TaxID     string     `json:"tax_id,omitempty"`
	LogoURL   string     `json:"logo_url,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

// User representa un usuario del sistema
type User struct {
	ID           uuid.UUID  `json:"id"`
	RestaurantID uuid.UUID  `json:"restaurant_id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	Role         string     `json:"role"` // admin, cajero
	Active       bool       `json:"active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"-" db:"deleted_at"`
}

// Category representa una categoría de productos
type Category struct {
	ID           uuid.UUID `json:"id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	SortOrder    int       `json:"sort_order"`
}

// Product representa un producto del menú
type Product struct {
	ID           uuid.UUID  `json:"id"`
	RestaurantID uuid.UUID  `json:"restaurant_id"`
	CategoryID   *uuid.UUID `json:"category_id,omitempty"`
	Name         string     `json:"name"`
	Description  string     `json:"description,omitempty"`
	Price        float64    `json:"price"`
	ImageURL     string     `json:"image_url,omitempty"`
	Active       bool       `json:"active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Sale representa una venta
type Sale struct {
	ID           uuid.UUID  `json:"id"`
	RestaurantID uuid.UUID  `json:"restaurant_id"`
	UserID       uuid.UUID  `json:"user_id"`
	Total        float64    `json:"total"`
	Status       string     `json:"status"` // pending, completed, cancelled
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// SaleItem representa un item en una venta
type SaleItem struct {
	ID        uuid.UUID  `json:"id"`
	SaleID    uuid.UUID  `json:"sale_id"`
	ProductID uuid.UUID  `json:"product_id"`
	Quantity  int        `json:"quantity"`
	UnitPrice float64    `json:"unit_price"`
	Subtotal  float64    `json:"subtotal"`
	Notes     string     `json:"notes,omitempty"`
	Toppings  []*Topping `json:"toppings,omitempty" db:"-"`
}

// Topping representa un adicional/topping
type Topping struct {
	ID        uuid.UUID `json:"id"`
	SaleItemID uuid.UUID `json:"sale_item_id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
}

// SalePayment representa un método de pago en una venta
type SalePayment struct {
	ID        uuid.UUID `json:"id"`
	SaleID    uuid.UUID `json:"sale_id"`
	Method    string    `json:"method"` // cash, card, transfer
	Amount    float64   `json:"amount"`
	Reference string    `json:"reference,omitempty"`
}
