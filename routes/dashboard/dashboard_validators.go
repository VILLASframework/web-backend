package dashboard

import (
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type validNewDashboard struct {
	Name       string `form:"Name" validate:"required"`
	Grid       int    `form:"Grid" validate:"required"`
	ScenarioID uint   `form:"ScenarioID" validate:"required"`
}

type validUpdatedDashboard struct {
	Name       string `form:"Name" validate:"omitempty"`
	Grid       int    `form:"Grid" validate:"omitempty"`
	ScenarioID uint   `form:"ScenarioID" validate:"omitempty"`
}

type addDashboardRequest struct {
	validNewDashboard `json:"dashboard"`
}

type updateDashboardRequest struct {
	validUpdatedDashboard `json:"dashboard"`
}

func (r *addDashboardRequest) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *validUpdatedDashboard) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *addDashboardRequest) createDashboard() Dashboard {
	var s Dashboard

	s.Name = r.Name
	s.Grid = r.Grid
	s.ScenarioID = r.ScenarioID

	return s
}

func (r *updateDashboardRequest) updatedDashboard(oldDashboard Dashboard) (Dashboard, error) {
	// Use the old Dashboard as a basis for the updated Dashboard `s`
	s := oldDashboard

	if r.Name != "" {
		s.Name = r.Name
	}

	if r.Grid != 0 {
		s.Grid = r.Grid
	}

	if r.ScenarioID != 0 {
		// TODO do we allow this case?
		//s.ScenarioID = r.ScenarioID
	}

	return s, nil
}
