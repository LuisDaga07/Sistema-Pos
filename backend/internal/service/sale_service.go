package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pos-saas/restaurant-pos/internal/models"
	"github.com/pos-saas/restaurant-pos/internal/repository"
)

type SaleService struct {
	saleRepo    *repository.SaleRepository
	productRepo *repository.ProductRepository
	authRepo    *repository.AuthRepository
}

func NewSaleService(saleRepo *repository.SaleRepository, productRepo *repository.ProductRepository, authRepo *repository.AuthRepository) *SaleService {
	return &SaleService{
		saleRepo:    saleRepo,
		productRepo: productRepo,
		authRepo:    authRepo,
	}
}

type SaleItemInput struct {
	ProductID string   `json:"product_id" binding:"required"`
	Quantity  int      `json:"quantity" binding:"required,gt=0"`
	Notes     string   `json:"notes"`
	Toppings  []ToppingInput `json:"toppings"`
}

type ToppingInput struct {
	Name     string  `json:"name" binding:"required"`
	Price    float64 `json:"price" binding:"gte=0"`
	Quantity int     `json:"quantity" binding:"gte=0"`
}

type SalePaymentInput struct {
	Method    string  `json:"method" binding:"required,oneof=cash card transfer"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	Reference string  `json:"reference"`
}

type CreateSaleInput struct {
	Items    []SaleItemInput    `json:"items" binding:"required,min=1,dive"`
	Payments []SalePaymentInput `json:"payments" binding:"required,min=1,dive"`
}

func (s *SaleService) Create(ctx context.Context, restaurantID, userID uuid.UUID, input CreateSaleInput) (*models.Sale, error) {
	var total float64
	saleID := uuid.New()

	// Validar productos y calcular total
	for _, it := range input.Items {
		productID, err := uuid.Parse(it.ProductID)
		if err != nil {
			return nil, NewValidationError("product_id", "UUID inválido")
		}

		product, err := s.productRepo.GetByID(ctx, restaurantID, productID)
		if err != nil {
			return nil, NewValidationError("product_id", "producto no encontrado")
		}
		if !product.Active {
			return nil, NewValidationError("product_id", "producto inactivo")
		}

		itemTotal := product.Price * float64(it.Quantity)
		for _, tp := range it.Toppings {
			itemTotal += tp.Price * float64(tp.Quantity)
		}
		total += itemTotal
	}

	// Validar que la suma de pagos coincida con el total
	var paymentsTotal float64
	for _, p := range input.Payments {
		paymentsTotal += p.Amount
	}
	if paymentsTotal < total-0.01 || paymentsTotal > total+0.01 { // tolerancia por decimales
		return nil, NewValidationError("payments", "la suma de pagos debe coincidir con el total")
	}

	sale := &models.Sale{
		ID:           saleID,
		RestaurantID: restaurantID,
		UserID:       userID,
		Total:        total,
		Status:       "completed",
	}
	if err := s.saleRepo.Create(ctx, sale); err != nil {
		return nil, err
	}

	for _, it := range input.Items {
		productID, _ := uuid.Parse(it.ProductID)
		product, _ := s.productRepo.GetByID(ctx, restaurantID, productID)

		itemTotal := product.Price * float64(it.Quantity)
		for _, tp := range it.Toppings {
			itemTotal += tp.Price * float64(tp.Quantity)
		}

		item := &models.SaleItem{
			ID:        uuid.New(),
			SaleID:    saleID,
			ProductID: productID,
			Quantity:  it.Quantity,
			UnitPrice: product.Price,
			Subtotal:  itemTotal,
			Notes:     it.Notes,
		}
		if err := s.saleRepo.CreateItem(ctx, item); err != nil {
			return nil, err
		}

		for _, tp := range it.Toppings {
			if tp.Quantity <= 0 {
				continue
			}
			topping := &models.Topping{
				ID:         uuid.New(),
				SaleItemID: item.ID,
				Name:       tp.Name,
				Price:      tp.Price,
				Quantity:   tp.Quantity,
			}
			if err := s.saleRepo.CreateItemTopping(ctx, topping); err != nil {
				return nil, err
			}
		}
	}

	for _, p := range input.Payments {
		payment := &models.SalePayment{
			ID:        uuid.New(),
			SaleID:    saleID,
			Method:    p.Method,
			Amount:    p.Amount,
			Reference: p.Reference,
		}
		if err := s.saleRepo.CreatePayment(ctx, payment); err != nil {
			return nil, err
		}
	}

	return sale, nil
}

func (s *SaleService) GetByID(ctx context.Context, restaurantID, saleID uuid.UUID) (*models.Sale, []*models.SaleItem, []*models.SalePayment, *models.Restaurant, error) {
	sale, err := s.saleRepo.GetByID(ctx, restaurantID, saleID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	items, err := s.saleRepo.GetItems(ctx, saleID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	for _, item := range items {
		toppings, _ := s.saleRepo.GetItemToppings(ctx, item.ID)
		item.Toppings = toppings
	}

	payments, err := s.saleRepo.GetPayments(ctx, saleID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	restaurant, err := s.authRepo.GetRestaurantByID(ctx, restaurantID)
	if err != nil {
		return sale, items, payments, nil, nil
	}

	return sale, items, payments, restaurant, nil
}
