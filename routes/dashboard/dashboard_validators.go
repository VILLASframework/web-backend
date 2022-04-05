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

package dashboard

import (
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type validNewDashboard struct {
	Name       string `form:"Name" validate:"required"`
	Grid       int    `form:"Grid" validate:"required"`
	Height     int    `form:"Height" validate:"omitempty"`
	ScenarioID uint   `form:"ScenarioID" validate:"required"`
}

type validUpdatedDashboard struct {
	Name   string `form:"Name" validate:"omitempty" json:"name"`
	Height int    `form:"Height" validate:"omitempty" json:"height"`
	Grid   int    `form:"Grid" validate:"omitempty" json:"grid"`
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
	s.Height = r.Dashboard.Height
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

	if r.Dashboard.Height > 0 {
		s.Height = r.Dashboard.Height
	}

	return s
}
