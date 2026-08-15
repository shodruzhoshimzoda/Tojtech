package product_usecase

import (
	"context"

	"github.com/google/uuid"
	domain_product "github.com/shodruzhoshimzoda/tojtech/internal/domain/product"
)

type ProductRepo interface {
	GetProductByUUID(ctx context.Context, uuid uuid.UUID) (*domain_product.Product, error)
	ProductList(ctx context.Context) ([]*domain_product.Product, error)
	CreateProduct(ctx context.Context, product *domain_product.Product) error
	DeleteProduct(ctx context.Context, uuid uuid.UUID) error
	UpdateProduct(ctx context.Context, product *domain_product.Product) error
	DeleteProductImage(ctx context.Context, productUUID, imageUUID uuid.UUID) error
	AddProductImage(ctx context.Context, productUUID uuid.UUID, imageURL string, isMain bool) (*domain_product.ProductImage, error)
	CountProductImages(ctx context.Context, productUUID uuid.UUID) (int, error)
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

func (p *ProductUsecase) GetProduct(ctx context.Context, id uuid.UUID) (*domain_product.Product, error) {

	return p.repo.GetProductByUUID(ctx, id)

}

func (p *ProductUsecase) ProductList(ctx context.Context) ([]*domain_product.Product, error) {

	return p.repo.ProductList(ctx)
}

func (p *ProductUsecase) CreateProduct(ctx context.Context, prod *domain_product.Product) error {
	return p.repo.CreateProduct(ctx, prod)
}

func (p *ProductUsecase) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return p.repo.DeleteProduct(ctx, id)
}

func (p *ProductUsecase) UpdateProduct(ctx context.Context, prod *domain_product.Product) error {
	return p.repo.UpdateProduct(ctx, prod)
}

const maxImagesPerProduct = 10

func (p *ProductUsecase) AddProductImage(ctx context.Context, productUUID uuid.UUID, imageURL string, isMain bool) (*domain_product.ProductImage, error) {
	if err := domain_product.ValidateImageURL(imageURL); err != nil {
		return nil, err
	}

	count, err := p.repo.CountProductImages(ctx, productUUID)
	if err != nil {
		return nil, err
	}
	if count >= maxImagesPerProduct {
		return nil, domain_product.ErrTooManyImages
	}

	return p.repo.AddProductImage(ctx, productUUID, imageURL, isMain)
}

func (p *ProductUsecase) DeleteProductImage(ctx context.Context, productUUID, imageUUID uuid.UUID) error {
	return p.repo.DeleteProductImage(ctx, productUUID, imageUUID)
}
