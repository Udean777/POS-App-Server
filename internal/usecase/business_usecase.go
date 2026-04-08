package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sajudin/pos-app-server/internal/domain"
)

type businessUsecase struct {
	businessRepo domain.BusinessRepository
}

func NewBusinessUsecase(br domain.BusinessRepository) domain.BusinessUsecase {
	return &businessUsecase{businessRepo: br}
}

func (u *businessUsecase) UpdateBusiness(ctx context.Context, businessID uuid.UUID, req domain.UpdateBusinessRequest) error {
	biz, err := u.businessRepo.GetByID(ctx, businessID)
	if err != nil {
		return errors.New("bisnis tidak ditemukan")
	}

	if req.Name != "" {
		biz.Name = req.Name
	}
	if req.Type != "" {
		biz.Type = req.Type
	}
	if req.Address != "" {
		biz.Address = req.Address
	}
	if req.Phone != "" {
		biz.Phone = req.Phone
	}
	if req.LogoURL != "" {
		biz.LogoURL = req.LogoURL
	}

	return u.businessRepo.Update(ctx, biz)
}
