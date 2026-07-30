package usecase

import (
	"context"
	"errors"

	domain "github.com/shodruzhoshimzoda/tojtech/internal/domain/product"
)

var (
	ErrProductNotFound = errors.New("product not  found")
)

type ProductRepo interface {
	GetById(ctx context.Context, id int64) (*domain.Product, error)
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
func (p *ProductUsecase) GetProduct(ctx context.Context, id int64) (*domain.Product, error) {

	prod, err := p.repo.GetById(ctx, int64(id))

	if err != nil {
		return nil, err
	}

	return prod, nil

}
