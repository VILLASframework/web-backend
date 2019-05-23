package model

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
)

type Model struct{
	common.Model
}

func FindAllModels(simID int) ([]common.Model, int, error) {
	db := common.GetDB()
	var models []common.Model
	sim, err := simulation.FindSimulation(simID)
	if err != nil {
		return models, 0, err
	}

	err = db.Model(sim).Related(&models, "Models").Error

	return models, len(models), err
}

func FindModel(modelID int) (common.Model, error){
	db := common.GetDB()
	var m common.Model
	err := db.First(&m, modelID).Error
	return m, err
}

func (m *Model) addToSimulation(simID int) error {
	db := common.GetDB()
	sim, err := simulation.FindSimulation(simID)
	if err != nil {
		return err
	}

	err = db.Model(&sim).Association("Models").Append(m).Error
	return err
}

func (m *Model) UpdateModel(modelID int) error {
	db := common.GetDB()
	model_to_update, err := FindModel(modelID)
	if err != nil {
		return err
	}
	// only Name and Start Params can be updated directly by the user
	err = db.Model(&model_to_update).Updates(map[string]interface{}{"Name": m.Name, "StartParameters": m.StartParameters}).Error
	return err
}


func (m *Model) UpdateSimulator(simulator *common.Simulator) error {
	db := common.GetDB()
	err := db.Model(m).Association("Simulator").Replace(simulator).Error
	return err
}

func (m *Model) UpdateSignals(signals []common.Signal, direction string) error {

	db := common.GetDB()
	var err error

	if direction == "in" {
		err = db.Model(m).Select("InputMapping").Update("InputMapping", signals).Error
	} else {
		err = db.Model(m).Select("OutputMapping").Update("OutputMapping", signals).Error
	}

	return err

}