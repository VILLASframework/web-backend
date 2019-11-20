/** Simulator package, validators.
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
	"encoding/json"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nsf/jsondiff"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type validNewSimulator struct {
	UUID       string         `form:"UUID" validate:"required"`
	Host       string         `form:"Host" validate:"required"`
	Modeltype  string         `form:"Modeltype" validate:"required"`
	Properties postgres.Jsonb `form:"Properties" validate:"required"`
	State      string         `form:"State"`
}

type validUpdatedSimulator struct {
	UUID       string         `form:"UUID" validate:"omitempty"`
	Host       string         `form:"Host" validate:"omitempty"`
	Modeltype  string         `form:"Modeltype" validate:"omitempty"`
	Properties postgres.Jsonb `form:"Properties" validate:"omitempty"`
	State      string         `form:"State" validate:"omitempty"`
}

type addSimulatorRequest struct {
	Simulator validNewSimulator `json:"simulator"`
}

type updateSimulatorRequest struct {
	Simulator validUpdatedSimulator `json:"simulator"`
}

func (r *addSimulatorRequest) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *validUpdatedSimulator) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *addSimulatorRequest) createSimulator() Simulator {
	var s Simulator

	s.UUID = r.Simulator.UUID
	s.Host = r.Simulator.Host
	s.Modeltype = r.Simulator.Modeltype
	s.Properties = r.Simulator.Properties
	if r.Simulator.State != "" {
		s.State = r.Simulator.State
	}
	return s
}

func (r *updateSimulatorRequest) updatedSimulator(oldSimulator Simulator) Simulator {
	// Use the old Simulator as a basis for the updated Simulator `s`
	s := oldSimulator

	if r.Simulator.UUID != "" {
		s.UUID = r.Simulator.UUID
	}

	if r.Simulator.Host != "" {
		s.Host = r.Simulator.Host
	}

	if r.Simulator.Modeltype != "" {
		s.Modeltype = r.Simulator.Modeltype
	}

	if r.Simulator.State != "" {
		s.State = r.Simulator.State
	}

	// only update props if not empty
	var emptyJson postgres.Jsonb
	// Serialize empty json and params
	emptyJson_ser, _ := json.Marshal(emptyJson)
	startParams_ser, _ := json.Marshal(r.Simulator.Properties)
	opts := jsondiff.DefaultConsoleOptions()
	diff, _ := jsondiff.Compare(emptyJson_ser, startParams_ser, &opts)
	if diff.String() != "FullMatch" {
		s.Properties = r.Simulator.Properties
	}

	return s
}
