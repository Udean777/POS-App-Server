package postgres

import (
	"context"

	"github.com/sajudin/pos-app-server/internal/domain"
	"gorm.io/gorm"
)

type gormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) domain.UserRepository {
	return &gormUserRepository{db}
}

func (r *gormUserRepository) Create(ctx context.Context, u *domain.User, businessName string) error {
	// GORM Transaction: Atomik dan sangat bersih
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Buat Bisnis Baru
		business := domain.Business{
			Name: businessName,
			Type: "RETAIL",
		}
		if err := tx.Create(&business).Error; err != nil {
			return err
		}

		// 2. Hubungkan User ke Bisnis yang baru dibuat
		u.BusinessID = business.ID
		if err := tx.Create(u).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *gormUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}
