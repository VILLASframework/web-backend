package visualization

import (
	"fmt"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
)

type Visualization struct {
	common.Visualization
}

func (v *Visualization) save() error {
	db := common.GetDB()
	err := db.Create(v).Error
	return err
}

func (v *Visualization) ByID(id uint) error {
	db := common.GetDB()
	err := db.Find(v, id).Error
	if err != nil {
		return fmt.Errorf("Visualization with id=%v does not exist", id)
	}
	return nil
}

func (v *Visualization) addToSimulation() error {
	db := common.GetDB()
	var sim simulation.Simulation
	err := sim.ByID(v.SimulationID)
	if err != nil {
		return err
	}

	// save visualization to DB
	err = v.save()
	if err != nil {
		return err
	}

	// associate visualization with simulation
	err = db.Model(&sim).Association("Visualizations").Append(v).Error

	return err
}

func (v *Visualization) update(modifiedVis Visualization) error {

	db := common.GetDB()

	err := db.Model(v).Updates(map[string]interface{}{
		"Name": modifiedVis.Name,
		"Grid": modifiedVis.Grid,
	}).Error

	return err
}

func (v *Visualization) delete() error {

	db := common.GetDB()
	var sim simulation.Simulation
	err := sim.ByID(v.SimulationID)
	if err != nil {
		return err
	}

	// remove association between Visualization and Simulation
	// Visualization itself is not deleted from DB, it remains as "dangling"
	err = db.Model(&sim).Association("Visualizations").Delete(v).Error

	return err
}
