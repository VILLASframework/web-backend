package queries

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func FindAllModels(simID int) ([]common.Model, int, error) {
	db := common.GetDB()
	var models []common.Model
	sim, err := FindSimulation(simID)
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

func AddModel(simID int, model *common.Model) error {
	db := common.GetDB()
	sim, err := FindSimulation(simID)
	if err != nil {
		return err
	}

	err = db.Model(&sim).Association("Models").Append(model).Error
	return err
}

func CloneModel(targetSim int, modelID int) error {

	// TODO TO BE IMPLEMENTED
	// Check if target sim exists
	// Check if model exists

	// Get all Signals of Model
	// Get Simulator of Model
	// Get Files of model

	// Add new model object to DB and associate with target sim
	// Add new signal objects to DB and associate with new model object (careful with directions)
	// Associate Simulator with new Model object
	var err error
	return err

}

func UpdateModel(modelID int , modelUpdate *common.Model) error {
	db := common.GetDB()
	m, err := FindModel(modelID)
	if err != nil {
		return err
	}
	// only Name and Start Params can be updated directly by the user
	err = db.Model(&m).Updates(map[string]interface{}{"Name": modelUpdate.Name, "StartParameters": modelUpdate.StartParameters}).Error
	return err
}

func DeleteModel(simID int , modelID int ) error {
	db := common.GetDB()
	sim, err := FindSimulation(simID)
	if err != nil {
		return err
	}

	m, err := FindModel(modelID)
	if err != nil {
		return err
	}

	//remove relationship between simulation and model
	err = db.Model(&sim).Association("Models").Delete(m).Error
	if err != nil {
		return err
	}

	// TODO remove File Associations and files on disk
	// TODO remove Signal Associations and Signals in DB
	// TODO how to treat simulator association?

	//remove model itself from DB
	//TODO: do we want this??
	err = db.Delete(&m).Error
	return err
}

func UpdateSimulatorOfModel(model *common.Model, simulator *common.Simulator) error {
	db := common.GetDB()
	err := db.Model(model).Association("Simulator").Replace(simulator).Error
	return err
}