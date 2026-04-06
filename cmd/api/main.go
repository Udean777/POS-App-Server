package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/sajudin/pos-app-server/internal/delivery/http"
	"github.com/sajudin/pos-app-server/internal/delivery/http/middleware"
	"github.com/sajudin/pos-app-server/internal/domain"
	repo "github.com/sajudin/pos-app-server/internal/repository/postgres"
	"github.com/sajudin/pos-app-server/internal/service"
	"github.com/sajudin/pos-app-server/internal/usecase"
	"gorm.io/driver/postgres"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system env")
	}

	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal konek ke database:", err)
	}

	db.AutoMigrate(
		&domain.Business{},
		&domain.User{},
		&domain.Product{},
		&domain.Variant{},
		&domain.Transaction{},
		&domain.TransactionItem{},
	)

	r := gin.Default()
	r.SetTrustedProxies(nil)
	secret := os.Getenv("JWT_SECRET")

	// Static files for local storage fallback
	r.Static("/uploads", "./uploads")

	// Initialize Storage Service (R2 or Local)
	storageService, err := service.NewS3StorageService(context.Background())
	if err != nil {
		log.Println("R2 storage not configured, using local storage:", err)
		publicURL := os.Getenv("APP_URL")
		if publicURL == "" {
			publicURL = "http://localhost:8080"
		}
		storageService = service.NewLocalStorageService(publicURL)
	}

	// Repositories
	userRepo := repo.NewGormUserRepository(db)
	productRepo := repo.NewGormProductRepository(db)
	txRepo := repo.NewGormTransactionRepository(db)

	// Usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, secret)
	productUsecase := usecase.NewProductUsecase(productRepo)
	txUsecase := usecase.NewTransactionUsecase(txRepo, productRepo)

	// Handlers
	authHandler := http.AuthHandler{AuthUsecase: authUsecase}
	productHandler := http.NewProductHandler(productUsecase, storageService)
	txHandler := http.NewTransactionHandler(txUsecase)

	// Public Routes
	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/register", authHandler.Register)
		v1.POST("/auth/login", authHandler.Login)
	}

	// Protected Routes
	protected := v1.Group("/")
	protected.Use(middleware.AuthMiddleware(secret))
	{
		protected.GET("/me", authHandler.GetProfile)

		// Staff Management (Owner Only)
		ownerOnly := protected.Group("/auth", middleware.RoleMiddleware("OWNER"))
		{
			ownerOnly.POST("/staff", authHandler.CreateStaff)
			ownerOnly.GET("/staff", authHandler.GetStaff)
		}

		// Product Routes
		protected.POST("/products/upload", productHandler.Upload)
		protected.POST("/products", productHandler.Create)
		protected.GET("/products", productHandler.GetAll)
		protected.GET("/products/:id", productHandler.GetByID)
		protected.PUT("/products/:id", productHandler.Update)
		protected.DELETE("/products/:id", productHandler.Delete)

		// Transaction Routes
		protected.POST("/transactions", txHandler.Checkout)
		protected.GET("/transactions", txHandler.GetAll)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server POS berjalan di port %s", port)
	r.Run(":" + port)
}
