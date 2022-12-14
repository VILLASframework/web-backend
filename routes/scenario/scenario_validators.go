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

package scenario

import (
	"encoding/json"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nsf/jsondiff"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type validNewScenario struct {
	Name            string         `form:"Name" validate:"required"`
	StartParameters postgres.Jsonb `form:"StartParameters" validate:"required"`
}

type validUpdatedScenario struct {
	Name            string         `form:"Name" validate:"omitempty"`
	IsLocked        bool           `form:"IsLocked" validate:"omitempty"`
	StartParameters postgres.Jsonb `form:"StartParameters" validate:"omitempty"`
}

type addScenarioRequest struct {
	Scenario validNewScenario `json:"scenario"`
}

type updateScenarioRequest struct {
	Scenario validUpdatedScenario `json:"scenario"`
}

func (r *addScenarioRequest) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *validUpdatedScenario) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *addScenarioRequest) createScenario() Scenario {
	var s Scenario

	s.Name = r.Scenario.Name
	s.IsLocked = false // new scenarios are not locked
	s.StartParameters = r.Scenario.StartParameters

	return s
}

func (r *updateScenarioRequest) updatedScenario(oldScenario Scenario, userRole string) Scenario {
	// Use the old Scenario as a basis for the updated Scenario `s`
	s := oldScenario

	if userRole == "Admin" { // only admin users can change isLocked status
		s.IsLocked = r.Scenario.IsLocked
	}

	if r.Scenario.Name != "" {
		s.Name = r.Scenario.Name
	}

	// only update Params if not empty
	var emptyJson postgres.Jsonb
	// Serialize empty json and params
	emptyJson_ser, _ := json.Marshal(emptyJson)
	startParams_ser, _ := json.Marshal(r.Scenario.StartParameters)
	opts := jsondiff.DefaultConsoleOptions()
	diff, _ := jsondiff.Compare(emptyJson_ser, startParams_ser, &opts)
	if diff.String() != "FullMatch" {
		s.StartParameters = r.Scenario.StartParameters
	}

	return s
}
