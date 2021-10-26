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
	"log"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
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
	var so database.Scenario
	err := db.Find(&so, m.ScenarioID).Error
	if err != nil {
		return err
	}

	// save component configuration to DB
	err = m.save()
	if err != nil {
		return err
	}

	// associate IC with component configuration
	var ic database.InfrastructureComponent
	err = db.Find(&ic, m.ICID).Error
	if err != nil {
		return err
	}
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
		var ic database.InfrastructureComponent
		var ic_old database.InfrastructureComponent
		err := db.Find(&ic, modifiedConfig.ICID).Error
		if err != nil {
			return err
		}
		err = db.Find(&ic_old, m.ICID).Error
		if err != nil {
			return err
		}

		// remove component configuration from old IC
		err = db.Model(&ic_old).Association("ComponentConfigurations").Delete(m).Error
		if err != nil {
			return err
		}
		// add component configuration to new IC
		err = db.Model(&ic).Association("ComponentConfigurations").Append(m).Error
		if err != nil {
			return err
		}
	}

	err := db.Model(m).Updates(map[string]interface{}{
		"Name":            modifiedConfig.Name,
		"StartParameters": modifiedConfig.StartParameters,
		"ICID":            modifiedConfig.ICID,
		"FileIDs":         modifiedConfig.FileIDs,
	}).Error

	return err
}

func (m *ComponentConfiguration) delete() error {

	db := database.GetDB()
	var so database.Scenario
	err := db.Find(&so, m.ScenarioID).Error
	if err != nil {
		return err
	}

	var ic database.InfrastructureComponent
	err = db.Find(&ic, m.ICID).Error
	if err != nil {
		return err
	}

	// remove association between ComponentConfiguration and Scenario
	err = db.Model(&so).Association("ComponentConfigurations").Delete(m).Error
	if err != nil {
		return err
	}

	// remove association between Infrastructure component and config
	err = db.Model(&ic).Association("ComponentConfigurations").Delete(m).Error
	if err != nil {
		return err
	}

	// Get Signals of InputMapping and delete them
	var InputMappingSignals []database.Signal
	err = db.Model(m).Related(&InputMappingSignals, "InputMapping").Error
	if err != nil {
		return err
	}
	for sig := range InputMappingSignals {
		err = db.Delete(&sig).Error
		if err != nil {
			return err
		}
	}

	// Get Signals of OutputMapping and delete them
	var OutputMappingSignals []database.Signal
	err = db.Model(m).Related(&OutputMappingSignals, "OutputMapping").Error
	if err != nil {
		return err
	}
	for sig := range OutputMappingSignals {
		err = db.Delete(&sig).Error
		if err != nil {
			return err
		}
	}

	// delete component configuration
	err = db.Delete(m).Error
	if err != nil {
		return err
	}

	// if IC has state gone and there is no component configuration associated with it: delete IC
	no_configs := db.Model(ic).Association("ComponentConfigurations").Count()
	if no_configs == 0 && ic.State == "gone" {
		log.Println("Deleting IC with state gone, last component config deleted", ic.UUID)
		err = db.Delete(ic).Error
		return err
	}

	return nil
}
