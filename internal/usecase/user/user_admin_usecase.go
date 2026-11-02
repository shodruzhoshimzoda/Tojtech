package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/shodruzhoshimzoda/tojtech/internal/domain/dto"
	domain_user "github.com/shodruzhoshimzoda/tojtech/internal/domain/user"
)



type AdminUserRepo interface {
	GetUsers(ctx context.Context) ([]domain_user.User, error) 
	UpdateUserRole(ctx context.Context, userUUID uuid.UUID, role string) error
}

type AdminUserUsecase struct{
	repo AdminUserRepo
}


func NewAdminUserUsecase(repo AdminUserRepo) *AdminUserUsecase {
	return &AdminUserUsecase{
		repo: repo,
	}
}

func (u *AdminUserUsecase) GetUsers(ctx context.Context)( []dto.AdminUserResponse, error ){

	users, err := u.repo.GetUsers(ctx)

	if err != nil {
		return nil, err
	}

	return dto.NewAdminUserResponses(users), nil

}

func (u *AdminUserUsecase) ChangeRole(ctx context.Context, callerUUID, targetUUID uuid.UUID, req dto.ChangeRoleDTO) error {
	// 1. the role must exist (also cleans spaces and letter case)
	if err := req.Validate(); err != nil {
		return err
	}

	// 2. an admin cannot change their own role, so the panel is never left without an admin
	if callerUUID == targetUUID {
		return domain_user.ErrCannotChangeOwnRole
	}

	// 3. the storage reports domain_user.ErrUserNotFound if there is no such user
	return u.repo.UpdateUserRole(ctx, targetUUID, req.Role)
}