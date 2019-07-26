package dashboard

import (
	"fmt"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/scenario"
)

type Dashboard struct {
	common.Dashboard
}

func (v *Dashboard) save() error {
	db := common.GetDB()
	err := db.Create(v).Error
	return err
}

func (d *Dashboard) ByID(id uint) error {
	db := common.GetDB()
	err := db.Find(d, id).Error
	if err != nil {
		return fmt.Errorf("Dashboard with id=%v does not exist", id)
	}
	return nil
}

func (d *Dashboard) addToScenario() error {
	db := common.GetDB()
	var sim scenario.Scenario
	err := sim.ByID(d.ScenarioID)
	if err != nil {
		return err
	}

	// save dashboard to DB
	err = d.save()
	if err != nil {
		return err
	}

	// associate dashboard with scenario
	err = db.Model(&sim).Association("Dashboards").Append(d).Error

	return err
}

func (d *Dashboard) update(modifiedDab common.DashboardResponse) error {

	db := common.GetDB()

	err := db.Model(d).Updates(map[string]interface{}{
		"Name": modifiedDab.Name,
		"Grid": modifiedDab.Grid,
	}).Error

	return err
}

func (d *Dashboard) delete() error {

	db := common.GetDB()
	var sim scenario.Scenario
	err := sim.ByID(d.ScenarioID)
	if err != nil {
		return err
	}

	// remove association between Dashboard and Scenario
	// Dashboard itself is not deleted from DB, it remains as "dangling"
	err = db.Model(&sim).Association("Dashboards").Delete(d).Error

	return err
}
