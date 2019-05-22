package user

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10

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
	common.User // check golang embedding types
}

func FindUserByUsername(username string) (User, error) {
	db := common.GetDB()
	var user User
	err := db.Find(&user, "Username = ?", username).Error
	return user, err
}

func (u *User) SetPassword(password string) error {
	if len(password) == 0 {
		return fmt.Errorf("Password cannot be empty")
	}
	newPassword, err :=
		bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return fmt.Errorf("Failed to generate hash from password")
	}
	u.Password = string(newPassword)
	return nil
}

func (u *User) validatePassword(password string) error {
	loginPassword := []byte(password)
	hashedPassword := []byte(u.Password)
	return bcrypt.CompareHashAndPassword(hashedPassword, loginPassword)
}

func (u *User) update(data interface{}) error {
	// TODO: Not implemented
	return nil
}
