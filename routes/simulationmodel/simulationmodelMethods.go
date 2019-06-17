package simulationmodel

import (
	"fmt"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulator"
)

type SimulationModel struct {
	common.SimulationModel
}

func (m *SimulationModel) save() error {
	db := common.GetDB()
	err := db.Create(m).Error
	return err
}

func (m *SimulationModel) ByID(id uint) error {
	db := common.GetDB()
	err := db.Find(m, id).Error
	if err != nil {
		return fmt.Errorf("Simulation Model with id=%v does not exist", id)
	}
	return nil
}

func (m *SimulationModel) addToSimulation() error {
	db := common.GetDB()
	var sim simulation.Simulation
	err := sim.ByID(m.SimulationID)
	if err != nil {
		return err
	}

	// save simulation model to DB
	err = m.save()
	if err != nil {
		return err
	}

	// associate simulator with simulation model
	var simltr simulator.Simulator
	err = simltr.ByID(m.SimulatorID)
	err = db.Model(m).Association("Simulator").Append(&simltr).Error

	// associate simulation model with simulation
	err = db.Model(&sim).Association("SimulationModels").Append(m).Error

	return err
}

func (m *SimulationModel) Update(modifiedSimulationModel SimulationModel) error {
	db := common.GetDB()

	if m.SimulatorID != modifiedSimulationModel.SimulatorID {
		// update simulator
		var s simulator.Simulator
		err := s.ByID(modifiedSimulationModel.SimulatorID)
		if err != nil {
			return err
		}
		err = db.Model(m).Association("Simulator").Replace(s).Error

	}

	err := db.Model(m).Updates(map[string]interface{}{
		"Name":            modifiedSimulationModel.Name,
		"OutputLength":    modifiedSimulationModel.OutputLength,
		"InputLength":     modifiedSimulationModel.InputLength,
		"StartParameters": modifiedSimulationModel.StartParameters,
		"SimulatorID":     modifiedSimulationModel.SimulatorID,
	}).Error

	return err
}

func (m *SimulationModel) delete() error {

	db := common.GetDB()
	var sim simulation.Simulation
	err := sim.ByID(m.SimulationID)
	if err != nil {
		return err
	}

	// remove association between SimulationModel and Simulation
	// SimulationModel itself is not deleted from DB, it remains as "dangling"
	err = db.Model(&sim).Association("SimulationModels").Delete(m).Error

	return err
}
