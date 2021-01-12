/** Dashboard package, methods.
*
* @author Sonja Happ <sonja.happ@eonerc.rwth-aachen.de>
* @copyright 2014-2019, Institute for Automation of Complex Power Systems, EONERC
* @license GNU General Public License (version 3)
*
* VILLASweb-backend-go
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/
package dashboard

import (
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
)

type Dashboard struct {
	database.Dashboard
}

func (d *Dashboard) save() error {
	db := database.GetDB()
	err := db.Create(d).Error
	return err
}

func (d *Dashboard) ByID(id uint) error {
	db := database.GetDB()
	err := db.Find(d, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *Dashboard) addToScenario() error {
	db := database.GetDB()
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

func (d *Dashboard) update(modifiedDab Dashboard) error {

	db := database.GetDB()

	err := db.Model(d).Updates(map[string]interface{}{
		"Name":   modifiedDab.Name,
		"Grid":   modifiedDab.Grid,
		"Height": modifiedDab.Height,
	}).Error

	return err
}

func (d *Dashboard) delete() error {

	db := database.GetDB()
	var sim scenario.Scenario
	err := sim.ByID(d.ScenarioID)
	if err != nil {
		return err
	}

	// remove association between Dashboard and Scenario
	err = db.Model(&sim).Association("Dashboards").Delete(d).Error

	// get all widgets of the dashboard
	var widgets []database.Widget
	err = db.Order("ID asc").Model(d).Related(&widgets, "Widgets").Error
	if err != nil {
		return err
	}

	// Delete widgets
	for widget, _ := range widgets {
		err = db.Delete(&widget).Error
	}

	// Delete dashboard
	err = db.Delete(d).Error

	return err
}
