package simulation

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

type Simulation struct{
	common.Simulation
}

func FindAllSimulations() ([]common.Simulation, int, error) {
	db := common.GetDB()
	var simulations []common.Simulation
	err := db.Order("ID asc").Find(&simulations).Error
	return simulations, len(simulations), err
}

func FindUserSimulations(user *common.User) ([]common.Simulation, int, error) {
	db := common.GetDB()
	var simulations []common.Simulation
	err := db.Order("ID asc").Model(user).Related(&simulations, "Simulations").Error
	return simulations, len(simulations), err
}

func FindSimulation(simID int) (common.Simulation, error) {
	db := common.GetDB()
	var sim common.Simulation
	err := db.First(&sim, simID).Error
	return sim, err
}
