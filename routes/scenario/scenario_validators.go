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
	Running         bool           `form:"Running" validate:"omitempty"`
	StartParameters postgres.Jsonb `form:"StartParameters" validate:"required"`
}

type validUpdatedScenario struct {
	Name            string         `form:"Name" validate:"omitempty"`
	Running         bool           `form:"Running" validate:"omitempty"`
	StartParameters postgres.Jsonb `form:"StartParameters" validate:"omitempty"`
}

type addScenarioRequest struct {
	validNewScenario `json:"scenario"`
}

type updateScenarioRequest struct {
	validUpdatedScenario `json:"scenario"`
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

	s.Name = r.Name
	s.Running = r.Running
	s.StartParameters = r.StartParameters

	return s
}

func (r *updateScenarioRequest) updatedScenario(oldScenario Scenario) (Scenario, error) {
	// Use the old Scenario as a basis for the updated Scenario `s`
	s := oldScenario

	if r.Name != "" {
		s.Name = r.Name
	}

	s.Running = r.Running

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

	return s, nil
}
