/** Signal package, methods.
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
package signal

import (
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/component-configuration"
)

type Signal struct {
	database.Signal
}

func (s *Signal) save() error {
	db := database.GetDB()
	err := db.Create(s).Error
	return err
}

func (s *Signal) byID(id uint) error {
	db := database.GetDB()
	err := db.Find(s, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Signal) addToConfig() error {
	db := database.GetDB()
	var m component_configuration.ComponentConfiguration
	err := m.ByID(s.ConfigID)
	if err != nil {
		return err
	}

	// save signal to DB
	err = s.save()
	if err != nil {
		return err
	}

	// associate signal with component configuration in correct direction
	if s.Direction == "in" {
		err = db.Model(&m).Association("InputMapping").Append(s).Error
		if err != nil {
			return err
		}

		// adapt length of mapping
		var newInputLength = db.Model(m).Where("Direction = ?", "in").Association("InputMapping").Count()
		err = db.Model(m).Update("InputLength", newInputLength).Error

	} else {
		err = db.Model(&m).Association("OutputMapping").Append(s).Error
		if err != nil {
			return err
		}

		// adapt length of mapping
		var newOutputLength = db.Model(m).Where("Direction = ?", "out").Association("OutputMapping").Count()
		err = db.Model(m).Update("OutputLength", newOutputLength).Error
	}
	return err
}

func (s *Signal) update(modifiedSignal Signal) error {
	db := database.GetDB()

	err := db.Model(s).Updates(map[string]interface{}{
		"Name":          modifiedSignal.Name,
		"Unit":          modifiedSignal.Unit,
		"Index":         modifiedSignal.Index,
		"ScalingFactor": modifiedSignal.ScalingFactor,
	}).Error

	return err

}

func (s *Signal) delete() error {

	db := database.GetDB()
	var m component_configuration.ComponentConfiguration
	err := m.ByID(s.ConfigID)
	if err != nil {
		return err
	}

	// remove association between Signal and ComponentConfiguration
	// Signal itself is not deleted from DB, it remains as "dangling"
	if s.Direction == "in" {
		err = db.Model(&m).Association("InputMapping").Delete(s).Error
		if err != nil {
			return err
		}

		// Reduce length of mapping by 1
		var newInputLength = m.InputLength - 1
		err = db.Model(m).Update("InputLength", newInputLength).Error

	} else {
		err = db.Model(&m).Association("OutputMapping").Delete(s).Error
		if err != nil {
			return err
		}

		// Reduce length of mapping by 1
		var newOutputLength = m.OutputLength - 1
		err = db.Model(m).Update("OutputLength", newOutputLength).Error
	}

	return err
}
