package simulator

import (
	"fmt"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

type Simulator struct {
	common.Simulator
}

func (s *Simulator) save() error {
	db := common.GetDB()
	err := db.Create(s).Error
	return err
}

func (s *Simulator) ByID(id uint) error {
	db := common.GetDB()
	err := db.Find(s, id).Error
	if err != nil {
		return fmt.Errorf("Simulator with id=%v does not exist", id)
	}
	return nil
}

func (s *Simulator) update(modifiedSimulator common.SimulatorResponse) error {

	db := common.GetDB()
	err := db.Model(s).Updates(modifiedSimulator).Error

	return err

}

func (s *Simulator) delete() error {
	db := common.GetDB()

	no_simulationmodels := db.Model(s).Association("SimulationModels").Count()

	if no_simulationmodels > 0 {
		return fmt.Errorf("Simulator cannot be deleted as it is still used in SimulationModels (active or dangling)")
	}

	// delete Simulator from DB (does NOT remain as dangling)
	err := db.Delete(s).Error
	return err
}

func (s *Simulator) getModels() ([]common.SimulationModel, int, error) {
	db := common.GetDB()
	var models []common.SimulationModel
	err := db.Order("ID asc").Model(s).Related(&models, "SimulationModels").Error
	return models, len(models), err
}
