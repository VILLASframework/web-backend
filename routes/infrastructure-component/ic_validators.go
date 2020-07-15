/** InfrastructureComponent package, validators.
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
package infrastructure_component

import (
	"encoding/json"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nsf/jsondiff"
	"gopkg.in/go-playground/validator.v9"
	"time"
)

var validate *validator.Validate

type validNewIC struct {
	UUID       string         `form:"UUID" validate:"required"`
	Host       string         `form:"Host" validate:"omitempty"`
	APIHost    string         `form:"APIHost" validate:"omitempty"`
	Type       string         `form:"Type" validate:"required"`
	Name       string         `form:"Name" validate:"required"`
	Category   string         `form:"Category" validate:"required"`
	Properties postgres.Jsonb `form:"Properties" validate:"omitempty"`
	State      string         `form:"State" validate:"omitempty"`
}

type validUpdatedIC struct {
	UUID       string         `form:"UUID" validate:"omitempty"`
	Host       string         `form:"Host" validate:"omitempty"`
	APIHost    string         `form:"APIHost" validate:"omitempty"`
	Type       string         `form:"Type" validate:"omitempty"`
	Name       string         `form:"Name" validate:"omitempty"`
	Category   string         `form:"Category" validate:"omitempty"`
	Properties postgres.Jsonb `form:"Properties" validate:"omitempty"`
	State      string         `form:"State" validate:"omitempty"`
}

type addICRequest struct {
	InfrastructureComponent validNewIC `json:"ic"`
}

type updateICRequest struct {
	InfrastructureComponent validUpdatedIC `json:"ic"`
}

func (r *addICRequest) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *validUpdatedIC) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *addICRequest) createIC() InfrastructureComponent {
	var s InfrastructureComponent

	s.UUID = r.InfrastructureComponent.UUID
	s.Host = r.InfrastructureComponent.Host
	s.APIHost = r.InfrastructureComponent.APIHost
	s.Type = r.InfrastructureComponent.Type
	s.Name = r.InfrastructureComponent.Name
	s.Category = r.InfrastructureComponent.Category
	s.Properties = r.InfrastructureComponent.Properties
	if r.InfrastructureComponent.State != "" {
		s.State = r.InfrastructureComponent.State
	} else {
		s.State = "unknown"
	}
	// set last update to creation time of IC
	s.StateUpdateAt = time.Now().Format(time.RFC1123)

	return s
}

func (r *updateICRequest) updatedIC(oldIC InfrastructureComponent) InfrastructureComponent {
	// Use the old InfrastructureComponent as a basis for the updated InfrastructureComponent `s`
	s := oldIC

	if r.InfrastructureComponent.UUID != "" {
		s.UUID = r.InfrastructureComponent.UUID
	}

	if r.InfrastructureComponent.Host != "" {
		s.Host = r.InfrastructureComponent.Host
	}

	if r.InfrastructureComponent.APIHost != "" {
		s.APIHost = r.InfrastructureComponent.APIHost
	}

	if r.InfrastructureComponent.Type != "" {
		s.Type = r.InfrastructureComponent.Type
	}

	if r.InfrastructureComponent.Name != "" {
		s.Name = r.InfrastructureComponent.Name
	}

	if r.InfrastructureComponent.Category != "" {
		s.Category = r.InfrastructureComponent.Category
	}

	if r.InfrastructureComponent.State != "" {
		s.State = r.InfrastructureComponent.State
	}

	// only update props if not empty
	var emptyJson postgres.Jsonb
	// Serialize empty json and params
	emptyJson_ser, _ := json.Marshal(emptyJson)
	startParams_ser, _ := json.Marshal(r.InfrastructureComponent.Properties)
	opts := jsondiff.DefaultConsoleOptions()
	diff, _ := jsondiff.Compare(emptyJson_ser, startParams_ser, &opts)
	if diff.String() != "FullMatch" {
		s.Properties = r.InfrastructureComponent.Properties
	}

	return s
}
