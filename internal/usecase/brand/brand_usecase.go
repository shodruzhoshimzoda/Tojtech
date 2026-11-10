package brand

import (
	"context"

	"github.com/google/uuid"
	domain_brand "github.com/shodruzhoshimzoda/tojtech/internal/domain/brand"
	"github.com/shodruzhoshimzoda/tojtech/internal/domain/dto"
)

// BrandRepo - what the usecase needs from the database layer.
// The real implementation is repo_brand.BrandRepository.
type BrandRepo interface {
	GetBrandByUUID(ctx context.Context, id uuid.UUID) (*domain_brand.Brand, error)
	CreateBrand(ctx context.Context, b *domain_brand.Brand) error
	UpdateBrand(ctx context.Context, id uuid.UUID, b *domain_brand.Brand) (*domain_brand.Brand, error)
	DeleteBrand(ctx context.Context, id uuid.UUID) error
	GetBrands(ctx context.Context) ([]*domain_brand.Brand, error)
	GetBrandsWithProducts(ctx context.Context, categoryUUID *uuid.UUID) ([]*domain_brand.Brand, error)
}

type BrandUsecase struct {
	repo BrandRepo
}

func NewBrandUsecase(repo BrandRepo) *BrandUsecase {
	return &BrandUsecase{repo: repo}
}

// GetBrands - for the admin panel: all brands, product_count = all their products.
func (u *BrandUsecase) GetBrands(ctx context.Context) ([]dto.BrandResponse, error) {
	brands, err := u.repo.GetBrands(ctx)
	if err != nil {
		return nil, err
	}
	return dto.NewBrandResponses(brands), nil
}

// GetShopBrands - for the filter panel of the shop: only brands with active products.
// categoryUUID == nil means "all categories".
func (u *BrandUsecase) GetShopBrands(ctx context.Context, categoryUUID *uuid.UUID) ([]dto.BrandResponse, error) {
	brands, err := u.repo.GetBrandsWithProducts(ctx, categoryUUID)
	if err != nil {
		return nil, err
	}
	return dto.NewBrandResponses(brands), nil
}

func (u *BrandUsecase) GetBrand(ctx context.Context, id uuid.UUID) (dto.BrandResponse, error) {
	b, err := u.repo.GetBrandByUUID(ctx, id)
	if err != nil {
		return dto.BrandResponse{}, err
	}
	return dto.NewBrandResponse(b), nil
}

func (u *BrandUsecase) CreateBrand(ctx context.Context, req dto.BrandRequest) (dto.BrandResponse, error) {
	b := req.ToBrand()
	b.Normalize() // trim spaces, lowercase slug, make the slug from the name if it is empty
	if err := b.Validate(); err != nil {
		return dto.BrandResponse{}, err
	}

	if err := u.repo.CreateBrand(ctx, b); err != nil {
		return dto.BrandResponse{}, err
	}
	return dto.NewBrandResponse(b), nil
}

func (u *BrandUsecase) UpdateBrand(ctx context.Context, id uuid.UUID, req dto.BrandRequest) (dto.BrandResponse, error) {
	b := req.ToBrand()
	b.Normalize()
	if err := b.Validate(); err != nil {
		return dto.BrandResponse{}, err
	}

	updated, err := u.repo.UpdateBrand(ctx, id, b)
	if err != nil {
		return dto.BrandResponse{}, err
	}
	return dto.NewBrandResponse(updated), nil
}

// DeleteBrand deletes the brand; its products stay and just lose the brand.
// The "are you sure, N products use it" question is asked by the admin panel
// (it knows product_count from the list), the server does not block the deletion.
func (u *BrandUsecase) DeleteBrand(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteBrand(ctx, id)
}