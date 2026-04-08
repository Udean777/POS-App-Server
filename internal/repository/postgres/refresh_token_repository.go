package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/sajudin/pos-app-server/internal/domain"
	"gorm.io/gorm"
)

type gormRefreshTokenRepository struct {
	db *gorm.DB
}

func NewGormRefreshTokenRepository(db *gorm.DB) domain.RefreshTokenRepository {
	return &gormRefreshTokenRepository{db}
}

func (r *gormRefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *gormRefreshTokenRepository) GetByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	var rt domain.RefreshToken
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *gormRefreshTokenRepository) RevokeByToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Model(&domain.RefreshToken{}).
		Where("token = ?", token).
		Update("is_revoked", true).Error
}

func (r *gormRefreshTokenRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&domain.RefreshToken{}).Error
}
