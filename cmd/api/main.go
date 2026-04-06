package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/sajudin/pos-app-server/internal/delivery/http"
	"github.com/sajudin/pos-app-server/internal/delivery/http/middleware"
	"github.com/sajudin/pos-app-server/internal/domain"
	repo "github.com/sajudin/pos-app-server/internal/repository/postgres"
	"github.com/sajudin/pos-app-server/internal/usecase"
	"gorm.io/driver/postgres"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal konek ke database:", err)
	}

	db.AutoMigrate(&domain.Business{}, &domain.User{}, &domain.Product{}, &domain.Variant{})

	r := gin.Default()
	r.SetTrustedProxies(nil)
	secret := os.Getenv("JWT_SECRET")

	// Repositories
	userRepo := repo.NewGormUserRepository(db)
	productRepo := repo.NewGormProductRepository(db)

	// Usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, secret)
	productUsecase := usecase.NewProductUsecase(productRepo)

	// Handlers
	authHandler := http.AuthHandler{AuthUsecase: authUsecase}
	productHandler := http.NewProductHandler(productUsecase)

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
		protected.POST("/products", productHandler.Create)
		protected.GET("/products", productHandler.GetAll)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server POS berjalan di port %s", port)
	r.Run(":" + port)
}
