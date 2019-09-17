package user

import (
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type loginRequest struct {
	Username string `form:"Username" validate:"required"`
	Password string `form:"Password" validate:"required"`
}

type validUpdatedRequest struct {
	Username string `form:"Username" validate:"omitempty,min=3"`
	Password string `form:"Password" validate:"omitempty,min=6"`
	Role     string `form:"Role" validate:"omitempty,oneof=Admin User Guest"`
	Mail     string `form:"Mail" validate:"omitempty,email"`
	Active   bool   `form:"Active" validate:"omitempty"`
}

type updateUserRequest struct {
	validUpdatedRequest `json:"user"`
}

type validNewUser struct {
	Username string `form:"Username" validate:"required,min=3"`
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
	if role != "Admin" && r.Role != "" {
		if r.Role != u.Role {
			return u, fmt.Errorf("Only Admin can update user's Role")
		}
	} else if role == "Admin" && r.Role != "" {
		u.Role = r.Role
	}

	// Only the Admin must be able to update users Active state
	if r.Active != u.Active {
		if role != "Admin" {
			return u, fmt.Errorf("Only Admin can update user's Active state")
		} else {
			u.Active = r.Active
		}
	}

	// Update the username making sure is NOT taken
	var testUser User
	if err := testUser.ByUsername(r.Username); err == nil {
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
	u.Active = true

	return u
}
