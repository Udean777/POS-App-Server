package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID            uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BusinessID    uuid.UUID         `gorm:"type:uuid;not null;index" json:"business_id"`
	StaffID       uuid.UUID         `gorm:"type:uuid;not null" json:"staff_id"`
	TotalAmount   float64           `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	AmountPaid    float64           `gorm:"type:decimal(10,2);not null" json:"amount_paid"`
	Change        float64           `gorm:"type:decimal(10,2);not null" json:"change"`
	PaymentMethod string            `gorm:"not null" json:"payment_method"` // CASH, QRIS
	Status        string            `gorm:"not null;default:'COMPLETED'" json:"status"`
	Items         []TransactionItem `gorm:"foreignKey:TransactionID" json:"items"`
	Staff         *User             `gorm:"foreignKey:StaffID" json:"staff,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
}

type TransactionItem struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TransactionID uuid.UUID `gorm:"type:uuid;not null" json:"transaction_id"`
	ProductID     uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	VariantID     uuid.UUID `gorm:"type:uuid;not null" json:"variant_id"`
	Quantity      int       `gorm:"not null" json:"quantity"`
	Price         float64   `gorm:"type:decimal(10,2);not null" json:"price"` // Snapshot harganya saat beli
	Subtotal      float64   `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	Product       *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Variant       *Variant  `gorm:"foreignKey:VariantID" json:"variant,omitempty"`
}

// Struct untuk mempermudah request dari frontend (Request Payload)
type CheckoutRequest struct {
	PaymentMethod string                 `json:"payment_method" binding:"required"`
	AmountPaid    float64                `json:"amount_paid" binding:"required"`
	Items         []CheckoutItemRequest `json:"items" binding:"required,min=1"`
}

type CheckoutItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	VariantID uuid.UUID `json:"variant_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *Transaction) error
	GetByBusinessID(ctx context.Context, businessID uuid.UUID) ([]Transaction, error)
}

type TransactionUsecase interface {
	ProcessCheckout(ctx context.Context, req CheckoutRequest, businessID uuid.UUID, staffID uuid.UUID) (*Transaction, error)
	GetTransactions(ctx context.Context, businessID uuid.UUID) ([]Transaction, error)
}
