package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/sajudin/pos-app-server/internal/domain"
	"gorm.io/gorm"
)

type gormProductRepository struct {
	db *gorm.DB
}

func NewGormProductRepository(db *gorm.DB) domain.ProductRepository {
	return &gormProductRepository{db}
}

func (r *gormProductRepository) Create(ctx context.Context, p *domain.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *gormProductRepository) Fetch(ctx context.Context, businessID uuid.UUID) ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.WithContext(ctx).Where("business_id = ?", businessID).Preload("Variants").Find(&products).Error
	return products, err
}
