/**
* This file is part of VILLASweb-backend-go
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
)

type Signal struct {
	database.Signal
}

func (s *Signal) save() error {
	db := database.GetDB()
	err := db.Create(s).Error
	return err
}

/*func (s *Signal) byID(id uint) error {
	db := database.GetDB()
	err := db.Find(s, id).Error
	if err != nil {
		return err
	}
	return nil
}*/

func (s *Signal) AddToConfig() error {
	db := database.GetDB()
	var err error
	var m database.ComponentConfiguration
	err = db.Find(&m, s.ConfigID).Error
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
	} else {
		err = db.Model(&m).Association("OutputMapping").Append(s).Error
	}

	if err != nil {
		return err
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
	var err error
	var m database.ComponentConfiguration
	err = db.Find(&m, s.ConfigID).Error
	if err != nil {
		return err
	}

	// remove association between Signal and ComponentConfiguration
	if s.Direction == "in" {
		err = db.Model(&m).Association("InputMapping").Delete(s).Error
	} else {
		err = db.Model(&m).Association("OutputMapping").Delete(s).Error
	}

	if err != nil {
		return err
	}

	// Delete signal
	err = db.Delete(s).Error

	return err
}
