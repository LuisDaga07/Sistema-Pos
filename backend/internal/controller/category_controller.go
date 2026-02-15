package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pos-saas/restaurant-pos/internal/models"
	"github.com/pos-saas/restaurant-pos/internal/repository"
)

type CategoryController struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryController(categoryRepo *repository.CategoryRepository) *CategoryController {
	return &CategoryController{categoryRepo: categoryRepo}
}

type CreateCategoryInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

func (c *CategoryController) getRestaurantID(ctx *gin.Context) (uuid.UUID, bool) {
	rid, ok := ctx.Get("restaurant_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return uuid.Nil, false
	}
	ridStr, ok := rid.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"})
		return uuid.Nil, false
	}
	parsed, err := uuid.Parse(ridStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "restaurant_id inválido"})
		return uuid.Nil, false
	}
	return parsed, true
}

func (c *CategoryController) Create(ctx *gin.Context) {
	restaurantID, ok := c.getRestaurantID(ctx)
	if !ok {
		return
	}

	var input CreateCategoryInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "datos inválidos: " + err.Error()})
		return
	}

	cat := &models.Category{
		ID:           uuid.New(),
		RestaurantID: restaurantID,
		Name:         input.Name,
		Description:  input.Description,
		SortOrder:    input.SortOrder,
	}
	if err := c.categoryRepo.Create(ctx.Request.Context(), cat); err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, cat)
}

func (c *CategoryController) List(ctx *gin.Context) {
	restaurantID, ok := c.getRestaurantID(ctx)
	if !ok {
		return
	}

	categories, err := c.categoryRepo.List(ctx.Request.Context(), restaurantID)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, categories)
}
