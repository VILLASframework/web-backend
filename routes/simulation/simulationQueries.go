package simulation

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func FindAllSimulations() ([]common.Simulation, int, error) {
	db := common.GetDB()
	var simulations []common.Simulation
	err := db.Find(&simulations).Error
	return simulations, len(simulations), err
}

func FindUserSimulations(user *common.User) ([]common.Simulation, int, error) {
	db := common.GetDB()
	var simulations []common.Simulation
	err := db.Model(user).Related(&simulations, "Simulations").Error
	return simulations, len(simulations), err
}
