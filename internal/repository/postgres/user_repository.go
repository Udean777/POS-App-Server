package postgres

import (
	"context"

	"github.com/google/uuid"
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

func (r *gormUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	// Preload "Business" agar data bisnis ikut terbawa sesuai relasi di domain
	err := r.db.WithContext(ctx).Preload("Business").First(&user, "id = ?", id).Error
	return &user, err
}

func (r *gormUserRepository) GetByBusinessID(ctx context.Context, businessID uuid.UUID) ([]domain.User, error) {
	var users []domain.User
	// Preload "Business" untuk konsistensi di struk/dashboard
	err := r.db.WithContext(ctx).
		Preload("Business").
		Where("business_id = ?", businessID).
		Find(&users).Error
	return users, err
}

func (r *gormUserRepository) AddUser(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}
