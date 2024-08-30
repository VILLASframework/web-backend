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
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type validNewUserGroup struct {
	Name             string                    `form:"name" validate:"required,min=3,max=100"`
	ScenarioMappings []validNewScenarioMapping `form:"scenarioMappings" validate:"dive"`
}

type validNewScenarioMapping struct {
	ScenarioID uint `form:"scenario_id" validate:"required"`
	Duplicate  bool `form:"duplicate" validate:"omitempty"`
}

type validUpdatedUserGroup struct {
	Name             string                        `form:"name" validate:"omitempty"`
	ScenarioMappings []validUpdatedScenarioMapping `form:"scenarioMappings" validate:"omitempty"`
}

type validUpdatedScenarioMapping struct {
	ScenarioID uint `form:"scenarioID" validate:"omitempty"`
	Duplicate  bool `form:"duplicate" validate:"omitempty"`
}

type addUserGroupRequest struct {
	UserGroup validNewUserGroup `json:"userGroup"`
}

type updateUserGroupRequest struct {
	UserGroup validUpdatedUserGroup `json:"userGroup"`
}

func (r *addUserGroupRequest) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *validUpdatedUserGroup) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *addUserGroupRequest) createUserGroup() UserGroup {
	var ug UserGroup
	ug.Name = r.UserGroup.Name
	ug.ScenarioMappings = convertScenarioMappings(r.UserGroup.ScenarioMappings)
	return ug
}

func convertScenarioMappings(validMappings []validNewScenarioMapping) []database.ScenarioMapping {
	scenarioMappings := make([]database.ScenarioMapping, len(validMappings))
	for i, v := range validMappings {
		scenarioMappings[i] = database.ScenarioMapping{
			ScenarioID: v.ScenarioID,
			Duplicate:  v.Duplicate,
		}
	}
	return scenarioMappings
}

func (r *updateUserGroupRequest) updatedUserGroup(oldUserGroup UserGroup) UserGroup {
	// Use the old UserGroup as a basis for the updated UserGroup `ug`
	ug := oldUserGroup

	if r.UserGroup.Name != "string" && r.UserGroup.Name != "" {
		ug.Name = r.UserGroup.Name
	}

	return ug
}
