package postgres

import (
	"context"

	"github.com/sajudin/pos-app-server/internal/domain"
	"gorm.io/gorm"
)

type gormVerificationCodeRepository struct {
	db *gorm.DB
}

func NewGormVerificationCodeRepository(db *gorm.DB) domain.VerificationCodeRepository {
	return &gormVerificationCodeRepository{db}
}

func (r *gormVerificationCodeRepository) Create(ctx context.Context, vc *domain.VerificationCode) error {
	return r.db.WithContext(ctx).Create(vc).Error
}

func (r *gormVerificationCodeRepository) GetLastByEmail(ctx context.Context, email, codeType string) (*domain.VerificationCode, error) {
	var vc domain.VerificationCode
	err := r.db.WithContext(ctx).
		Where("email = ? AND type = ?", email, codeType).
		Order("created_at DESC").
		First(&vc).Error
	return &vc, err
}

func (r *gormVerificationCodeRepository) DeleteByEmail(ctx context.Context, email, codeType string) error {
	return r.db.WithContext(ctx).
		Where("email = ? AND type = ?", email, codeType).
		Delete(&domain.VerificationCode{}).Error
}
