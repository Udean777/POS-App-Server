package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sajudin/pos-app-server/internal/domain"
)

type productUsecase struct {
	productRepo domain.ProductRepository
}

func NewProductUsecase(pr domain.ProductRepository) domain.ProductUsecase {
	return &productUsecase{productRepo: pr}
}

func (u *productUsecase) AddProduct(ctx context.Context, p *domain.Product) error {
	if p.Name == "" {
		return errors.New("nama produk tidak boleh kosong")
	}

	if len(p.Variants) == 0 {
		return errors.New("produk harus memiliki minimal satu varian")
	}

	for _, v := range p.Variants {
		if v.Name == "" {
			return errors.New("nama varian tidak boleh kosong")
		}
		if v.Price < 0 {
			return errors.New("harga tidak boleh negatif")
		}
		if v.Stock < 0 {
			return errors.New("stok tidak boleh negatif")
		}
		if v.SKU == "" {
			return errors.New("SKU tidak boleh kosong")
		}
	}

	return u.productRepo.Create(ctx, p)
}

func (u *productUsecase) GetAllProducts(ctx context.Context, businessID uuid.UUID) ([]domain.Product, error) {
	return u.productRepo.Fetch(ctx, businessID)
}

func (u *productUsecase) GetProductByID(ctx context.Context, id uuid.UUID, businessID uuid.UUID) (*domain.Product, error) {
	return u.productRepo.GetByID(ctx, id, businessID)
}

func (u *productUsecase) UpdateProduct(ctx context.Context, p *domain.Product) error {
	if p.Name == "" {
		return errors.New("nama produk tidak boleh kosong")
	}
	return u.productRepo.Update(ctx, p)
}

func (u *productUsecase) DeleteProduct(ctx context.Context, id uuid.UUID, businessID uuid.UUID) error {
	return u.productRepo.Delete(ctx, id, businessID)
}
