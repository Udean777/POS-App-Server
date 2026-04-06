package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BusinessID  uuid.UUID `gorm:"type:uuid;not null" json:"business_id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Variants    []Variant `gorm:"foreignKey:ProductID" json:"variants"`
	CreatedAt   time.Time `json:"created_at"`
}

type Variant struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProductID  uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	BusinessID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_business_sku" json:"business_id"`
	Name       string    `gorm:"not null" json:"name"`
	Price      float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock      int       `gorm:"not null" json:"stock"`
	SKU        string    `gorm:"uniqueIndex:idx_business_sku" json:"sku"`
}

type ProductRepository interface {
	Create(ctx context.Context, p *Product) error
	Fetch(ctx context.Context, businessID uuid.UUID) ([]Product, error)
}

type ProductUsecase interface {
	AddProduct(ctx context.Context, p *Product) error
	GetAllProducts(ctx context.Context, businessID uuid.UUID) ([]Product, error)
}
