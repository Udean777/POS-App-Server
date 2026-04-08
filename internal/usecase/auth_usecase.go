package usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/sajudin/pos-app-server/internal/domain"
	"github.com/sajudin/pos-app-server/pkg/mail"
	"github.com/sajudin/pos-app-server/pkg/utils"
)

type authUsecase struct {
	userRepo         domain.UserRepository
	refreshTokenRepo domain.RefreshTokenRepository
	vcRepo           domain.VerificationCodeRepository
	mailer           mail.Mailer
	secret           string
}

func NewAuthUsecase(
	ur domain.UserRepository,
	rtr domain.RefreshTokenRepository,
	vcr domain.VerificationCodeRepository,
	m mail.Mailer,
	secret string,
) domain.AuthUsecase {
	return &authUsecase{
		userRepo:         ur,
		refreshTokenRepo: rtr,
		vcRepo:           vcr,
		mailer:           m,
		secret:           secret,
	}
}

func (u *authUsecase) Login(ctx context.Context, email string, password string) (*domain.TokenResponse, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if !utils.CheckPassword(password, user.Password) {
		return nil, errors.New("password salah")
	}

	// Cek Verifikasi
	if !user.IsVerified {
		return nil, errors.New("EMAIL_NOT_VERIFIED")
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
	// Pengecekan proaktif: Apakah email sudah digunakan?
	existingUser, _ := u.userRepo.GetByEmail(ctx, email)
	if existingUser != nil && existingUser.Email != "" {
		return domain.ErrEmailAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:      email,
		Password:   hashedPassword,
		Role:       "OWNER",
		IsVerified: false,
	}

	if err := u.userRepo.Create(ctx, user, bizName); err != nil {
		return err
	}

	// Generate & Kirim OTP
	return u.generateAndSendOTP(ctx, email, "REGISTER")
}

func (u *authUsecase) VerifyOTP(ctx context.Context, email, code string) (*domain.TokenResponse, error) {
	vc, err := u.vcRepo.GetLastByEmail(ctx, email, "REGISTER")
	if err != nil {
		return nil, errors.New("kode verifikasi tidak ditemukan atau sudah kadaluarsa")
	}

	if vc.Code != code {
		return nil, errors.New("kode verifikasi tidak valid")
	}

	if vc.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("kode verifikasi sudah kadaluarsa")
	}

	// Update User
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user.IsVerified = true
	user.VerifiedAt = &now

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Hapus OTP setelah berhasil
	if err := u.vcRepo.DeleteByEmail(ctx, email, "REGISTER"); err != nil {
		// Log error but don't fail verification
		fmt.Printf("failed to delete verification code: %v\n", err)
	}

	// --- GENERATE TOKENS FOR DIRECT LOGIN ---
	accessToken, err := utils.GenerateAccessToken(user.ID.String(), user.BusinessID.String(), user.Role, u.secret)
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := utils.GenerateRefreshToken(user.ID.String(), user.BusinessID.String(), user.Role, u.secret)
	if err != nil {
		return nil, err
	}

	// Simpan refresh token ke DB
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

func (u *authUsecase) ResendOTP(ctx context.Context, email string) error {
	// Cek apakah user ada & belum verified
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	if user.IsVerified {
		return errors.New("akun sudah terverifikasi")
	}

	// Generate & Kirim OTP baru
	return u.generateAndSendOTP(ctx, email, "REGISTER")
}

func (u *authUsecase) generateAndSendOTP(ctx context.Context, email, codeType string) error {
	// Generate 6 digit OTP
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	code := fmt.Sprintf("%06d", n.Int64()+100000)

	vc := &domain.VerificationCode{
		Email:     email,
		Code:      code,
		Type:      codeType,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := u.vcRepo.Create(ctx, vc); err != nil {
		return err
	}

	// Kirim Email
	return u.mailer.SendOTP(email, code)
}

func (u *authUsecase) generateAndSendForgotPasswordOTP(ctx context.Context, email string) error {
	// Generate 6 digit OTP
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	code := fmt.Sprintf("%06d", n.Int64()+100000)

	vc := &domain.VerificationCode{
		Email:     email,
		Code:      code,
		Type:      "FORGOT_PASSWORD",
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := u.vcRepo.Create(ctx, vc); err != nil {
		return err
	}

	// Kirim Email Khusus Forgot Password
	return u.mailer.SendForgotPasswordOTP(email, code)
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
		IsVerified:      user.IsVerified,
	}, nil
}

func (u *authUsecase) ForgotPassword(ctx context.Context, email string) error {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	if user.Role == "OWNER" {
		return u.generateAndSendForgotPasswordOTP(ctx, email)
	}

	// Untuk STAFF/ADMIN, kirim notifikasi ke OWNER (Ajuan)
	users, err := u.userRepo.GetByBusinessID(ctx, user.BusinessID)
	if err != nil {
		return err
	}

	var ownerEmail string
	for _, u := range users {
		if u.Role == "OWNER" {
			ownerEmail = u.Email
			break
		}
	}

	if ownerEmail != "" {
		return u.mailer.SendStaffResetRequest(ownerEmail, user.Email)
	}

	return errors.New("owner tidak ditemukan untuk bisnis ini")
}

func (u *authUsecase) ResetPassword(ctx context.Context, email, code, newPassword string) (*domain.TokenResponse, error) {
	vc, err := u.vcRepo.GetLastByEmail(ctx, email, "FORGOT_PASSWORD")
	if err != nil {
		return nil, errors.New("kode verifikasi tidak ditemukan atau sudah kadaluarsa")
	}

	if vc.Code != code {
		return nil, errors.New("kode verifikasi tidak valid")
	}

	if vc.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("kode verifikasi sudah kadaluarsa")
	}

	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return nil, err
	}

	if err := u.userRepo.UpdatePassword(ctx, user.ID, hashedPassword); err != nil {
		return nil, err
	}

	// Hapus OTP
	u.vcRepo.DeleteByEmail(ctx, email, "FORGOT_PASSWORD")

	// Revoke Refresh Token
	u.refreshTokenRepo.DeleteByUserID(ctx, user.ID)

	// Login otomatis setelah reset
	accessToken, err := utils.GenerateAccessToken(user.ID.String(), user.BusinessID.String(), user.Role, u.secret)
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := utils.GenerateRefreshToken(user.ID.String(), user.BusinessID.String(), user.Role, u.secret)
	if err != nil {
		return nil, err
	}

	u.refreshTokenRepo.Create(ctx, &domain.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})

	return &domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
	}, nil
}
