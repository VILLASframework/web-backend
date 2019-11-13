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
	Name string `form:"Name" validate:"omitempty" json:"name"`
	Grid int    `form:"Grid" validate:"omitempty" json:"grid"`
}

type addDashboardRequest struct {
	Dashboard validNewDashboard `json:"dashboard"`
}

type updateDashboardRequest struct {
	Dashboard validUpdatedDashboard `json:"dashboard"`
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

	s.Name = r.Dashboard.Name
	s.Grid = r.Dashboard.Grid
	s.ScenarioID = r.Dashboard.ScenarioID

	return s
}

func (r *updateDashboardRequest) updatedDashboard(oldDashboard Dashboard) Dashboard {
	// Use the old Dashboard as a basis for the updated Dashboard `s`
	s := oldDashboard

	if r.Dashboard.Name != "" {
		s.Name = r.Dashboard.Name
	}

	if r.Dashboard.Grid != 0 {
		s.Grid = r.Dashboard.Grid
	}

	return s
}
