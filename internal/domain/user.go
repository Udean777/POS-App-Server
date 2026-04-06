package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email      string    `gorm:"uniqueIndex;not null" json:"email"`
	Password   string    `gorm:"not null" json:"-"`
	BusinessID uuid.UUID `gorm:"type:uuid;not null" json:"business_id"`
	Business   Business  `gorm:"foreignKey:BusinessID" json:"business"`
	CreatedAt  time.Time `json:"created_at"`
}

type Business struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Type      string    `gorm:"not null" json:"type"`
	Address   string    `json:"address"`
	Users     []User    `gorm:"foreignKey:BusinessID" json:"users"`
	CreatedAt time.Time `json:"created_at"`
}

type UserResponse struct {
	ID              uuid.UUID `json:"id"`
	Email           string    `json:"email"`
	BusinessID      uuid.UUID `json:"business_id"`
	BusinessName    string    `json:"business_name"`
	BusinessType    string    `json:"business_type"`
	BusinessAddress string    `json:"business_address"`
	Role            string    `json:"role"`
}

type UserRepository interface {
	Create(ctx context.Context, u *User, businessName string) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
}

type AuthUsecase interface {
	Login(ctx context.Context, email string, password string) (string, error)
	Register(ctx context.Context, email, password, bizName string) error
	GetProfile(ctx context.Context, userID uuid.UUID) (*UserResponse, error)
}
