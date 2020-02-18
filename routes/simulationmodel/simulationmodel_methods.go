/** Simulationmodel package, methods.
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
package simulationmodel

import (
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/simulator"
)

type SimulationModel struct {
	database.SimulationModel
}

func (m *SimulationModel) save() error {
	db := database.GetDB()
	err := db.Create(m).Error
	return err
}

func (m *SimulationModel) ByID(id uint) error {
	db := database.GetDB()
	err := db.Find(m, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (m *SimulationModel) addToScenario() error {
	db := database.GetDB()
	var so scenario.Scenario
	err := so.ByID(m.ScenarioID)
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
	err = db.Model(&simltr).Association("SimulationModels").Append(m).Error
	if err != nil {
		return err
	}

	// associate simulation model with scenario
	err = db.Model(&so).Association("SimulationModels").Append(m).Error

	return err
}

func (m *SimulationModel) Update(modifiedSimulationModel SimulationModel) error {
	db := database.GetDB()

	// check if simulator has been updated
	if m.SimulatorID != modifiedSimulationModel.SimulatorID {
		// update simulator
		var s simulator.Simulator
		var s_old simulator.Simulator
		err := s.ByID(modifiedSimulationModel.SimulatorID)
		if err != nil {
			return err
		}
		err = s_old.ByID(m.SimulatorID)
		if err != nil {
			return err
		}
		// remove simulation model from old simulator
		err = db.Model(&s_old).Association("SimulationModels").Delete(m).Error
		if err != nil {
			return err
		}
		// add simulation model to new simulator
		err = db.Model(&s).Association("SimulationModels").Append(m).Error
		if err != nil {
			return err
		}
	}

	err := db.Model(m).Updates(map[string]interface{}{
		"Name":                modifiedSimulationModel.Name,
		"StartParameters":     modifiedSimulationModel.StartParameters,
		"SimulatorID":         modifiedSimulationModel.SimulatorID,
		"SelectedModelFileID": modifiedSimulationModel.SelectedModelFileID,
	}).Error

	return err
}

func (m *SimulationModel) delete() error {

	db := database.GetDB()
	var so scenario.Scenario
	err := so.ByID(m.ScenarioID)
	if err != nil {
		return err
	}

	// remove association between SimulationModel and Scenario
	// SimulationModel itself is not deleted from DB, it remains as "dangling"
	err = db.Model(&so).Association("SimulationModels").Delete(m).Error

	return err
}
