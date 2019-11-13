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
	User validUpdatedRequest `json:"user"`
}

type validNewUser struct {
	Username string `form:"Username" validate:"required,min=3"`
	Password string `form:"Password" validate:"required,min=6"`
	Role     string `form:"Role" validate:"required,oneof=Admin User Guest"`
	Mail     string `form:"Mail" validate:"required,email"`
}

type addUserRequest struct {
	User validNewUser `json:"user"`
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

func (r *updateUserRequest) updatedUser(callerID interface{}, role interface{}, oldUser User) (User, error) {

	// Use the old User as a basis for the updated User `u`
	u := oldUser

	// Only the Admin must be able to update user's role
	if role != "Admin" && r.User.Role != "" {
		if r.User.Role != u.Role {
			return u, fmt.Errorf("Only Admin can update user's Role")
		}
	} else if role == "Admin" && r.User.Role != "" {
		u.Role = r.User.Role
	}

	// Only the Admin must be able to update users Active state
	if (r.User.Active == "yes" && u.Active == false) || (r.User.Active == "no" && u.Active == true) {
		if role != "Admin" {
			return u, fmt.Errorf("Only Admin can update user's Active state")
		} else {
			u.Active = !u.Active
		}
	}

	// Update the username making sure it is NOT taken
	var testUser User
	if err := testUser.ByUsername(r.User.Username); err == nil {
		return u, fmt.Errorf("Username is alreaday taken")
	}

	if r.User.Username != "" {
		u.Username = r.User.Username
	}

	// If there is a new password then hash it and update it
	if r.User.Password != "" {

		if r.User.OldPassword == "" { // admin or old password has to be present for pw change
			return u, fmt.Errorf("old or admin password is missing in request")
		}

		if role == "Admin" { // admin has to enter admin password
			var adminUser User
			err := adminUser.ByID(callerID.(uint))
			if err != nil {
				return u, err
			}

			err = adminUser.validatePassword(r.User.OldPassword)
			if err != nil {
				return u, fmt.Errorf("admin password not correct, pw not changed")
			}

		} else { //normal or guest user has to enter old password

			err := oldUser.validatePassword(r.User.OldPassword)
			if err != nil {
				return u, fmt.Errorf("previous password not correct, pw not changed")
			}
		}

		err := u.setPassword(r.User.Password)
		if err != nil {
			return u, fmt.Errorf("unable to encrypt new password")
		}
	}

	// Update mail
	if r.User.Mail != "" {
		u.Mail = r.User.Mail
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

	u.Username = r.User.Username
	u.Password = r.User.Password
	u.Mail = r.User.Mail
	u.Role = r.User.Role
	u.Active = true

	return u
}
