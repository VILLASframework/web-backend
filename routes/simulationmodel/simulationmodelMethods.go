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

func (m *SimulationModel) update(modifiedSimulationModel SimulationModel) error {
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

	err := db.Model(m).Updates(map[string]interface{}{"Name": modifiedSimulationModel.Name,
		"OutputLength":    modifiedSimulationModel.OutputLength,
		"InputLength":     modifiedSimulationModel.InputLength,
		"StartParameters": modifiedSimulationModel.StartParameters,
		"SimulatorID":     modifiedSimulationModel.SimulatorID,
	}).Error
	if err != nil {
		return err
	}

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

func (m *SimulationModel) addSignal(signal common.Signal) error {

	db := common.GetDB()
	var err error

	if signal.Direction == "in" {
		err = db.Model(m).Association("InputMapping").Append(signal).Error
		if err != nil {
			return err
		}
		// adapt length of mapping
		m.InputLength = db.Model(m).Where("Direction = ?", "in").Association("InputMapping").Count()
		err = m.update(*m)

	} else {
		err = db.Model(m).Association("OutputMapping").Append(signal).Error
		if err != nil {
			return err
		}

		// adapt length of mapping
		m.OutputLength = db.Model(m).Where("Direction = ?", "out").Association("OutputMapping").Count()
		err = m.update(*m)

	}

	return err
}

func (m *SimulationModel) deleteSignals(direction string) error {

	db := common.GetDB()
	var err error

	var columnName string

	if direction == "in" {
		columnName = "InputMapping"

	} else {
		columnName = "OutputMapping"
	}

	var signals []common.Signal
	err = db.Order("ID asc").Model(m).Where("Direction = ?", direction).Related(&signals, columnName).Error
	if err != nil {
		return err
	}

	// remove association to each signal and delete each signal from db
	for _, sig := range signals {
		err = db.Model(m).Association(columnName).Delete(sig).Error
		if err != nil {
			return err
		}
		err = db.Delete(sig).Error
		if err != nil {
			return err
		}
	}

	// set length of mapping to 0
	if columnName == "InputMapping" {
		m.InputLength = 0
	} else {
		m.OutputLength = 0
	}
	err = m.update(*m)

	return err
}
