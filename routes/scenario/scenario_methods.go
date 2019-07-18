package scenario

import (
	"fmt"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
)

type Scenario struct {
	common.Scenario
}

func (s *Scenario) ByID(id uint) error {
	db := common.GetDB()
	err := db.Find(s, id).Error
	if err != nil {
		return fmt.Errorf("scenario with id=%v does not exist", id)
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

func (s *Scenario) update(modifiedScenario Scenario) error {
	db := common.GetDB()
	err := db.Model(s).Update(modifiedScenario).Error
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
		return fmt.Errorf("cannot delete last user from scenario without deleting scenario itself, doing nothing")
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
