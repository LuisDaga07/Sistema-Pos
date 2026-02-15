package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pos-saas/restaurant-pos/internal/service"
)

type SaleController struct {
	saleService *service.SaleService
	pdfService  *service.PDFService
}

func NewSaleController(saleService *service.SaleService, pdfService *service.PDFService) *SaleController {
	return &SaleController{saleService: saleService, pdfService: pdfService}
}

func (c *SaleController) getIDs(ctx *gin.Context) (restaurantID, userID uuid.UUID, ok bool) {
	rid, ok1 := ctx.Get("restaurant_id")
	uid, ok2 := ctx.Get("user_id")
	if !ok1 || !ok2 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return uuid.Nil, uuid.Nil, false
	}
	ridStr, ok1 := rid.(string)
	uidStr, ok2 := uid.(string)
	if !ok1 || !ok2 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"})
		return uuid.Nil, uuid.Nil, false
	}
	parsedRid, err := uuid.Parse(ridStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "restaurant_id inválido"})
		return uuid.Nil, uuid.Nil, false
	}
	parsedUid, err := uuid.Parse(uidStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id inválido"})
		return uuid.Nil, uuid.Nil, false
	}
	return parsedRid, parsedUid, true
}

func (c *SaleController) Create(ctx *gin.Context) {
	restaurantID, userID, ok := c.getIDs(ctx)
	if !ok {
		return
	}

	var input service.CreateSaleInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "datos inválidos: " + err.Error()})
		return
	}

	sale, err := c.saleService.Create(ctx.Request.Context(), restaurantID, userID, input)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, sale)
}

func (c *SaleController) GetByID(ctx *gin.Context) {
	restaurantID, _, ok := c.getIDs(ctx)
	if !ok {
		return
	}

	saleID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	sale, items, payments, restaurant, err := c.saleService.GetByID(ctx.Request.Context(), restaurantID, saleID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"sale":       sale,
		"items":      items,
		"payments":   payments,
		"restaurant": restaurant,
	})
}

func (c *SaleController) GeneratePDF(ctx *gin.Context) {
	restaurantID, _, ok := c.getIDs(ctx)
	if !ok {
		return
	}

	saleID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	pdfBytes, err := c.pdfService.GenerateInvoice(ctx.Request.Context(), restaurantID, saleID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename=factura-"+saleID.String()+".pdf")
	ctx.Header("Content-Type", "application/pdf")
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}
