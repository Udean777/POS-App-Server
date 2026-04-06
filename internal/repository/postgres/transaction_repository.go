package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/sajudin/pos-app-server/internal/domain"
	"gorm.io/gorm"
)

type gormTransactionRepository struct {
	db *gorm.DB
}

func NewGormTransactionRepository(db *gorm.DB) domain.TransactionRepository {
	return &gormTransactionRepository{db}
}

// Create menjalankan operasi simpan transaksi dan pengurangan stok dalam satu DB Transaction
func (r *gormTransactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	return r.db.WithContext(ctx).Transaction(func(dbTx *gorm.DB) error {
		// 1. Simpan Header Transaksi
		if err := dbTx.Create(tx).Error; err != nil {
			return err
		}

		// 2. Loop setiap item untuk update stok
		for _, item := range tx.Items {
			// Kurangi stok variant. Gunakan gorm.Expr untuk mengindari race condition.
			// Stok bisa minus jika memang diizinkan dari level usecase/frontend.
			res := dbTx.Model(&domain.Variant{}).
				Where("id = ?", item.VariantID).
				UpdateColumn("stock", gorm.Expr("stock - ?", item.Quantity))

			if res.Error != nil {
				return res.Error
			}
		}

		return nil
	})
}

func (r *gormTransactionRepository) GetByBusinessID(ctx context.Context, businessID uuid.UUID) ([]domain.Transaction, error) {
	var txs []domain.Transaction
	// Preload Items and Staff relations
	err := r.db.WithContext(ctx).
		Preload("Staff").
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Variant").
		Where("business_id = ?", businessID).
		Order("created_at desc").
		Find(&txs).Error
	return txs, err
}
