package user

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func FindAllUsers() ([]common.User, int, error) {
	db := common.GetDB()
	var users []common.User
	err := db.Order("ID asc").Find(&users).Error
	return users, len(users), err
}

func FindAllUsersSim(sim *common.Simulation) ([]common.User, int, error) {
	db := common.GetDB()
	var users []common.User
	err := db.Order("ID asc").Model(sim).Related(&users, "Users").Error
	return users, len(users), err
}

func FindUserByName(username string) (common.User, error){
	db := common.GetDB()
	var user common.User
	err := db.Where("Username = ?", username).Find(&user).Error
	return user, err
}

func AddUserToSim(sim *common.Simulation, user *common.User) error {
	db := common.GetDB()
	err := db.Model(sim).Association("Users").Append(user).Error
	return err
}

func RemoveUserFromSim(sim *common.Simulation, username string) error {
	db := common.GetDB()

	user, err := FindUserByName(username)
	if err != nil {
		return err
	}

	// remove user from simulation
	err = db.Model(sim).Association("Users").Delete(&user).Error
	if err != nil {
		return err
	}

	// remove simulation from user
	err = db.Model(&user).Association("Simulations").Delete(sim).Error

	return err
}