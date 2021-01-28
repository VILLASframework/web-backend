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
	"github.com/google/uuid"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nsf/jsondiff"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"time"
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
	StartParameterScheme postgres.Jsonb `form:"StartParameterScheme" validate:"omitempty"`
	StatusUpdateRaw      postgres.Jsonb `form:"StatusUpdateRaw" validate:"omitempty"`
	ManagedExternally    *bool          `form:"ManagedExternally" validate:"required"`
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
	StartParameterScheme postgres.Jsonb `form:"StartParameterScheme" validate:"omitempty"`
	StatusUpdateRaw      postgres.Jsonb `form:"StatusUpdateRaw" validate:"omitempty"`
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
	return errs
}

func (r *UpdateICRequest) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *AddICRequest) createIC(receivedViaAMQP bool) (InfrastructureComponent, error) {
	var s InfrastructureComponent
	var err error
	err = nil

	// case distinction for externally managed IC
	if *r.InfrastructureComponent.ManagedExternally && !receivedViaAMQP {
		var action Action
		action.Act = "create"
		action.When = time.Now().Unix()
		action.Properties.Type = new(string)
		action.Properties.Name = new(string)
		action.Properties.Category = new(string)

		*action.Properties.Type = r.InfrastructureComponent.Type
		*action.Properties.Name = r.InfrastructureComponent.Name
		*action.Properties.Category = r.InfrastructureComponent.Category

		// set optional properties
		action.Properties.Description = new(string)
		*action.Properties.Description = r.InfrastructureComponent.Description

		action.Properties.Location = new(string)
		*action.Properties.Location = r.InfrastructureComponent.Location

		action.Properties.API_url = new(string)
		*action.Properties.API_url = r.InfrastructureComponent.APIURL

		action.Properties.WS_url = new(string)
		*action.Properties.WS_url = r.InfrastructureComponent.WebsocketURL

		action.Properties.UUID = new(string)
		*action.Properties.UUID = r.InfrastructureComponent.UUID

		log.Println("AMQP: Sending request to create new IC")
		err = sendActionAMQP(action)
	}

	s.UUID = r.InfrastructureComponent.UUID
	s.WebsocketURL = r.InfrastructureComponent.WebsocketURL
	s.APIURL = r.InfrastructureComponent.APIURL
	s.Type = r.InfrastructureComponent.Type
	s.Name = r.InfrastructureComponent.Name
	s.Category = r.InfrastructureComponent.Category
	s.Location = r.InfrastructureComponent.Location
	s.Description = r.InfrastructureComponent.Description
	s.StartParameterScheme = r.InfrastructureComponent.StartParameterScheme
	s.StatusUpdateRaw = r.InfrastructureComponent.StatusUpdateRaw
	s.ManagedExternally = *r.InfrastructureComponent.ManagedExternally
	s.Uptime = -1.0 // no uptime available
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

	s.UUID = r.InfrastructureComponent.UUID
	s.WebsocketURL = r.InfrastructureComponent.WebsocketURL
	s.APIURL = r.InfrastructureComponent.APIURL
	s.Location = r.InfrastructureComponent.Location
	s.Description = r.InfrastructureComponent.Description

	// set last update time
	s.StateUpdateAt = time.Now().Format(time.RFC1123Z)

	// only update props if not empty
	var emptyJson postgres.Jsonb
	// Serialize empty json and params
	emptyJson_ser, _ := json.Marshal(emptyJson)
	opts := jsondiff.DefaultConsoleOptions()

	startParams_ser, _ := json.Marshal(r.InfrastructureComponent.StartParameterScheme)
	diff, _ := jsondiff.Compare(emptyJson_ser, startParams_ser, &opts)
	if diff.String() != "FullMatch" {
		s.StartParameterScheme = r.InfrastructureComponent.StartParameterScheme
	}

	statusUpdateRaw_ser, _ := json.Marshal(r.InfrastructureComponent.StatusUpdateRaw)
	diff, _ = jsondiff.Compare(emptyJson_ser, statusUpdateRaw_ser, &opts)
	if diff.String() != "FullMatch" {
		s.StatusUpdateRaw = r.InfrastructureComponent.StatusUpdateRaw
	}

	return s
}
