/** infrastructure-component package, AMQP client.
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
	"fmt"
	"log"
	"strings"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/streadway/amqp"
)

type ICStatus struct {
	State     string  `json:"state"`
	Version   string  `json:"version"`
	Uptime    float64 `json:"uptime"`
	Result    string  `json:"result"`
	Error     string  `json:"error"`
	ManagedBy string  `json:"managed_by"`
}

type ICProperties struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Owner       string `json:"owner"`
	WS_url      string `json:"ws_url"`
	API_url     string `json:"api_url"`
	Category    string `json:"category"`
	Type        string `json:"type"`
}

type ICSchema struct {
	StartParameterSchema   json.RawMessage `json:"start"`
	CreateParametersSchema json.RawMessage `json:"create"`
}

type ICUpdate struct {
	Status     ICStatus     `json:"status"`
	Properties ICProperties `json:"properties"`
	Schema     ICSchema     `json:"schema"`
	When       float64      `json:"when"`
	Action     string       `json:"action"`
}

func StartAMQP(AMQPurl string, api *gin.RouterGroup) error {
	if AMQPurl != "" {
		log.Println("Starting AMQP client")

		err := helper.ConnectAMQP(AMQPurl, ProcessMessage)
		if err != nil {
			return err
		}

		// register IC action endpoint only if AMQP client is used
		RegisterAMQPEndpoint(api.Group("/ic"))

		log.Printf("Connected AMQP client to %s", AMQPurl)
	}

	return nil
}

func ProcessMessage(message amqp.Delivery) error {

	var payload ICUpdate
	err := json.Unmarshal(message.Body, &payload)
	if err != nil {
		return fmt.Errorf("AMQP: Could not unmarshal message to JSON: %v err: %v", string(message.Body), err)
	}

	if payload.Action != "" {
		// if a message contains an action, it is not intended for the backend
		//log.Println("AMQP: Ignoring action message ", payload)
		return nil
	}

	ICUUID := payload.Properties.UUID
	_, err = uuid.Parse(ICUUID)
	if err != nil {
		return fmt.Errorf("AMQP: UUID not valid: %v, message ignored: %v \n", ICUUID, string(message.Body))
	}

	var sToBeUpdated InfrastructureComponent
	err = sToBeUpdated.byUUID(ICUUID)

	if err == gorm.ErrRecordNotFound {
		// create new record
		err = createExternalIC(payload, ICUUID, message.Body)
	} else if err != nil {
		// database error
		err = fmt.Errorf("AMQP: Database error for IC %v DB error message: %v", ICUUID, err)
	} else {
		// update record based on payload
		err = sToBeUpdated.updateExternalIC(payload, message.Body)
	}

	return err
}

func createExternalIC(payload ICUpdate, ICUUID string, body []byte) error {

	var newICReq AddICRequest
	newICReq.InfrastructureComponent.UUID = ICUUID
	newICReq.InfrastructureComponent.Name = payload.Properties.Name
	newICReq.InfrastructureComponent.Category = payload.Properties.Category
	newICReq.InfrastructureComponent.Type = payload.Properties.Type

	// add optional params
	if payload.Status.State != "" {
		newICReq.InfrastructureComponent.State = payload.Status.State
	} else {
		newICReq.InfrastructureComponent.State = "unknown"
	}
	if newICReq.InfrastructureComponent.State == "gone" {
		// Check if state is "gone" and abort creation of IC in this case
		log.Println("AMQP: Aborting creation of IC with state gone")
		return nil
	}

	newICReq.InfrastructureComponent.UUID = payload.Properties.UUID
	newICReq.InfrastructureComponent.Uptime = payload.Status.Uptime
	newICReq.InfrastructureComponent.WebsocketURL = payload.Properties.WS_url
	newICReq.InfrastructureComponent.APIURL = payload.Properties.API_url
	newICReq.InfrastructureComponent.Location = payload.Properties.Location
	newICReq.InfrastructureComponent.Description = payload.Properties.Description
	// set managed externally to true because this IC is created via AMQP
	newICReq.InfrastructureComponent.ManagedExternally = newTrue()
	newICReq.InfrastructureComponent.Manager = payload.Status.ManagedBy
	newICReq.InfrastructureComponent.StartParameterSchema = postgres.Jsonb{RawMessage: payload.Schema.StartParameterSchema}
	newICReq.InfrastructureComponent.CreateParameterSchema = postgres.Jsonb{RawMessage: payload.Schema.CreateParametersSchema}
	// set raw status update if IC
	newICReq.InfrastructureComponent.StatusUpdateRaw = postgres.Jsonb{RawMessage: body}

	// Validate the new IC
	err := newICReq.validate()
	if err != nil {
		return fmt.Errorf("AMQP: Validation of new IC failed: %v", err)
	}

	// Create the new IC
	newIC, err := newICReq.createIC()
	if err != nil {
		return fmt.Errorf("AMQP: Creating new IC failed: %v", err)
	}

	// save IC
	err = newIC.save()
	if err != nil {
		return fmt.Errorf("AMQP: Saving new IC to DB failed: %v", err)
	}

	log.Println("AMQP: Created IC with UUID ", newIC.UUID)

	// send ping to get full status update of this IC
	err = helper.SendPing(ICUUID)
	return err
}

func (s *InfrastructureComponent) updateExternalIC(payload ICUpdate, body []byte) error {

	var updatedICReq UpdateICRequest

	if payload.Status.State != "" {
		updatedICReq.InfrastructureComponent.State = payload.Status.State

		if updatedICReq.InfrastructureComponent.State == "gone" {
			// remove IC from DB
			log.Println("AMQP: Deleting IC with state gone", s.UUID)
			err := s.delete()
			if err != nil {
				// if component could not be deleted there are still configurations using it in the DB
				// continue with the update to save the new state of the component and get back to the deletion later
				if strings.Contains(err.Error(), "postponed") {
					log.Println(err) // print log message
				} else {
					return err // return upon DB error
				}
			} else {
				// if delete was successful, return here and do not run the update
				return nil
			}
		}
	} else {
		updatedICReq.InfrastructureComponent.State = "unknown"
	}

	updatedICReq.InfrastructureComponent.UUID = payload.Properties.UUID
	updatedICReq.InfrastructureComponent.Uptime = payload.Status.Uptime
	updatedICReq.InfrastructureComponent.Type = payload.Properties.Type
	updatedICReq.InfrastructureComponent.Category = payload.Properties.Category
	updatedICReq.InfrastructureComponent.Name = payload.Properties.Name
	updatedICReq.InfrastructureComponent.WebsocketURL = payload.Properties.WS_url
	updatedICReq.InfrastructureComponent.APIURL = payload.Properties.API_url
	updatedICReq.InfrastructureComponent.Location = payload.Properties.Location
	updatedICReq.InfrastructureComponent.Description = payload.Properties.Description
	updatedICReq.InfrastructureComponent.Manager = payload.Status.ManagedBy
	updatedICReq.InfrastructureComponent.StartParameterSchema = postgres.Jsonb{RawMessage: payload.Schema.StartParameterSchema}
	updatedICReq.InfrastructureComponent.CreateParameterSchema = postgres.Jsonb{RawMessage: payload.Schema.CreateParametersSchema}
	// set raw status update if IC
	updatedICReq.InfrastructureComponent.StatusUpdateRaw = postgres.Jsonb{RawMessage: body}

	// Validate the updated IC
	err := updatedICReq.validate()
	if err != nil {
		return fmt.Errorf("AMQP: Validation of updated IC failed: %v", err)
	}

	// Create the updated IC from old IC
	updatedIC := updatedICReq.updatedIC(*s)

	// Finally update the IC in the DB
	err = s.update(updatedIC)
	if err != nil {
		return fmt.Errorf("AMQP: Unable to update IC %v in DB: %v", s.Name, err)
	}

	log.Println("AMQP: Updated IC with UUID ", s.UUID)
	return err
}

func newTrue() *bool {
	b := true
	return &b
}

func newFalse() *bool {
	b := false
	return &b
}
