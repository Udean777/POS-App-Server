package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sajudin/pos-app-server/internal/domain"
	"github.com/sajudin/pos-app-server/pkg/utils"
)

type authUsecase struct {
	userRepo domain.UserRepository
	// businessRepo domain.BusinessRepository
	secret string
}

func NewAuthUsecase(ur domain.UserRepository, secret string) domain.AuthUsecase {
	return &authUsecase{userRepo: ur, secret: secret}
}

func (u *authUsecase) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("user tidak ditemukan")
	}

	if !utils.CheckPassword(password, user.Password) {
		return "", errors.New("password salah")
	}

	token, err := utils.GenerateToken(user.ID.String(), user.BusinessID.String(), u.secret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *authUsecase) Register(ctx context.Context, email, password, bizName string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:    email,
		Password: hashedPassword,
	}

	return u.userRepo.Create(ctx, user, bizName)
}

func (u *authUsecase) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.UserResponse, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("profil tidak ditemukan")
	}

	return &domain.UserResponse{
		ID:              user.ID,
		Email:           user.Email,
		BusinessID:      user.Business.ID,
		BusinessName:    user.Business.Name,
		BusinessType:    user.Business.Type,
		BusinessAddress: user.Business.Address,
		Role:            "OWNER",
	}, nil
}
