package signal

import (
	"fmt"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulationmodel"
)

type Signal struct {
	common.Signal
}

func (s *Signal) save() error {
	db := common.GetDB()
	err := db.Create(s).Error
	return err
}

func (s *Signal) byID(id uint) error {
	db := common.GetDB()
	err := db.Find(s, id).Error
	if err != nil {
		return fmt.Errorf("Signal with id=%v does not exist", id)
	}
	return nil
}

func (s *Signal) addToSimulationModel() error {
	db := common.GetDB()
	var m simulationmodel.SimulationModel
	err := m.ByID(s.SimulationModelID)
	if err != nil {
		return err
	}

	// save signal to DB
	err = s.save()
	if err != nil {
		return err
	}

	// associate signal with simulation model in correct direction
	if s.Direction == "in" {
		err = db.Model(&m).Association("InputMapping").Append(s).Error
		if err != nil {
			return err
		}

		// adapt length of mapping
		m.InputLength = db.Model(m).Where("Direction = ?", "in").Association("InputMapping").Count()
		err = m.Update(m)
	} else {
		err = db.Model(&m).Association("OutputMapping").Append(s).Error
		if err != nil {
			return err
		}

		// adapt length of mapping
		m.OutputLength = db.Model(m).Where("Direction = ?", "out").Association("OutputMapping").Count()
		err = m.Update(m)
	}
	return err
}

func (s *Signal) update(modifiedSignal Signal) error {
	db := common.GetDB()

	err := db.Model(s).Updates(map[string]interface{}{
		"Name":  modifiedSignal.Name,
		"Unit":  modifiedSignal.Unit,
		"Index": modifiedSignal.Index,
	}).Error

	return err

}

func (s *Signal) delete() error {

	db := common.GetDB()
	var m simulationmodel.SimulationModel
	err := m.ByID(s.SimulationModelID)
	if err != nil {
		return err
	}

	// remove association between Signal and SimulationModel
	// Signal itself is not deleted from DB, it remains as "dangling"
	if s.Direction == "in" {
		err = db.Model(&m).Association("InputMapping").Delete(s).Error
		if err != nil {
			return err
		}

		// Reduce length of mapping by 1
		m.InputLength--
		err = m.Update(m)

	} else {
		err = db.Model(&m).Association("OutputMapping").Delete(s).Error
		if err != nil {
			return err
		}

		// Reduce length of mapping by 1
		m.OutputLength--
		err = m.Update(m)
	}

	return err
}
