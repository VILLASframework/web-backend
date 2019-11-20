/** Simulator package, methods.
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
package simulator

import (
	"fmt"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
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
