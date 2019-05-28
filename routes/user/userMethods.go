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

func (u *User) save() error {
	db := common.GetDB()
	err := db.Create(u).Error
	return err
}

func (u *User) ByUsername(username string) error {
	db := common.GetDB()
	err := db.Find(u, "Username = ?", username).Error
	if err != nil {
		return fmt.Errorf("User with username=%v does not exist", username)
	}
	return nil
}

func (u *User) byID(id uint) error {
	db := common.GetDB()
	err := db.Find(u, id).Error
	if err != nil {
		return fmt.Errorf("User with id=%v does not exist", id)
	}
	return nil
}

func (u *User) setPassword(password string) error {
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

func (u *User) update(modifiedUser User) error {
	db := common.GetDB()
	err := db.Model(u).Update(modifiedUser).Error
	return err
}
