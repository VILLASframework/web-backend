package user

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/database"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10

// This is ugly but no other way to keep each model on the corresponding
// package since we have circular dependencies. Methods of a type must
// live in the same package.
type User struct {
	database.User // check golang embedding types
}

func (u *User) save() error {
	db := database.GetDB()
	err := db.Create(u).Error
	return err
}

func (u *User) remove() error {
	db := database.GetDB()
	err := db.Delete(u).Error
	return err
}

func (u *User) ByUsername(username string) error {
	db := database.GetDB()
	err := db.Find(u, "Username = ?", username).Error
	return err
}

func (u *User) ByID(id uint) error {
	db := database.GetDB()
	err := db.Find(u, id).Error
	return err
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

	db := database.GetDB()
	err := db.Model(u).Update(updatedUser).Error
	return err
}
