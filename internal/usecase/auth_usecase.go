package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sajudin/pos-app-server/internal/domain"
	"github.com/sajudin/pos-app-server/pkg/utils"
)

type authUsecase struct {
	userRepo     domain.UserRepository
	businessRepo domain.BusinessRepository
	secret       string
}

func NewAuthUsecase(ur domain.UserRepository, br domain.BusinessRepository, secret string) domain.AuthUsecase {
	return &authUsecase{userRepo: ur, businessRepo: br, secret: secret}
}

func (u *authUsecase) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("user tidak ditemukan")
	}

	if !utils.CheckPassword(password, user.Password) {
		return "", errors.New("password salah")
	}

	// Update token generation to include Role
	token, err := utils.GenerateToken(user.ID.String(), user.BusinessID.String(), user.Role, u.secret)
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
		Role:     "OWNER", // Default role for registration is OWNER
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
		BusinessPhone:   user.Business.Phone,
		BusinessLogoURL: user.Business.LogoURL,
		Role:            user.Role,
	}, nil
}

func (u *authUsecase) CreateStaff(ctx context.Context, email, password, role string, businessID uuid.UUID) error {
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

func (u *authUsecase) GetStaff(ctx context.Context, businessID uuid.UUID) ([]domain.UserResponse, error) {
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

func (u *authUsecase) UpdateBusiness(ctx context.Context, businessID uuid.UUID, req domain.UpdateBusinessRequest) error {
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
