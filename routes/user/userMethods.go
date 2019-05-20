package user

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"golang.org/x/crypto/bcrypt"
)

// TODO: use validator
type Credentials struct {
	Username string `form:"Username"`
	Password string `form:"Password"`
	Role     string `form:"Role"`
	Mail     string `form:"Mail"`
}

// This is ugly but no other way to keep each model on the corresponding
// package since we have circular dependencies. Methods of a type must
// live in the same package.
type User struct {
	Model common.User
}

func FindUserByUsername(username string) (User, error) {
	db := common.GetDB()
	var user User
	err := db.Find(&user.Model, "Username = ?", username).Error
	return user, err
}

func (u *User) setPassword(password string) error {
	// TODO: Not implemented
	return nil
}

func (u *User) validatePassword(password string) error {
	loginPassword := []byte(password)
	hashedPassword := []byte(u.Model.Password)
	return bcrypt.CompareHashAndPassword(hashedPassword, loginPassword)
}

func (u *User) update(data interface{}) error {
	// TODO: Not implemented
	return nil
}
