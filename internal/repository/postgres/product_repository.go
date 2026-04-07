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

func (r *gormProductRepository) GetByID(ctx context.Context, id uuid.UUID, businessID uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).
		Preload("Variants").
		Where("id = ? AND business_id = ?", id, businessID).
		First(&product).Error
	return &product, err
}

func (r *gormProductRepository) Update(ctx context.Context, p *domain.Product) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var incomingVariantIDs []uuid.UUID
		for _, v := range p.Variants {
			if v.ID != uuid.Nil {
				incomingVariantIDs = append(incomingVariantIDs, v.ID)
			}
		}

		if len(incomingVariantIDs) > 0 {
			if err := tx.Where("product_id = ? AND id NOT IN ?", p.ID, incomingVariantIDs).Delete(&domain.Variant{}).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Where("product_id = ?", p.ID).Delete(&domain.Variant{}).Error; err != nil {
				return err
			}
		}

		return tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(p).Error
	})
}

func (r *gormProductRepository) Delete(ctx context.Context, id uuid.UUID, businessID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id = ?", id).Delete(&domain.Variant{}).Error; err != nil {
			return err
		}

		return tx.Where("id = ? AND business_id = ?", id, businessID).Delete(&domain.Product{}).Error
	})
}

func (r *gormProductRepository) RestockVariant(ctx context.Context, variantID uuid.UUID, businessID uuid.UUID, quantity int) error {
	return r.db.WithContext(ctx).
		Model(&domain.Variant{}).
		Where("id = ? AND business_id = ?", variantID, businessID).
		UpdateColumn("stock", gorm.Expr("stock + ?", quantity)).Error
}
