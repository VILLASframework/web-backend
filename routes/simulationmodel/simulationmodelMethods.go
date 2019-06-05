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

func (m *SimulationModel) addToSimulation(simID uint) error {
	db := common.GetDB()
	var sim simulation.Simulation
	err := sim.ByID(simID)
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
	err := db.Model(m).Update(modifiedSimulationModel).Error
	if err != nil {
		return err
	}

	if m.SimulatorID != modifiedSimulationModel.SimulatorID {
		// update simulator
		var s simulator.Simulator
		err = s.ByID(modifiedSimulationModel.SimulatorID)

		err = db.Model(m).Association("Simulator").Replace(s).Error

	}

	return err
}

func (m *SimulationModel) updateSignals(signals []common.Signal, direction string) error {

	db := common.GetDB()
	var err error

	if direction == "in" {
		err = db.Model(m).Select("InputMapping").Update("InputMapping", signals).Error
	} else {
		err = db.Model(m).Select("OutputMapping").Update("OutputMapping", signals).Error
	}

	return err
}

func (m *SimulationModel) addSignal(signal common.Signal, direction string) error {

	db := common.GetDB()
	var err error

	if direction == "in" {
		err = db.Model(m).Association("InputMapping").Append(signal).Error
		if err != nil {
			return err
		}
		// adapt length of mapping
		m.InputLength = db.Model(m).Association("InputMapping").Count()

	} else {
		err = db.Model(m).Association("OutputMapping").Append(signal).Error

		// adapt length of mapping
		m.OutputLength = db.Model(m).Association("OutputMapping").Count()

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
	err = db.Order("ID asc").Model(m).Related(&signals, columnName).Error
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

	return err
}
