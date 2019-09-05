package simulationmodel

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/scenario"
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
		return err
	}
	return nil
}

func (m *SimulationModel) addToScenario() error {
	db := common.GetDB()
	var so scenario.Scenario
	err := so.ByID(m.ScenarioID)
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
	err = db.Model(&simltr).Association("SimulationModels").Append(m).Error
	if err != nil {
		return err
	}

	// associate simulation model with scenario
	err = db.Model(&so).Association("SimulationModels").Append(m).Error

	return err
}

func (m *SimulationModel) Update(modifiedSimulationModel SimulationModel) error {
	db := common.GetDB()

	if m.SimulatorID != modifiedSimulationModel.SimulatorID {
		// update simulator
		var s simulator.Simulator
		var s_old simulator.Simulator
		err := s.ByID(modifiedSimulationModel.SimulatorID)
		if err != nil {
			return err
		}
		err = s_old.ByID(m.SimulatorID)
		if err != nil {
			return err
		}
		// remove simulation model from old simulator
		err = db.Model(&s_old).Association("SimulationModels").Delete(m).Error
		if err != nil {
			return err
		}
		// add simulation model to new simulator
		err = db.Model(&s).Association("SimulationModels").Append(m).Error
		if err != nil {
			return err
		}
	}

	if m.ScenarioID != modifiedSimulationModel.ScenarioID {
		// TODO do we allow this case?
	}

	err := db.Model(m).Updates(map[string]interface{}{
		"Name":            modifiedSimulationModel.Name,
		"StartParameters": modifiedSimulationModel.StartParameters,
		"SimulatorID":     modifiedSimulationModel.SimulatorID,
		//"ScenarioID":     modifiedSimulationModel.ScenarioID,
	}).Error

	return err
}

func (m *SimulationModel) delete() error {

	db := common.GetDB()
	var so scenario.Scenario
	err := so.ByID(m.ScenarioID)
	if err != nil {
		return err
	}

	// remove association between SimulationModel and Scenario
	// SimulationModel itself is not deleted from DB, it remains as "dangling"
	err = db.Model(&so).Association("SimulationModels").Delete(m).Error

	return err
}
