package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sajudin/pos-app-server/internal/domain"
)

type transactionUsecase struct {
	txRepo      domain.TransactionRepository
	productRepo domain.ProductRepository
}

func NewTransactionUsecase(tr domain.TransactionRepository, pr domain.ProductRepository) domain.TransactionUsecase {
	return &transactionUsecase{
		txRepo:      tr,
		productRepo: pr,
	}
}

func (u *transactionUsecase) ProcessCheckout(ctx context.Context, req domain.CheckoutRequest, businessID uuid.UUID, staffID uuid.UUID) (*domain.Transaction, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("keranjang belanja kosong")
	}

	var transactionItems []domain.TransactionItem
	var totalAmount float64

	// Validasi dan Hitung Total (Di Backend agar aman)
	for _, reqItem := range req.Items {
		// Ambil produk dan varian dari database
		product, err := u.productRepo.GetByID(ctx, reqItem.ProductID, businessID)
		if err != nil {
			return nil, errors.New("produk tidak valid atau tidak ditemukan")
		}

		var selectedVariant *domain.Variant
		for _, v := range product.Variants {
			if v.ID == reqItem.VariantID {
				selectedVariant = &v
				break
			}
		}

		if selectedVariant == nil {
			return nil, errors.New("varian tidak valid")
		}

		// Kalkulasi subtotal
		subtotal := selectedVariant.Price * float64(reqItem.Quantity)
		totalAmount += subtotal

		transactionItems = append(transactionItems, domain.TransactionItem{
			ProductID: product.ID,
			VariantID: selectedVariant.ID,
			Quantity:  reqItem.Quantity,
			Price:     selectedVariant.Price, // Snapshot price
			Subtotal:  subtotal,
		})
	}

	// Cek pembayaran
	if req.AmountPaid < totalAmount {
		return nil, errors.New("jumlah pembayaran kurang dari total tagihan")
	}

	change := req.AmountPaid - totalAmount

	tx := &domain.Transaction{
		BusinessID:    businessID,
		StaffID:       staffID,
		TotalAmount:   totalAmount,
		AmountPaid:    req.AmountPaid,
		Change:        change,
		PaymentMethod: req.PaymentMethod,
		Status:        "COMPLETED",
		Items:         transactionItems,
	}

	// Simpan transaksi (stok akan berkurang di layer repo)
	if err := u.txRepo.Create(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (u *transactionUsecase) GetTransactions(ctx context.Context, businessID uuid.UUID) ([]domain.Transaction, error) {
	return u.txRepo.GetByBusinessID(ctx, businessID)
}
