package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	domain "github.com/shodruzhoshimzoda/tojtech/internal/domain/product"
)

var (
	ErrProductNotFound = errors.New("product not  found")
)

type ProductRepo interface {
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.Product, error)
	List(ctx context.Context) ([]domain.Product, error)
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

	prod, err := p.repo.GetByUUID(ctx, id)

	if err != nil {
		return nil, err
	}

	return prod, nil

}

func (p *ProductUsecase) List(ctx context.Context) ([]domain.Product, error) {

	products, err := p.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}
