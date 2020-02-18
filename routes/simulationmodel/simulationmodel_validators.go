/** Simulationmodel package, validators.
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
	"encoding/json"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nsf/jsondiff"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type validNewSimulationModel struct {
	Name                string         `form:"Name" validate:"required"`
	ScenarioID          uint           `form:"ScenarioID" validate:"required"`
	SimulatorID         uint           `form:"SimulatorID" validate:"required"`
	StartParameters     postgres.Jsonb `form:"StartParameters" validate:"required"`
	SelectedModelFileID uint           `form:"SelectedModelFilID" validate:"omitempty"`
}

type validUpdatedSimulationModel struct {
	Name                string         `form:"Name" validate:"omitempty"`
	SimulatorID         uint           `form:"SimulatorID" validate:"omitempty"`
	StartParameters     postgres.Jsonb `form:"StartParameters" validate:"omitempty"`
	SelectedModelFileID uint           `form:"SelectedModelFileID" validate:"omitempty"`
}

type addSimulationModelRequest struct {
	SimulationModel validNewSimulationModel `json:"simulationModel"`
}

type updateSimulationModelRequest struct {
	SimulationModel validUpdatedSimulationModel `json:"simulationModel"`
}

func (r *addSimulationModelRequest) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *validUpdatedSimulationModel) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *addSimulationModelRequest) createSimulationModel() SimulationModel {
	var s SimulationModel

	s.Name = r.SimulationModel.Name
	s.ScenarioID = r.SimulationModel.ScenarioID
	s.SimulatorID = r.SimulationModel.SimulatorID
	s.StartParameters = r.SimulationModel.StartParameters
	s.SelectedModelFileID = r.SimulationModel.SelectedModelFileID

	return s
}

func (r *updateSimulationModelRequest) updatedSimulationModel(oldSimulationModel SimulationModel) SimulationModel {
	// Use the old SimulationModel as a basis for the updated Simulation model
	s := oldSimulationModel

	if r.SimulationModel.Name != "" {
		s.Name = r.SimulationModel.Name
	}

	if r.SimulationModel.SimulatorID != 0 {
		s.SimulatorID = r.SimulationModel.SimulatorID
	}

	if r.SimulationModel.SelectedModelFileID != 0 {
		s.SelectedModelFileID = r.SimulationModel.SelectedModelFileID
	}

	// only update Params if not empty
	var emptyJson postgres.Jsonb
	// Serialize empty json and params
	emptyJson_ser, _ := json.Marshal(emptyJson)
	startParams_ser, _ := json.Marshal(r.SimulationModel.StartParameters)
	opts := jsondiff.DefaultConsoleOptions()
	diff, _ := jsondiff.Compare(emptyJson_ser, startParams_ser, &opts)
	if diff.String() != "FullMatch" {
		s.StartParameters = r.SimulationModel.StartParameters
	}

	return s
}
