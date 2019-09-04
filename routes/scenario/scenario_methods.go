package scenario

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
	"github.com/jinzhu/gorm"
)

type Scenario struct {
	common.Scenario
}

func (s *Scenario) ByID(id uint) error {
	db := common.GetDB()
	err := db.Find(s, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Scenario) getUsers() ([]common.User, int, error) {
	db := common.GetDB()
	var users []common.User
	err := db.Order("ID asc").Model(s).Related(&users, "Users").Error
	return users, len(users), err
}

func (s *Scenario) save() error {
	db := common.GetDB()
	err := db.Create(s).Error
	return err
}

func (s *Scenario) update(updatedScenario Scenario) error {

	// TODO: if the field is empty member shouldn't be updated
	s.Name = updatedScenario.Name
	s.Running = updatedScenario.Running
	s.StartParameters = updatedScenario.StartParameters

	db := common.GetDB()
	err := db.Model(s).Update(updatedScenario).Error
	return err
}

func (s *Scenario) addUser(u *common.User) error {

	db := common.GetDB()
	err := db.Model(s).Association("Users").Append(u).Error
	return err
}

func (s *Scenario) deleteUser(username string) error {
	db := common.GetDB()

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
	db := common.GetDB()

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
		db := common.GetDB()
		u := common.User{}
		u.Username = ""
		err := db.Order("ID asc").Model(s).Where("ID = ?", userID).Related(&u, "Users").Error
		if err != nil {
			return false
		} else {
			if u.Username != "" {
				return true
			} else {
				return false
			}
		}
	}

}
