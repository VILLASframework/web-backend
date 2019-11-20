/** User package, methods.
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
package widget

import (
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/dashboard"
)

type Widget struct {
	database.Widget
}

func (w *Widget) save() error {
	db := database.GetDB()
	err := db.Create(w).Error
	return err
}

func (w *Widget) ByID(id uint) error {
	db := database.GetDB()
	err := db.Find(w, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (w *Widget) addToDashboard() error {
	db := database.GetDB()
	var dab dashboard.Dashboard
	err := dab.ByID(uint(w.DashboardID))
	if err != nil {
		return err
	}

	// save widget to DB
	err = w.save()
	if err != nil {
		return err
	}

	// associate dashboard with simulation
	err = db.Model(&dab).Association("Widgets").Append(w).Error

	return err
}

func (w *Widget) update(modifiedWidget Widget) error {

	db := database.GetDB()
	err := db.Model(w).Updates(map[string]interface{}{
		"Name":             modifiedWidget.Name,
		"Type":             modifiedWidget.Type,
		"Width":            modifiedWidget.Width,
		"Height":           modifiedWidget.Height,
		"MinWidth":         modifiedWidget.MinWidth,
		"MinHeight":        modifiedWidget.MinHeight,
		"X":                modifiedWidget.X,
		"Y":                modifiedWidget.Y,
		"Z":                modifiedWidget.Z,
		"IsLocked":         modifiedWidget.IsLocked,
		"CustomProperties": modifiedWidget.CustomProperties,
	}).Error

	return err
}

func (w *Widget) delete() error {

	db := database.GetDB()
	var dab dashboard.Dashboard
	err := dab.ByID(w.DashboardID)
	if err != nil {
		return err
	}

	// remove association between Dashboard and Widget
	// Widget itself is not deleted from DB, it remains as "dangling"
	err = db.Model(&dab).Association("Widgets").Delete(w).Error

	// TODO: What about files that belong to a widget? Keep them or remove them here?

	return err
}
