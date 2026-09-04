package product_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/shodruzhoshimzoda/tojtech/internal/domain/dto"
	product_domain "github.com/shodruzhoshimzoda/tojtech/internal/domain/product"
)

type ProductRepo interface {
	GetProduct(ctx context.Context, uuid uuid.UUID) (*product_domain.Product, error)
	GetProducts(ctx context.Context) ([]*product_domain.Product, error)
	CreateProduct(ctx context.Context, product *product_domain.Product) error
	DeleteProduct(ctx context.Context, uuid uuid.UUID) error
	UpdateProduct(ctx context.Context, product *product_domain.Product) error
	DeleteProductImage(ctx context.Context, productUUID, imageUUID uuid.UUID) error
	AddProductImage(ctx context.Context, productUUID uuid.UUID, imageURL string, isMain bool) (*product_domain.ProductImage, error)
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

func (p *ProductUsecase) GetProduct(ctx context.Context, id uuid.UUID) (*product_domain.Product, error) {

	return p.repo.GetProduct(ctx, id)

}

func (p *ProductUsecase) ProductList(ctx context.Context) ([]dto.ProductDTO, error) {

	products, err := p.repo.GetProducts(ctx)
	if err != nil {
		return nil, err
	}

	var productSlice = make([]dto.ProductDTO, 0, len(products))
	for _, p := range products {
		prodDTO := dto.NewProductDTO(p)
		productSlice = append(productSlice, prodDTO)
	}

	return productSlice, err
}

func (p *ProductUsecase) CreateProduct(ctx context.Context, prod *product_domain.Product) error {
	return p.repo.CreateProduct(ctx, prod)
}

func (p *ProductUsecase) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return p.repo.DeleteProduct(ctx, id)
}

func (p *ProductUsecase) UpdateProduct(ctx context.Context, prod *product_domain.Product) error {
	return p.repo.UpdateProduct(ctx, prod)
}

const maxImagesPerProduct = 10

func (p *ProductUsecase) AddProductImage(ctx context.Context, productUUID uuid.UUID, imageURL string, isMain bool) (*product_domain.ProductImage, error) {
	if err := product_domain.ValidateImageURL(imageURL); err != nil {
		return nil, err
	}

	count, err := p.repo.CountProductImages(ctx, productUUID)
	if err != nil {
		return nil, err
	}
	if count >= maxImagesPerProduct {
		return nil, product_domain.ErrTooManyImages
	}

	return p.repo.AddProductImage(ctx, productUUID, imageURL, isMain)
}

func (p *ProductUsecase) DeleteProductImage(ctx context.Context, productUUID, imageUUID uuid.UUID) error {
	return p.repo.DeleteProductImage(ctx, productUUID, imageUUID)
}
