package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEmailAlreadyExists             = errors.New("email sudah terdaftar")
	ErrEmailRegisteredByOtherBusiness = errors.New("Email sudah terdaftar oleh bisnis lain")
)

type User struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email      string     `gorm:"uniqueIndex;not null" json:"email"`
	Password   string     `gorm:"not null" json:"-"`
	BusinessID uuid.UUID  `gorm:"type:uuid;not null" json:"business_id"`
	Business   Business   `gorm:"foreignKey:BusinessID" json:"business"`
	Role       string     `gorm:"not null;default:'OWNER'" json:"role"`
	IsVerified bool       `gorm:"default:false" json:"is_verified"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

type Business struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Type      string    `gorm:"not null" json:"type"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	LogoURL   string    `json:"logo_url"`
	Users     []User    `gorm:"foreignKey:BusinessID" json:"users"`
	CreatedAt time.Time `json:"created_at"`
}

type VerificationCode struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"index;not null" json:"email"`
	Code      string    `gorm:"not null" json:"code"`
	Type      string    `gorm:"not null" json:"type"` // e.g., 'REGISTER', 'FORGOT_PASSWORD'
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type UserResponse struct {
	ID              uuid.UUID `json:"id"`
	Email           string    `json:"email"`
	BusinessID      uuid.UUID `json:"business_id"`
	BusinessName    string    `json:"business_name"`
	BusinessType    string    `json:"business_type"`
	BusinessAddress string    `json:"business_address"`
	BusinessPhone   string    `json:"business_phone"`
	BusinessLogoURL string    `json:"business_logo_url"`
	Role            string    `json:"role"`
	IsVerified      bool      `json:"is_verified"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserRepository interface {
	Create(ctx context.Context, u *User, businessName string) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByBusinessID(ctx context.Context, businessID uuid.UUID) ([]User, error)
	AddUser(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error
}

type VerificationCodeRepository interface {
	Create(ctx context.Context, vc *VerificationCode) error
	GetLastByEmail(ctx context.Context, email, codeType string) (*VerificationCode, error)
	DeleteByEmail(ctx context.Context, email, codeType string) error
}

type BusinessRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Business, error)
	Update(ctx context.Context, b *Business) error
}

type AuthUsecase interface {
	Login(ctx context.Context, email string, password string) (*TokenResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*TokenResponse, error)
	Register(ctx context.Context, email, password, bizName string) error
	VerifyOTP(ctx context.Context, email, code string) (*TokenResponse, error)
	ResendOTP(ctx context.Context, email string) error
	GetProfile(ctx context.Context, userID uuid.UUID) (*UserResponse, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, email, code, newPassword string) (*TokenResponse, error)
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type StaffUsecase interface {
	CreateStaff(ctx context.Context, email, password, role string, businessID uuid.UUID) error
	GetStaff(ctx context.Context, businessID uuid.UUID) ([]UserResponse, error)
}

type BusinessUsecase interface {
	UpdateBusiness(ctx context.Context, businessID uuid.UUID, req UpdateBusinessRequest) error
}

type UpdateBusinessRequest struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	LogoURL string `json:"logo_url"`
}
