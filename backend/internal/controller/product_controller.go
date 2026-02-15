package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pos-saas/restaurant-pos/internal/service"
)

type ProductController struct {
	productService *service.ProductService
}

func NewProductController(productService *service.ProductService) *ProductController {
	return &ProductController{productService: productService}
}

func (c *ProductController) getRestaurantID(ctx *gin.Context) (uuid.UUID, bool) {
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

func (c *ProductController) Create(ctx *gin.Context) {
	restaurantID, ok := c.getRestaurantID(ctx)
	if !ok {
		return
	}

	var input service.CreateProductInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "datos inválidos: " + err.Error()})
		return
	}

	product, err := c.productService.Create(ctx.Request.Context(), restaurantID, input)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, product)
}

func (c *ProductController) GetByID(ctx *gin.Context) {
	restaurantID, ok := c.getRestaurantID(ctx)
	if !ok {
		return
	}

	productID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	product, err := c.productService.GetByID(ctx.Request.Context(), restaurantID, productID)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, product)
}

func (c *ProductController) List(ctx *gin.Context) {
	restaurantID, ok := c.getRestaurantID(ctx)
	if !ok {
		return
	}

	var categoryID *uuid.UUID
	if catID := ctx.Query("category_id"); catID != "" {
		id, err := uuid.Parse(catID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "category_id inválido"})
			return
		}
		categoryID = &id
	}

	activeOnly := ctx.Query("active") != "false"

	products, err := c.productService.List(ctx.Request.Context(), restaurantID, categoryID, activeOnly)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, products)
}

func (c *ProductController) Update(ctx *gin.Context) {
	restaurantID, ok := c.getRestaurantID(ctx)
	if !ok {
		return
	}

	productID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input service.UpdateProductInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "datos inválidos: " + err.Error()})
		return
	}

	product, err := c.productService.Update(ctx.Request.Context(), restaurantID, productID, input)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, product)
}

func (c *ProductController) Delete(ctx *gin.Context) {
	restaurantID, ok := c.getRestaurantID(ctx)
	if !ok {
		return
	}

	productID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := c.productService.Delete(ctx.Request.Context(), restaurantID, productID); err != nil {
		handleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
