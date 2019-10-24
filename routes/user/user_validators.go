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
	Username    string `form:"Username" validate:"omitempty,min=3"`
	Password    string `form:"Password" validate:"omitempty,min=6"`
	OldPassword string `form:"OldPassword" validate:"omitempty,min=6"`
	Role        string `form:"Role" validate:"omitempty,oneof=Admin User Guest"`
	Mail        string `form:"Mail" validate:"omitempty,email"`
	Active      string `form:"Active" validate:"omitempty,oneof=yes no"`
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
	if errs != nil {
		return errs
	}

	return nil
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
	if (r.Active == "yes" && u.Active == false) || (r.Active == "no" && u.Active == true) {
		if role != "Admin" {
			return u, fmt.Errorf("Only Admin can update user's Active state")
		} else {
			u.Active = !u.Active
		}
	}

	// Update the username making sure it is NOT taken
	var testUser User
	if err := testUser.ByUsername(r.Username); err == nil {
		return u, fmt.Errorf("Username is alreaday taken")
	}

	if r.Username != "" {
		u.Username = r.Username
	}

	// If there is a new password then hash it and update it
	if r.Password != "" {
		if role != "Admin" { // if requesting user is NOT admin, old password needs to be validated

			if r.OldPassword == "" {
				return u, fmt.Errorf("old password is missing in request")
			}

			err := oldUser.validatePassword(r.OldPassword)
			if err != nil {
				return u, fmt.Errorf("previous password not correct, pw not changed")
			}
		}

		err := u.setPassword(r.Password)
		if err != nil {
			return u, fmt.Errorf("unable to encrypt new password")
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
