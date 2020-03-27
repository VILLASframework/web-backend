/** component_configuration package, methods.
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
package component_configuration

import (
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/infrastructure-component"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
)

type ComponentConfiguration struct {
	database.ComponentConfiguration
}

func (m *ComponentConfiguration) save() error {
	db := database.GetDB()
	err := db.Create(m).Error
	return err
}

func (m *ComponentConfiguration) ByID(id uint) error {
	db := database.GetDB()
	err := db.Find(m, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (m *ComponentConfiguration) addToScenario() error {
	db := database.GetDB()
	var so scenario.Scenario
	err := so.ByID(m.ScenarioID)
	if err != nil {
		return err
	}

	// save component configuration to DB
	err = m.save()
	if err != nil {
		return err
	}

	// associate IC with component configuration
	var ic infrastructure_component.InfrastructureComponent
	err = ic.ByID(m.ICID)
	err = db.Model(&ic).Association("ComponentConfigurations").Append(m).Error
	if err != nil {
		return err
	}

	// associate component configuration with scenario
	err = db.Model(&so).Association("ComponentConfigurations").Append(m).Error

	return err
}

func (m *ComponentConfiguration) Update(modifiedConfig ComponentConfiguration) error {
	db := database.GetDB()

	// check if IC has been updated
	if m.ICID != modifiedConfig.ICID {
		// update IC
		var s infrastructure_component.InfrastructureComponent
		var s_old infrastructure_component.InfrastructureComponent
		err := s.ByID(modifiedConfig.ICID)
		if err != nil {
			return err
		}
		err = s_old.ByID(m.ICID)
		if err != nil {
			return err
		}
		// remove component configuration from old IC
		err = db.Model(&s_old).Association("ComponentConfigurations").Delete(m).Error
		if err != nil {
			return err
		}
		// add component configuration to new IC
		err = db.Model(&s).Association("ComponentConfigurations").Append(m).Error
		if err != nil {
			return err
		}
	}

	err := db.Model(m).Updates(map[string]interface{}{
		"Name":            modifiedConfig.Name,
		"StartParameters": modifiedConfig.StartParameters,
		"ICID":            modifiedConfig.ICID,
		"SelectedFileID":  modifiedConfig.SelectedFileID,
	}).Error

	return err
}

func (m *ComponentConfiguration) delete() error {

	db := database.GetDB()
	var so scenario.Scenario
	err := so.ByID(m.ScenarioID)
	if err != nil {
		return err
	}

	// remove association between ComponentConfiguration and Scenario
	// ComponentConfiguration itself is not deleted from DB, it remains as "dangling"
	err = db.Model(&so).Association("ComponentConfigurations").Delete(m).Error

	return err
}