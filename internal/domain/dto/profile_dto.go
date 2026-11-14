package dto

import (
	"strings"
	domain_user "github.com/shodruzhoshimzoda/tojtech/internal/domain/user"


)



type UpdateProfileDTO struct {
	Username 	string `json:"username"`
}


// Type for changing password
type ChangePasswordDTO struct {
	CurrentPassword  	string  	`json:"current_password"`
	NewPassword 		string	    `json:"new_password"`
}


//  ---------- Validate - these constructs are implemented for  username and password validation ------------ 


func (d *ChangePasswordDTO) Validate() error {
	if d.CurrentPassword == "" || d.NewPassword == "" {
		return ErrPasswordRequired
	}
	return domain_user.ValidatePassword(d.NewPassword)
}




func (d *UpdateProfileDTO) Validate() error {

	d.Username = strings.TrimSpace(d.Username)

	return domain_user.ValidateUsername(d.Username)

}

func NewUserResponse(u domain_user.User) UserResponseDTO {
	return UserResponseDTO{
		UUID: u.UUID.String(),
		Email: u.Email,
		Username: u.Username,
		Role: u.Role,
	}
}







