package usecase_category

import (
	"context"

	"github.com/google/uuid"
	domain_category "github.com/shodruzhoshimzoda/tojtech/internal/domain/category"
)

type CategoryRepo interface {
	GetCategoryByUUID(ctx context.Context, uuid uuid.UUID) (*domain_category.Category, error)
	CreateCategory(ctx context.Context, category *domain_category.Category) (uuid.UUID, error)
}

type CategoryUseCase struct {
	categoryRepo CategoryRepo
}

func NewCategoryUseCase(categoryRepo CategoryRepo) *CategoryUseCase {
	return &CategoryUseCase{
		categoryRepo: categoryRepo,
	}
}

func (u *CategoryUseCase) GetCategory(ctx context.Context, uuid uuid.UUID) (*domain_category.Category, error) {
	return u.categoryRepo.GetCategoryByUUID(ctx, uuid)
}

func (u *CategoryUseCase) CreateCategory(ctx context.Context, category *domain_category.Category) (uuid.UUID, error) {

	return u.categoryRepo.CreateCategory(ctx, category)

}
