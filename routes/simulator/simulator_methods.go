package simulator

import (
	"fmt"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/database"
)

type Simulator struct {
	database.Simulator
}

func (s *Simulator) save() error {
	db := database.GetDB()
	err := db.Create(s).Error
	return err
}

func (s *Simulator) ByID(id uint) error {
	db := database.GetDB()
	err := db.Find(s, id).Error
	return err
}

func (s *Simulator) update(updatedSimulator Simulator) error {

	db := database.GetDB()
	err := db.Model(s).Updates(updatedSimulator).Error

	return err
}

func (s *Simulator) delete() error {
	db := database.GetDB()

	no_simulationmodels := db.Model(s).Association("SimulationModels").Count()

	if no_simulationmodels > 0 {
		return fmt.Errorf("Simulator cannot be deleted as it is still used in SimulationModels (active or dangling)")
	}

	// delete Simulator from DB (does NOT remain as dangling)
	err := db.Delete(s).Error
	return err
}

func (s *Simulator) getModels() ([]database.SimulationModel, int, error) {
	db := database.GetDB()
	var models []database.SimulationModel
	err := db.Order("ID asc").Model(s).Related(&models, "SimulationModels").Error
	return models, len(models), err
}
