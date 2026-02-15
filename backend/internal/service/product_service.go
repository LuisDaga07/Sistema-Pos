package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pos-saas/restaurant-pos/internal/errors"
	"github.com/pos-saas/restaurant-pos/internal/models"
	"github.com/pos-saas/restaurant-pos/internal/repository"
)

type ProductService struct {
	productRepo  *repository.ProductRepository
	categoryRepo *repository.CategoryRepository
}

func NewProductService(productRepo *repository.ProductRepository, categoryRepo *repository.CategoryRepository) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

type CreateProductInput struct {
	CategoryID  *string  `json:"category_id"`
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Price       float64  `json:"price" binding:"required,gt=0"`
	ImageURL    string   `json:"image_url"`
	Active      bool     `json:"active"`
}

type UpdateProductInput struct {
	CategoryID  *string  `json:"category_id"`
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	ImageURL    *string  `json:"image_url"`
	Active      *bool    `json:"active"`
}

func (s *ProductService) Create(ctx context.Context, restaurantID uuid.UUID, input CreateProductInput) (*models.Product, error) {
	var categoryID *uuid.UUID
	if input.CategoryID != nil && *input.CategoryID != "" {
		id, err := uuid.Parse(*input.CategoryID)
		if err != nil {
			return nil, NewValidationError("category_id", "UUID inválido")
		}
		if _, err := s.categoryRepo.GetByID(ctx, restaurantID, id); err != nil {
			if errors.Is(err, errors.ErrNotFound) {
				return nil, NewValidationError("category_id", "categoría no encontrada")
			}
			return nil, err
		}
		categoryID = &id
	}

	product := &models.Product{
		ID:           uuid.New(),
		RestaurantID: restaurantID,
		CategoryID:   categoryID,
		Name:         input.Name,
		Description:  input.Description,
		Price:        input.Price,
		ImageURL:     input.ImageURL,
		Active:       input.Active,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) GetByID(ctx context.Context, restaurantID, productID uuid.UUID) (*models.Product, error) {
	return s.productRepo.GetByID(ctx, restaurantID, productID)
}

func (s *ProductService) List(ctx context.Context, restaurantID uuid.UUID, categoryID *uuid.UUID, activeOnly bool) ([]*models.Product, error) {
	return s.productRepo.List(ctx, restaurantID, categoryID, activeOnly)
}

func (s *ProductService) Update(ctx context.Context, restaurantID, productID uuid.UUID, input UpdateProductInput) (*models.Product, error) {
	product, err := s.productRepo.GetByID(ctx, restaurantID, productID)
	if err != nil {
		return nil, err
	}

	if input.CategoryID != nil {
		if *input.CategoryID == "" {
			product.CategoryID = nil
		} else {
			id, err := uuid.Parse(*input.CategoryID)
			if err != nil {
				return nil, NewValidationError("category_id", "UUID inválido")
			}
			if _, err := s.categoryRepo.GetByID(ctx, restaurantID, id); err != nil && errors.Is(err, errors.ErrNotFound) {
				return nil, NewValidationError("category_id", "categoría no encontrada")
			}
			product.CategoryID = &id
		}
	}
	if input.Name != nil {
		product.Name = *input.Name
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.Price != nil && *input.Price >= 0 {
		product.Price = *input.Price
	}
	if input.ImageURL != nil {
		product.ImageURL = *input.ImageURL
	}
	if input.Active != nil {
		product.Active = *input.Active
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) Delete(ctx context.Context, restaurantID, productID uuid.UUID) error {
	return s.productRepo.Delete(ctx, restaurantID, productID)
}
