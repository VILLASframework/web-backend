package simulation

import (
	"fmt"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
)

type Simulation struct {
	common.Simulation
}

func (s *Simulation) ByID(id uint) error {
	db := common.GetDB()
	err := db.Find(s, id).Error
	if err != nil {
		return fmt.Errorf("simulation with id=%v does not exist", id)
	}
	return nil
}

func (s *Simulation) getUsers() ([]common.User, int, error) {
	db := common.GetDB()
	var users []common.User
	err := db.Order("ID asc").Model(s).Related(&users, "Users").Error
	return users, len(users), err
}

func (s *Simulation) save() error {
	db := common.GetDB()
	err := db.Create(s).Error
	return err
}

func (s *Simulation) update(modifiedSimulation Simulation) error {
	db := common.GetDB()
	err := db.Model(s).Update(modifiedSimulation).Error
	return err
}

func (s *Simulation) addUser(u *common.User) error {

	db := common.GetDB()
	err := db.Model(s).Association("Users").Append(u).Error
	return err
}

func (s *Simulation) deleteUser(username string) error {
	db := common.GetDB()

	var deletedUser user.User
	err := deletedUser.ByUsername(username)
	if err != nil {
		return err
	}

	no_users := db.Model(s).Association("Users").Count()

	if no_users > 1 {
		// remove user from simulation
		err = db.Model(s).Association("Users").Delete(&deletedUser.User).Error
		if err != nil {
			return err
		}
		// remove simulation from user
		err = db.Model(&deletedUser.User).Association("Simulations").Delete(s).Error
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("cannot delete last user from simulation without deleting simulation itself, doing nothing")
	}

	return nil
}

func (s *Simulation) delete() error {
	db := common.GetDB()
	no_models := db.Model(s).Association("SimulationModels").Count()
	no_visualizations := db.Model(s).Association("Visualizations").Count()

	if no_models > 0 || no_visualizations > 0 {
		return fmt.Errorf("cannot delete simulation that contains models and/ or visualizations, doing nothing")
	} else {
		// delete simulation from all users and vice versa

		users, no_users, err := s.getUsers()
		if err != nil {
			return err
		}

		if no_users > 0 {
			for _, u := range users {
				fmt.Println("User in delete: ", u)
				// remove user from simulation
				err = db.Model(s).Association("Users").Delete(&u).Error
				if err != nil {
					return err
				}
				// remove simulation from user
				err = db.Model(&u).Association("Simulations").Delete(s).Error
				if err != nil {
					return err
				}
			}
		}

		// Delete simulation
		err = db.Delete(s).Error
		if err != nil {
			return err
		}

	}

	return nil
}

func (s *Simulation) checkAccess(userID uint, userRole string) bool {

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
