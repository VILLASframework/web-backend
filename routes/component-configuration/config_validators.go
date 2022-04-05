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

package component_configuration

import (
	"encoding/json"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nsf/jsondiff"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type validNewConfig struct {
	Name            string         `form:"Name" validate:"required"`
	ScenarioID      uint           `form:"ScenarioID" validate:"required"`
	ICID            uint           `form:"ICID" validate:"omitempty"`
	StartParameters postgres.Jsonb `form:"StartParameters" validate:"required"`
	FileIDs         []int64        `form:"FileIDs" validate:"omitempty"`
}

type validUpdatedConfig struct {
	Name            string         `form:"Name" validate:"omitempty"`
	ICID            uint           `form:"ICID" validate:"omitempty"`
	StartParameters postgres.Jsonb `form:"StartParameters" validate:"omitempty"`
	FileIDs         []int64        `form:"FileIDs" validate:"omitempty"`
}

type addConfigRequest struct {
	Config validNewConfig `json:"config"`
}

type updateConfigRequest struct {
	Config validUpdatedConfig `json:"config"`
}

func (r *addConfigRequest) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *validUpdatedConfig) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *addConfigRequest) createConfig() ComponentConfiguration {
	var s ComponentConfiguration

	s.Name = r.Config.Name
	s.ScenarioID = r.Config.ScenarioID
	s.ICID = r.Config.ICID
	s.StartParameters = r.Config.StartParameters
	s.FileIDs = r.Config.FileIDs

	return s
}

func (r *updateConfigRequest) updateConfig(oldConfig ComponentConfiguration) ComponentConfiguration {
	// Use the old ComponentConfiguration as a basis for the updated config
	s := oldConfig

	if r.Config.Name != "" {
		s.Name = r.Config.Name
	}

	if r.Config.ICID != 0 {
		s.ICID = r.Config.ICID
	}

	s.FileIDs = r.Config.FileIDs

	// only update Params if not empty
	var emptyJson postgres.Jsonb
	// Serialize empty json and params
	emptyJson_ser, _ := json.Marshal(emptyJson)
	startParams_ser, _ := json.Marshal(r.Config.StartParameters)
	opts := jsondiff.DefaultConsoleOptions()
	diff, _ := jsondiff.Compare(emptyJson_ser, startParams_ser, &opts)
	if diff.String() != "FullMatch" {
		s.StartParameters = r.Config.StartParameters
	}

	return s
}
