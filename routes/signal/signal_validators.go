package signal

import (
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type validNewSignal struct {
	Name              string `form:"Name" validate:"required"`
	Unit              string `form:"unit" validate:"omitempty"`
	Index             uint   `form:"index" validate:"required"`
	Direction         string `form:"direction" validate:"required,oneof=in out"`
	SimulationModelID uint   `form:"simulationModelID" validate:"required"`
}

type validUpdatedSignal struct {
	Name  string `form:"Name" validate:"omitempty"`
	Unit  string `form:"unit" validate:"omitempty"`
	Index uint   `form:"index" validate:"omitempty"`
}

type addSignalRequest struct {
	validNewSignal `json:"signal"`
}

type updateSignalRequest struct {
	validUpdatedSignal `json:"signal"`
}

func (r *addSignalRequest) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *validUpdatedSignal) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *addSignalRequest) createSignal() Signal {
	var s Signal

	s.Name = r.Name
	s.Unit = r.Unit
	s.Index = r.Index
	s.Direction = r.Direction
	s.SimulationModelID = r.SimulationModelID

	return s
}

func (r *updateSignalRequest) updatedSignal(oldSignal Signal) Signal {
	// Use the old Signal as a basis for the updated Signal `s`
	s := oldSignal

	if r.Name != "" {
		s.Name = r.Name
	}

	if r.Index != 0 {
		// TODO this implies that we start indexing at 1
		s.Index = r.Index
	}

	if r.Unit != "" {
		s.Unit = r.Unit
	}

	return s
}
