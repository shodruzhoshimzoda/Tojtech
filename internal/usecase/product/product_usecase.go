package usecase_product

import (
	"context"

	"github.com/google/uuid"
	domain_product "github.com/shodruzhoshimzoda/tojtech/internal/domain/product"
)

type ProductRepo interface {
	GetProductByUUID(ctx context.Context, uuid uuid.UUID) (*domain_product.Product, error)
	ProductList(ctx context.Context) ([]*domain_product.Product, error)
	CreateProduct(ctx context.Context, product *domain_product.Product) error
}

// ProductUsecase struct contain repository which means its doesn't matter which DB we use this
// Db must have this methods
type ProductUsecase struct {
	repo ProductRepo
}

func NewProductUsecase(repo ProductRepo) *ProductUsecase {
	return &ProductUsecase{
		repo: repo,
	}
}

// GetProduct Business logic for
func (p *ProductUsecase) GetProduct(ctx context.Context, id uuid.UUID) (*domain_product.Product, error) {

	return p.repo.GetProductByUUID(ctx, id)

}

func (p *ProductUsecase) ProductList(ctx context.Context) ([]*domain_product.Product, error) {

	return p.repo.ProductList(ctx)
}

func (p *ProductUsecase) CreateProduct(ctx context.Context, prod *domain_product.Product) error {
	return p.repo.CreateProduct(ctx, prod)
}
