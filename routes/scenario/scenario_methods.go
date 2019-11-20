/** Scenario package, methods.
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
package scenario

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"
	"github.com/jinzhu/gorm"
)

type Scenario struct {
	database.Scenario
}

func (s *Scenario) ByID(id uint) error {
	db := database.GetDB()
	err := db.Find(s, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Scenario) getUsers() ([]database.User, int, error) {
	db := database.GetDB()
	var users []database.User
	err := db.Order("ID asc").Model(s).Where("Active = ?", true).Related(&users, "Users").Error
	return users, len(users), err
}

func (s *Scenario) save() error {
	db := database.GetDB()
	err := db.Create(s).Error
	return err
}

func (s *Scenario) update(updatedScenario Scenario) error {

	// TODO: if the field is empty member shouldn't be updated
	s.Name = updatedScenario.Name
	s.Running = updatedScenario.Running
	s.StartParameters = updatedScenario.StartParameters

	db := database.GetDB()
	err := db.Model(s).Update(updatedScenario).Error
	return err
}

func (s *Scenario) addUser(u *database.User) error {

	db := database.GetDB()
	err := db.Model(s).Association("Users").Append(u).Error
	return err
}

func (s *Scenario) deleteUser(username string) error {
	db := database.GetDB()

	var deletedUser user.User
	err := deletedUser.ByUsername(username)
	if err != nil {
		return err
	}

	no_users := db.Model(s).Association("Users").Count()

	if no_users > 1 {
		// remove user from scenario
		err = db.Model(s).Association("Users").Delete(&deletedUser.User).Error
		if err != nil {
			return err
		}
		// remove scenario from user
		err = db.Model(&deletedUser.User).Association("Scenarios").Delete(s).Error
		if err != nil {
			return err
		}
	} else {
		// There is only one associated user
		var remainingUser user.User
		err = db.Model(s).Related(&remainingUser, "Users").Error
		if remainingUser.Username == username {
			// if the remaining user is the one to be deleted
			return fmt.Errorf("cannot delete last user from scenario without deleting scenario itself, doing nothing")
		} else {
			// the remaining user is NOT the one to be deleted
			// that means the user to be deleted is not associated with the scenario
			return gorm.ErrRecordNotFound
		}
	}

	return nil
}

func (s *Scenario) delete() error {
	db := database.GetDB()

	// delete scenario from all users and vice versa

	users, no_users, err := s.getUsers()
	if err != nil {
		return err
	}

	if no_users > 0 {
		for _, u := range users {
			// remove user from scenario
			err = db.Model(s).Association("Users").Delete(&u).Error
			if err != nil {
				return err
			}
			// remove scenario from user
			err = db.Model(&u).Association("Scenarios").Delete(s).Error
			if err != nil {
				return err
			}
		}
	}

	// Scenario is not deleted from DB, only associations with users are removed
	// Scenario remains "dangling" in DB

	// Delete scenario
	//err = db.Delete(s).Error
	//if err != nil {
	//	return err
	//}

	return nil
}

func (s *Scenario) checkAccess(userID uint, userRole string) bool {

	if userRole == "Admin" {
		return true
	} else {
		db := database.GetDB()
		u := database.User{}
		u.Username = ""
		err := db.Order("ID asc").Model(s).Where("ID = ?", userID).Related(&u, "Users").Error
		if err != nil || !u.Active {
			return false
		} else {
			return true
		}
	}

}
