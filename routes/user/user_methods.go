/** User package, methods.
*
* @author Sonja Happ <sonja.happ@eonerc.rwth-aachen.de>
* @copyright 2014-2019, Institute for Automation of Complex Power Systems, EONERC
* @license GNU General Public License (version 3)
*
* VILLASweb-backend-go
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/
package user

import (
	"fmt"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10

// This is ugly but no other way to keep each model on the corresponding
// package since we have circular dependencies. Methods of a type must
// live in the same package.
type User struct {
	database.User // check golang embedding types
}

func NewUser(username, password, mail, role string, active bool) (User, error) {
	var newUser User

	// Check that the username is NOT taken
	err := newUser.ByUsername(username)
	if err == nil {
		return newUser, &UsernameAlreadyTaken{Username: username}
	}

	newUser.Username = username
	newUser.Mail = mail
	newUser.Role = role
	newUser.Active = active

	if password == "" {
		// This user is authenticated via some external system
		newUser.Password = ""
	} else {
		// Hash the password before saving it to the DB
		err = newUser.setPassword(password)
		if err != nil {
			return newUser, err
		}
	}

	// Save the user in the DB
	err = newUser.save()
	if err != nil {
		return newUser, err
	}

	return newUser, nil
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
		return fmt.Errorf("password cannot be empty")
	}
	newPassword, err :=
		bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return fmt.Errorf("failed to generate hash from password")
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

	u.Username = updatedUser.Username
	u.Password = updatedUser.Password
	u.Mail = updatedUser.Mail
	u.Role = updatedUser.Role
	u.Active = updatedUser.Active

	db := database.GetDB()
	err := db.Model(u).Update(updatedUser).Error
	if err != nil {
		return err
	}

	// extra update for bool Active since it is ignored if false
	err = db.Model(u).Updates(map[string]interface{}{"Active": updatedUser.Active}).Error

	return err
}
