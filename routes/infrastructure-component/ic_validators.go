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
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type validNewIC struct {
	UUID                 string         `form:"UUID" validate:"omitempty"`
	WebsocketURL         string         `form:"WebsocketURL" validate:"omitempty"`
	APIURL               string         `form:"APIURL" validate:"omitempty"`
	Type                 string         `form:"Type" validate:"required"`
	Name                 string         `form:"Name" validate:"required"`
	Category             string         `form:"Category" validate:"required"`
	State                string         `form:"State" validate:"omitempty"`
	Location             string         `form:"Location" validate:"omitempty"`
	Description          string         `form:"Description" validate:"omitempty"`
	StartParameterSchema postgres.Jsonb `form:"StartParameterSchema" validate:"omitempty"`
	StatusUpdateRaw      postgres.Jsonb `form:"StatusUpdateRaw" validate:"omitempty"`
	ManagedExternally    *bool          `form:"ManagedExternally" validate:"required"`
	Manager              string         `form:"Manager" validate:"omitempty"`
	Uptime               float64        `form:"Uptime" validate:"omitempty"`
}

type validUpdatedIC struct {
	UUID                 string         `form:"UUID" validate:"omitempty"`
	WebsocketURL         string         `form:"WebsocketURL" validate:"omitempty"`
	APIURL               string         `form:"APIURL" validate:"omitempty"`
	Type                 string         `form:"Type" validate:"omitempty"`
	Name                 string         `form:"Name" validate:"omitempty"`
	Category             string         `form:"Category" validate:"omitempty"`
	State                string         `form:"State" validate:"omitempty"`
	Location             string         `form:"Location" validate:"omitempty"`
	Description          string         `form:"Description" validate:"omitempty"`
	StartParameterSchema postgres.Jsonb `form:"StartParameterSchema" validate:"omitempty"`
	StatusUpdateRaw      postgres.Jsonb `form:"StatusUpdateRaw" validate:"omitempty"`
	Manager              string         `form:"Manager" validate:"omitempty"`
	Uptime               float64        `form:"Uptime" validate:"omitempty"`
}

type AddICRequest struct {
	InfrastructureComponent validNewIC `json:"ic"`
}

type UpdateICRequest struct {
	InfrastructureComponent validUpdatedIC `json:"ic"`
}

func (r *AddICRequest) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	if errs != nil {
		return errs
	}

	// check if uuid is valid
	_, errs = uuid.Parse(r.InfrastructureComponent.UUID)
	if errs != nil {
		return errs
	}

	if *r.InfrastructureComponent.ManagedExternally == true {
		// check if valid manager UUID is provided
		_, errs = uuid.Parse(r.InfrastructureComponent.Manager)
		if errs != nil {
			return errs
		}
	}

	return errs
}

func (r *UpdateICRequest) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	if errs != nil {
		return errs
	}

	if r.InfrastructureComponent.Manager != "" {
		//check if valid manager UUID is provided
		_, errs = uuid.Parse(r.InfrastructureComponent.Manager)
		if errs != nil {
			return errs
		}
	}

	return errs
}

func (r *AddICRequest) createIC() (InfrastructureComponent, error) {
	var s InfrastructureComponent
	var err error
	err = nil

	s.UUID = r.InfrastructureComponent.UUID
	s.WebsocketURL = r.InfrastructureComponent.WebsocketURL
	s.APIURL = r.InfrastructureComponent.APIURL
	s.Type = r.InfrastructureComponent.Type
	s.Name = r.InfrastructureComponent.Name
	s.Category = r.InfrastructureComponent.Category
	s.Location = r.InfrastructureComponent.Location
	s.Description = r.InfrastructureComponent.Description
	s.StartParameterSchema = r.InfrastructureComponent.StartParameterSchema
	s.StatusUpdateRaw = r.InfrastructureComponent.StatusUpdateRaw
	s.ManagedExternally = *r.InfrastructureComponent.ManagedExternally
	s.Manager = r.InfrastructureComponent.Manager
	s.Uptime = math.Round(r.InfrastructureComponent.Uptime) // round required for backward compatibility of data model
	if r.InfrastructureComponent.State != "" {
		s.State = r.InfrastructureComponent.State
	} else {
		s.State = "unknown"
	}
	// set last update to creation time of IC
	s.StateUpdateAt = time.Now().Format(time.RFC1123Z)

	return s, err
}

func (r *UpdateICRequest) updatedIC(oldIC InfrastructureComponent) InfrastructureComponent {
	// Use the old InfrastructureComponent as a basis for the updated InfrastructureComponent `s`
	s := oldIC
	s.Type = r.InfrastructureComponent.Type
	s.Name = r.InfrastructureComponent.Name
	s.Category = r.InfrastructureComponent.Category
	s.State = r.InfrastructureComponent.State
	s.UUID = r.InfrastructureComponent.UUID
	s.WebsocketURL = r.InfrastructureComponent.WebsocketURL
	s.APIURL = r.InfrastructureComponent.APIURL
	s.Location = r.InfrastructureComponent.Location
	s.Description = r.InfrastructureComponent.Description
	s.Uptime = math.Round(r.InfrastructureComponent.Uptime) // round required for backward compatibility of data model
	s.Manager = r.InfrastructureComponent.Manager
	s.StartParameterSchema = r.InfrastructureComponent.StartParameterSchema
	s.StatusUpdateRaw = r.InfrastructureComponent.StatusUpdateRaw

	// set last update time
	s.StateUpdateAt = time.Now().Format(time.RFC1123Z)
	return s
}
