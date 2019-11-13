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
	Signal validNewSignal `json:"signal"`
}

type updateSignalRequest struct {
	Signal validUpdatedSignal `json:"signal"`
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

	s.Name = r.Signal.Name
	s.Unit = r.Signal.Unit
	s.Index = r.Signal.Index
	s.Direction = r.Signal.Direction
	s.SimulationModelID = r.Signal.SimulationModelID

	return s
}

func (r *updateSignalRequest) updatedSignal(oldSignal Signal) Signal {
	// Use the old Signal as a basis for the updated Signal `s`
	s := oldSignal

	if r.Signal.Name != "" {
		s.Name = r.Signal.Name
	}

	if r.Signal.Index != 0 {
		// TODO this implies that we start indexing at 1
		s.Index = r.Signal.Index
	}

	if r.Signal.Unit != "" {
		s.Unit = r.Signal.Unit
	}

	return s
}
