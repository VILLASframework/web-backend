package simulator

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func FindAllSimulators() ([]common.Simulator, int, error) {
	db := common.GetDB()
	var simulators []common.Simulator
	err := db.Find(&simulators).Error
	return simulators, len(simulators), err
}

func FindSimulator(simulatorID int) (common.Simulator, error) {
	db := common.GetDB()
	var simulator common.Simulator
	err := db.First(&simulator, simulatorID).Error
	return simulator, err
}

