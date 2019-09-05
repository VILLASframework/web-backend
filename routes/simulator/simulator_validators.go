package simulator

import (
	"github.com/jinzhu/gorm/dialects/postgres"
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
	validNewSimulator `json:"simulator"`
}

type updateSimulatorRequest struct {
	validUpdatedSimulator `json:"simulator"`
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

	s.UUID = r.UUID
	s.Host = r.Host
	s.Modeltype = r.Modeltype
	s.Properties = r.Properties
	if r.State != "" {
		s.State = r.State
	}
	return s
}

func (r *updateSimulatorRequest) updatedSimulator(oldSimulator Simulator) (Simulator, error) {
	// Use the old Simulator as a basis for the updated Simulator `s`
	s := oldSimulator

	if r.UUID != "" {
		s.UUID = r.UUID
	}

	if r.Host != "" {
		s.Host = r.Host
	}

	if r.Modeltype != "" {
		s.Modeltype = r.Modeltype
	}

	if r.State != "" {
		s.State = r.State
	}
	// TODO check for empty properties?
	s.Properties = r.Properties

	return s, nil
}
