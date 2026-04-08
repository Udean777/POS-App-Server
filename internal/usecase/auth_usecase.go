package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sajudin/pos-app-server/internal/domain"
	"github.com/sajudin/pos-app-server/pkg/utils"
)

type authUsecase struct {
	userRepo         domain.UserRepository
	refreshTokenRepo domain.RefreshTokenRepository
	secret           string
}

func NewAuthUsecase(ur domain.UserRepository, rtr domain.RefreshTokenRepository, secret string) domain.AuthUsecase {
	return &authUsecase{userRepo: ur, refreshTokenRepo: rtr, secret: secret}
}

func (u *authUsecase) Login(ctx context.Context, email string, password string) (*domain.TokenResponse, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if !utils.CheckPassword(password, user.Password) {
		return nil, errors.New("password salah")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID.String(), user.BusinessID.String(), user.Role, u.secret)
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := utils.GenerateRefreshToken(user.ID.String(), user.BusinessID.String(), user.Role, u.secret)
	if err != nil {
		return nil, err
	}

	// Simpan refresh token ke DB untuk tracking/rotation
	refreshToken := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := u.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, err
	}

	return &domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
	}, nil
}

func (u *authUsecase) Refresh(ctx context.Context, refreshTokenString string) (*domain.TokenResponse, error) {
	// 1. Validasi Token secara JWT
	claims, err := utils.ValidateToken(refreshTokenString, u.secret)
	if err != nil {
		return nil, errors.New("refresh token tidak valid")
	}

	// 2. Cek di Database
	storedToken, err := u.refreshTokenRepo.GetByToken(ctx, refreshTokenString)
	if err != nil {
		return nil, errors.New("refresh token tidak ditemukan")
	}

	if storedToken.IsRevoked {
		return nil, errors.New("refresh token sudah dicabut")
	}

	if storedToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token kadaluarsa")
	}

	// 3. Generate Pasangan Token Baru (Token Rotation)
	accessToken, err := utils.GenerateAccessToken(claims.UserID, claims.BusinessID, claims.Role, u.secret)
	if err != nil {
		return nil, err
	}

	newRefreshTokenString, err := utils.GenerateRefreshToken(claims.UserID, claims.BusinessID, claims.Role, u.secret)
	if err != nil {
		return nil, err
	}

	// 4. Update Database: Revoke token lama, Simpan yang baru
	if err := u.refreshTokenRepo.RevokeByToken(ctx, refreshTokenString); err != nil {
		return nil, err
	}

	userID, _ := uuid.Parse(claims.UserID)
	newRefreshToken := &domain.RefreshToken{
		UserID:    userID,
		Token:     newRefreshTokenString,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := u.refreshTokenRepo.Create(ctx, newRefreshToken); err != nil {
		return nil, err
	}

	return &domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshTokenString,
	}, nil
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
