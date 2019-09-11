package simulationmodel

import (
	"encoding/json"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nsf/jsondiff"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type validNewSimulationModel struct {
	Name            string         `form:"Name" validate:"required"`
	ScenarioID      uint           `form:"ScenarioID" validate:"required"`
	SimulatorID     uint           `form:"SimulatorID" validate:"required"`
	StartParameters postgres.Jsonb `form:"StartParameters" validate:"required"`
}

type validUpdatedSimulationModel struct {
	Name            string         `form:"Name" validate:"omitempty"`
	SimulatorID     uint           `form:"SimulatorID" validate:"omitempty"`
	StartParameters postgres.Jsonb `form:"StartParameters" validate:"omitempty"`
}

type addSimulationModelRequest struct {
	validNewSimulationModel `json:"model"`
}

type updateSimulationModelRequest struct {
	validUpdatedSimulationModel `json:"model"`
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

	s.Name = r.Name
	s.ScenarioID = r.ScenarioID
	s.SimulatorID = r.SimulatorID
	s.StartParameters = r.StartParameters

	return s
}

func (r *updateSimulationModelRequest) updatedSimulationModel(oldSimulationModel SimulationModel) SimulationModel {
	// Use the old SimulationModel as a basis for the updated Simulation model
	s := oldSimulationModel

	if r.Name != "" {
		s.Name = r.Name
	}

	if r.SimulatorID != 0 {
		s.SimulatorID = r.SimulatorID
	}

	// only update Params if not empty
	var emptyJson postgres.Jsonb
	// Serialize empty json and params
	emptyJson_ser, _ := json.Marshal(emptyJson)
	startParams_ser, _ := json.Marshal(r.StartParameters)
	opts := jsondiff.DefaultConsoleOptions()
	diff, _ := jsondiff.Compare(emptyJson_ser, startParams_ser, &opts)
	if diff.String() != "FullMatch" {
		s.StartParameters = r.StartParameters
	}

	return s
}
