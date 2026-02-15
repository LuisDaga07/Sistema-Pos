package service

import (
	"bytes"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/pos-saas/restaurant-pos/internal/repository"
)

type PDFService struct {
	saleRepo    *repository.SaleRepository
	productRepo *repository.ProductRepository
	authRepo    *repository.AuthRepository
}

func NewPDFService(saleRepo *repository.SaleRepository, productRepo *repository.ProductRepository, authRepo *repository.AuthRepository) *PDFService {
	return &PDFService{
		saleRepo:    saleRepo,
		productRepo: productRepo,
		authRepo:    authRepo,
	}
}

func (s *PDFService) GenerateInvoice(ctx context.Context, restaurantID, saleID uuid.UUID) ([]byte, error) {
	sale, err := s.saleRepo.GetByID(ctx, restaurantID, saleID)
	if err != nil {
		return nil, err
	}

	items, err := s.saleRepo.GetItems(ctx, saleID)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		toppings, _ := s.saleRepo.GetItemToppings(ctx, item.ID)
		item.Toppings = toppings
	}

	payments, err := s.saleRepo.GetPayments(ctx, saleID)
	if err != nil {
		return nil, err
	}

	restaurant, err := s.authRepo.GetRestaurantByID(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	itemDetails := make([]struct {
		Name     string
		Qty      int
		Price    float64
		Subtotal float64
		Toppings []string
	}, len(items))

	for i, item := range items {
		product, _ := s.productRepo.GetByID(ctx, restaurantID, item.ProductID)
		name := "Producto"
		if product != nil {
			name = product.Name
		}
		toppingStrs := make([]string, 0)
		for _, t := range item.Toppings {
			toppingStrs = append(toppingStrs, fmt.Sprintf("  + %s x%d $%.2f", t.Name, t.Quantity, t.Price*float64(t.Quantity)))
		}
		itemDetails[i] = struct {
			Name     string
			Qty      int
			Price    float64
			Subtotal float64
			Toppings []string
		}{Name: name, Qty: item.Quantity, Price: item.UnitPrice, Subtotal: item.Subtotal, Toppings: toppingStrs}
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 12)

	// Encabezado
	pdf.SetX(20)
	pdf.SetY(20)
	pdf.SetFont("Helvetica", "B", 16)
	pdf.CellFormat(0, 8, restaurant.Name, "", 0, "L", false, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Helvetica", "", 10)
	if restaurant.Address != "" {
		pdf.CellFormat(0, 6, restaurant.Address, "", 0, "L", false, 0, "")
		pdf.Ln(5)
	}
	if restaurant.TaxID != "" {
		pdf.CellFormat(0, 6, "RFC/NIT: "+restaurant.TaxID, "", 0, "L", false, 0, "")
		pdf.Ln(5)
	}
	if restaurant.Phone != "" {
		pdf.CellFormat(0, 6, "Tel: "+restaurant.Phone, "", 0, "L", false, 0, "")
		pdf.Ln(10)
	}

	pdf.SetFont("Helvetica", "B", 14)
	pdf.CellFormat(0, 8, "FACTURA / TICKET DE VENTA", "", 0, "L", false, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(0, 6, fmt.Sprintf("Venta #%s", saleID.String()[:8]), "", 0, "L", false, 0, "")
	pdf.Ln(4)
	pdf.CellFormat(0, 6, fmt.Sprintf("Fecha: %s", sale.CreatedAt.Format("02/01/2006 15:04")), "", 0, "L", false, 0, "")
	pdf.Ln(12)

	// Tabla
	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(80, 7, "Producto", "B", 0, "L", false, 0, "")
	pdf.CellFormat(20, 7, "Cant", "B", 0, "R", false, 0, "")
	pdf.CellFormat(35, 7, "P.U.", "B", 0, "R", false, 0, "")
	pdf.CellFormat(50, 7, "Subtotal", "B", 0, "R", false, 0, "")
	pdf.Ln(8)

	pdf.SetFont("Helvetica", "", 10)
	for _, it := range itemDetails {
		pdf.CellFormat(80, 6, it.Name, "", 0, "L", false, 0, "")
		pdf.CellFormat(20, 6, fmt.Sprintf("%d", it.Qty), "", 0, "R", false, 0, "")
		pdf.CellFormat(35, 6, fmt.Sprintf("$%.2f", it.Price), "", 0, "R", false, 0, "")
		pdf.CellFormat(50, 6, fmt.Sprintf("$%.2f", it.Subtotal), "", 0, "R", false, 0, "")
		pdf.Ln(5)
		for _, tp := range it.Toppings {
			pdf.CellFormat(80, 5, tp, "", 0, "L", false, 0, "")
			pdf.Ln(5)
		}
	}

	pdf.Ln(8)
	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(135, 8, "TOTAL:", "", 0, "R", false, 0, "")
	pdf.CellFormat(50, 8, fmt.Sprintf("$%.2f", sale.Total), "", 0, "R", false, 0, "")
	pdf.Ln(12)

	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(0, 6, "Metodos de pago:", "", 0, "L", false, 0, "")
	pdf.Ln(6)
	pdf.SetFont("Helvetica", "", 10)
	for _, p := range payments {
		method := p.Method
		switch p.Method {
		case "cash":
			method = "Efectivo"
		case "card":
			method = "Tarjeta"
		case "transfer":
			method = "Transferencia"
		}
		line := fmt.Sprintf("  - %s: $%.2f", method, p.Amount)
		if p.Reference != "" {
			line += " (Ref: " + p.Reference + ")"
		}
		pdf.CellFormat(0, 5, line, "", 0, "L", false, 0, "")
		pdf.Ln(5)
	}

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
