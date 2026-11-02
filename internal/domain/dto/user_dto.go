package dto

import (
	"strings"
	"time"

	domain_user "github.com/shodruzhoshimzoda/tojtech/internal/domain/user"
)


type AdminUserResponse struct {
	UUID 	 	string  `json:"uuid"`
	Role     	string `json:"role"`
	Email    	string `json:"email"`
	CreatedAt 	time.Time `json:"created_at"`
	UpdatedAt 	time.Time `json:"updated_at"`
}


func NewAdminUserResponse (u domain_user.User ) AdminUserResponse {
	return AdminUserResponse{
		UUID: u.UUID.String(),
		Role: u.Role,
		Email: u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func NewAdminUserResponses(users []domain_user.User) []AdminUserResponse {
	items := make([]AdminUserResponse, 0, len(users))

	for _, u := range users {
		
		items = append(items, NewAdminUserResponse(u))
	}

	return items
}

type ChangeRoleDTO struct {
	Role string `json:"role"`
}

// Validate - cleans the value and checks that the role exists
func (c *ChangeRoleDTO) Validate() error {
	c.Role = strings.ToLower(strings.TrimSpace(c.Role))

	if !domain_user.IsValidRole(c.Role) {
		return domain_user.ErrInvalidRole
	}

	return nil
}