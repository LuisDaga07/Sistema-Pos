package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pos-saas/restaurant-pos/config"
	"github.com/pos-saas/restaurant-pos/internal/controller"
	"github.com/pos-saas/restaurant-pos/internal/database"
	"github.com/pos-saas/restaurant-pos/internal/middleware"
	"github.com/pos-saas/restaurant-pos/internal/repository"
	"github.com/pos-saas/restaurant-pos/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	pool, err := database.NewPool(cfg.Database.DSN())
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	gin.SetMode(cfg.Server.GinMode)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Repositories
	authRepo := repository.NewAuthRepository(pool)
	productRepo := repository.NewProductRepository(pool)
	categoryRepo := repository.NewCategoryRepository(pool)
	saleRepo := repository.NewSaleRepository(pool)

	// Services
	authService := service.NewAuthService(authRepo, cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	productService := service.NewProductService(productRepo, categoryRepo)
	saleService := service.NewSaleService(saleRepo, productRepo, authRepo)
	pdfService := service.NewPDFService(saleRepo, productRepo, authRepo)

	// Controllers
	authCtrl := controller.NewAuthController(authService)
	productCtrl := controller.NewProductController(productService)
	categoryCtrl := controller.NewCategoryController(categoryRepo)
	saleCtrl := controller.NewSaleController(saleService, pdfService)

	// Public routes
	api := r.Group("/api/v1")
	api.POST("/auth/register", authCtrl.Register)
	api.POST("/auth/login", authCtrl.Login)

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthRequired(cfg.JWT.Secret))
	{
		protected.GET("/categories", categoryCtrl.List)
		protected.POST("/categories", categoryCtrl.Create)

		protected.GET("/products", productCtrl.List)
		protected.GET("/products/:id", productCtrl.GetByID)
		protected.POST("/products", productCtrl.Create)
		protected.PUT("/products/:id", productCtrl.Update)
		protected.DELETE("/products/:id", productCtrl.Delete)

		protected.POST("/sales", saleCtrl.Create)
		protected.GET("/sales/:id", saleCtrl.GetByID)
		protected.GET("/sales/:id/pdf", saleCtrl.GeneratePDF)
	}

	addr := ":" + cfg.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server: %v", err)
	}
}
