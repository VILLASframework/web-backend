package queries

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func FindAllSimulators() ([]common.Simulator, int, error) {
	db := common.GetDB()
	var simulators []common.Simulator
	err := db.Find(&simulators).Error
	return simulators, len(simulators), err
}

