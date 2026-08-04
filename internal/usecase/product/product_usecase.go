package usecase_product

import (
	"context"

	"github.com/google/uuid"
	domain "github.com/shodruzhoshimzoda/tojtech/internal/domain/product"
)

type ProductRepo interface {
	GetProductByUUID(ctx context.Context, uuid uuid.UUID) (*domain.Product, error)
	ProductList(ctx context.Context) ([]*domain.Product, error)
}

// Productusecase struct contain repository wcich means its doesnt matter which DB we use this
// Db must have this methods
type ProductUsecase struct {
	repo ProductRepo
}

func NewProductUsecase(repo ProductRepo) *ProductUsecase {
	return &ProductUsecase{
		repo: repo,
	}
}

// Bussines logic for GetProduct
func (p *ProductUsecase) GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error) {

	return p.repo.GetProductByUUID(ctx, id)

}

func (p *ProductUsecase) ProductList(ctx context.Context) ([]*domain.Product, error) {

	return p.repo.ProductList(ctx)
}
