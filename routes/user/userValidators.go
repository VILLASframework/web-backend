package user

import (
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type loginRequest struct {
	Username string `form:"Username" validate:"required"`
	Password string `form:"Password" validate:"required,min=6"`
}

type updateUserRequest struct {
	Username string `form:"Username" validate:"omitempty"`
	Password string `form:"Password" validate:"min=6"`
	Role     string `form:"Role" validate:"omitempty,oneof=Admin User Guest"`
	Mail     string `form:"Mail" validate:"omitempty,email"`
}

type addUserRequest struct {
	Username string `form:"Username" validate:"required"`
	Password string `form:"Password" validate:"required,min=6"`
	Role     string `form:"Role" validate:"required,oneof=Admin User Guest"`
	Mail     string `form:"Mail" validate:"required,email"`
}

func (r *loginRequest) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
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
