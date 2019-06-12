package user

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10

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

func (u *User) remove() error {
	db := common.GetDB()
	err := db.Delete(u).Error
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

func (u *User) ByID(id uint) error {
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

func (u *User) update(updatedUser User) error {

	// TODO: if the field is empty string member shouldn't be updated
	u.Username = updatedUser.Username
	u.Password = updatedUser.Password
	u.Mail = updatedUser.Mail
	u.Role = updatedUser.Role

	db := common.GetDB()
	err := db.Model(u).Update(updatedUser).Error
	return err
}
