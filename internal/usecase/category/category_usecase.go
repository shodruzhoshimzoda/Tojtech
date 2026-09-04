package usecase_category

import (
	"context"

	"github.com/google/uuid"
	category_domain "github.com/shodruzhoshimzoda/tojtech/internal/domain/category"
	"github.com/shodruzhoshimzoda/tojtech/internal/domain/dto"
)

type CategoryRepo interface {
	GetCategory(ctx context.Context, uuid uuid.UUID) (*category_domain.Category, error)
	GetCategories(ctx context.Context) ([]*category_domain.Category, error)
	CreateCategory(ctx context.Context, category *category_domain.Category) (uuid.UUID, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, category *category_domain.Category) (cat *category_domain.Category, e error)
	DeleteCategory(ctx context.Context, id uuid.UUID) (e error)
}

type CategoryUseCase struct {
	categoryRepo CategoryRepo
}

func NewCategoryUseCase(categoryRepo CategoryRepo) *CategoryUseCase {
	return &CategoryUseCase{
		categoryRepo: categoryRepo,
	}
}

func (u *CategoryUseCase) GetCategory(ctx context.Context, rawUUID string) (dto.CategoryResponse, error) {

	id, err := uuid.Parse(rawUUID)
	if err != nil {
		return dto.CategoryResponse{}, category_domain.ErrInvalidUUID
	}

	category, err := u.categoryRepo.GetCategory(ctx, id)
	if err != nil {
		return dto.CategoryResponse{}, err
	}

	cat := dto.NewCategoryResponse(category)
	return cat, nil

	// return u.categoryRepo.GetCategory(ctx, uuid)
}

func (u *CategoryUseCase) CreateCategory(ctx context.Context, category *category_domain.Category) (dto.CategoryResponse, error) {

	id, err := u.categoryRepo.CreateCategory(ctx, category)

	if err != nil {
		return dto.CategoryResponse{}, err
	}
	category.UUID = id
	catDTO := dto.NewCategoryResponse(category)

	return catDTO, nil

}

func (u *CategoryUseCase) UpdateCategory(ctx context.Context, id uuid.UUID, category *category_domain.Category) (*category_domain.Category, error) {

	if err := category.Validate(); err != nil {
		return nil, err
	}

	return u.categoryRepo.UpdateCategory(ctx, id, category)

	// return u.categoryRepo.UpdateCategory(ctx, rawUUID, category)

}

func (u *CategoryUseCase) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return u.categoryRepo.DeleteCategory(ctx, id)
}

func (u *CategoryUseCase) GetCategories(ctx context.Context) ([]dto.CategoryResponse, error) {

	categories, err := u.categoryRepo.GetCategories(ctx)

	if err != nil {
		return nil, err
	}
	var categoriesDTO = make([]dto.CategoryResponse, 0, len(categories)) // capacity its equel to lenth of categories

	for _, cat := range categories {

		category := dto.NewCategoryResponse(cat)

		categoriesDTO = append(categoriesDTO, category)
	}

	return categoriesDTO, nil

}
