package simulationmodel

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func FindAllSimulationModels() ([]common.SimulationModel, int, error) {
	db := common.GetDB()
	var simulationmodels []common.SimulationModel
	err := db.Find(&simulationmodels).Error
	return simulationmodels, len(simulationmodels), err
}
