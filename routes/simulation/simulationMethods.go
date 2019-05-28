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
		return fmt.Errorf("Simulation with id=%v does not exist", id)
	}
	return nil
}

func (s *Simulation) getUsers() ([]common.User, int, error) {
	db := common.GetDB()
	var users []common.User
	err := db.Order("ID asc").Model(s).Related(&users, "Users").Error
	return users, len(users), err
}

func (s *Simulation) update(modifiedSimulation Simulation) error {
	db := common.GetDB()
	err := db.Model(s).Update(modifiedSimulation).Error
	return err
}

func (s *Simulation) addUser(username string) error {

	var newUser user.User
	err := newUser.ByUsername(username)
	if err != nil {
		return err
	}

	db := common.GetDB()
	err = db.Model(s).Association("Users").Append(&newUser).Error
	return err
}

func (s *Simulation) deleteUser(username string) error {
	db := common.GetDB()

	var deletedUser user.User
	err := deletedUser.ByUsername(username)
	if err != nil {
		return err
	}

	// remove user from simulation
	err = db.Model(s).Association("Users").Delete(&deletedUser).Error
	if err != nil {
		return err
	}

	// remove simulation from user
	err = db.Model(&deletedUser).Association("Simulations").Delete(s).Error

	return err
}
