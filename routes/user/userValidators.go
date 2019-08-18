package user

import (
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type loginRequest struct {
	Username string `form:"Username" validate:"required"`
	Password string `form:"Password" validate:"required,min=6"`
}

type validUpdatedRequest struct {
	Username string `form:"Username" validate:"omitempty"`
	Password string `form:"Password" validate:"omitempty,min=6"`
	Role     string `form:"Role" validate:"omitempty,oneof=Admin User Guest"`
	Mail     string `form:"Mail" validate:"omitempty,email"`
}

type updateUserRequest struct {
	validUpdatedRequest `json:"user"`
}

type validNewUser struct {
	Username string `form:"Username" validate:"required"`
	Password string `form:"Password" validate:"required,min=6"`
	Role     string `form:"Role" validate:"required,oneof=Admin User Guest"`
	Mail     string `form:"Mail" validate:"required,email"`
}

type addUserRequest struct {
	validNewUser `json:"user"`
}

func (r *loginRequest) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *updateUserRequest) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *updateUserRequest) updatedUser(role interface{},
	oldUser User) (User, error) {

	// Use the old User as a basis for the updated User `u`
	u := oldUser

	// Only the Admin must be able to update user's role
	if role != "Admin" && r.Role != u.Role {
		return u, fmt.Errorf("Only Admin can update user's Role")
	} else if role == "Admin" && r.Role != "" {
		u.Role = r.Role
	}

	// Update the username making sure is NOT taken
	if err := u.ByUsername(r.Username); err == nil {
		return u, fmt.Errorf("Username is alreaday taken")
	}

	if r.Username != "" {
		u.Username = r.Username
	}

	// If there is a new password then hash it and update it
	if r.Password != "" {
		err := u.setPassword(r.Password)
		if err != nil {
			return u, fmt.Errorf("Unable to encrypt new password")
		}
	}

	// Update mail
	if r.Mail != "" {
		u.Mail = r.Mail
	}

	return u, nil
}

func (r *addUserRequest) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *addUserRequest) createUser() User {
	var u User

	u.Username = r.Username
	u.Password = r.Password
	u.Mail = r.Mail
	u.Role = r.Role

	return u
}
