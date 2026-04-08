package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/sajudin/pos-app-server/internal/domain"
	"github.com/sajudin/pos-app-server/pkg/utils"
)

type staffUsecase struct {
	userRepo domain.UserRepository
}

func NewStaffUsecase(ur domain.UserRepository) domain.StaffUsecase {
	return &staffUsecase{userRepo: ur}
}

func (u *staffUsecase) CreateStaff(ctx context.Context, email, password, role string, businessID uuid.UUID) error {
	// Pengecekan proaktif: Apakah email sudah digunakan?
	existingUser, _ := u.userRepo.GetByEmail(ctx, email)
	if existingUser != nil && existingUser.Email != "" {
		// Jika terdaftar di bisnis yang sama
		if existingUser.BusinessID == businessID {
			return domain.ErrEmailAlreadyExists
		}
		// Jika terdaftar di bisnis lain
		return domain.ErrEmailRegisteredByOtherBusiness
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:      email,
		Password:   hashedPassword,
		BusinessID: businessID,
		Role:       role,
	}

	return u.userRepo.AddUser(ctx, user)
}

func (u *staffUsecase) GetStaff(ctx context.Context, businessID uuid.UUID) ([]domain.UserResponse, error) {
	users, err := u.userRepo.GetByBusinessID(ctx, businessID)
	if err != nil {
		return nil, err
	}

	var responses []domain.UserResponse
	for _, user := range users {
		responses = append(responses, domain.UserResponse{
			ID:              user.ID,
			Email:           user.Email,
			BusinessID:      user.BusinessID,
			Role:            user.Role,
			BusinessName:    user.Business.Name,
			BusinessType:    user.Business.Type,
			BusinessAddress: user.Business.Address,
			BusinessPhone:   user.Business.Phone,
			BusinessLogoURL: user.Business.LogoURL,
		})
	}

	return responses, nil
}
