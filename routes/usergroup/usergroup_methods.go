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

package usergroup

import (
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
)

type UserGroup struct {
	database.UserGroup
}

type ScenarioMapping struct {
	database.ScenarioMapping
}

func (ug *UserGroup) save() error {
	db := database.GetDB()
	err := db.Create(ug).Error
	return err
}

func (ug *UserGroup) update(updatedUserGroup UserGroup, reqScenarioMappings []validUpdatedScenarioMapping) error {
	ug.Name = updatedUserGroup.Name

	db := database.GetDB()
	err := db.Model(ug).Update(updatedUserGroup).Error
	if err != nil {
		return err
	}
	/*
		err = db.Model(ug).Updates(database.UserGroup{Name: ug.Name}).Error
		if err != nil {
			return err
		}
	*/

	return updateScenarioMappings(ug.ID, reqScenarioMappings)
}

func updateScenarioMappings(groupID uint, reqScenarioMappings []validUpdatedScenarioMapping) error {
	var oldMappings []database.ScenarioMapping
	db := database.GetDB()
	err := db.Where("user_group_id = ?", groupID).Find(&oldMappings).Error
	if err != nil {
		return err
	}

	oldMappingsMap := make(map[uint]database.ScenarioMapping)
	for _, mapping := range oldMappings {
		oldMappingsMap[mapping.ScenarioID] = mapping
	}

	// Handle ScenarioMappings (add/update/delete)
	for _, reqMapping := range reqScenarioMappings {
		if oldMapping, exists := oldMappingsMap[reqMapping.ScenarioID]; exists {
			// Update
			oldMapping.Duplicate = reqMapping.Duplicate
			err = db.Save(&oldMapping).Error
			if err != nil {
				return err
			}
			delete(oldMappingsMap, reqMapping.ScenarioID)
		} else {
			// Add
			newMapping := database.ScenarioMapping{
				ScenarioID:  reqMapping.ScenarioID,
				UserGroupID: groupID,
				Duplicate:   reqMapping.Duplicate,
			}
			err = db.Create(&newMapping).Error
			if err != nil {
				return err
			}
		}
	}

	// Delete old mappings that were not in the request
	for _, mapping := range oldMappingsMap {
		err = db.Delete(&mapping).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (ug *UserGroup) byID(id uint) error {
	db := database.GetDB()
	err := db.Find(ug, id).Error
	return err
}

func (u *UserGroup) remove() error {
	db := database.GetDB()
	err := db.Delete(u).Error
	return err
}
