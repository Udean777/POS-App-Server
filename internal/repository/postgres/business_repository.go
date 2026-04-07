package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/sajudin/pos-app-server/internal/domain"
	"gorm.io/gorm"
)

type gormBusinessRepository struct {
	db *gorm.DB
}

func NewGormBusinessRepository(db *gorm.DB) domain.BusinessRepository {
	return &gormBusinessRepository{db}
}

func (r *gormBusinessRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Business, error) {
	var business domain.Business
	err := r.db.WithContext(ctx).First(&business, "id = ?", id).Error
	return &business, err
}

func (r *gormBusinessRepository) Update(ctx context.Context, b *domain.Business) error {
	return r.db.WithContext(ctx).Save(b).Error
}
