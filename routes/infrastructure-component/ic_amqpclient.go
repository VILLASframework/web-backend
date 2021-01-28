/** AMQP package, client.
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
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/streadway/amqp"
	"log"
	"time"
)

const VILLAS_EXCHANGE = "villas"

type AMQPclient struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	replies    <-chan amqp.Delivery
}

type Action struct {
	Act        string `json:"action"`
	When       int64  `json:"when"`
	Properties struct {
		UUID        *string `json:"uuid"`
		Name        *string `json:"name"`
		Category    *string `json:"category"`
		Type        *string `json:"type"`
		Location    *string `json:"location"`
		WS_url      *string `json:"ws_url"`
		API_url     *string `json:"api_url"`
		Description *string `json:"description"`
	} `json:"properties"`
}

type ICStatus struct {
	State   *string  `json:"state"`
	Version *string  `json:"version"`
	Uptime  *float64 `json:"uptime"`
	Result  *string  `json:"result"`
	Error   *string  `json:"error"`
}

type ICProperties struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Location    *string `json:"location"`
	Owner       *string `json:"owner"`
	WS_url      *string `json:"ws_url"`
	API_url     *string `json:"api_url"`
	Category    *string `json:"category"`
	Type        *string `json:"type"`
}

type ICUpdate struct {
	Status     *ICStatus     `json:"status"`
	Properties *ICProperties `json:"properties"`
	When       *float64      `json:"when"`
	// TODO add JSON start parameter scheme
}

var client AMQPclient

func ConnectAMQP(uri string) error {

	var err error

	// connect to broker
	client.connection, err = amqp.Dial(uri)
	if err != nil {
		return fmt.Errorf("AMQP: failed to connect to RabbitMQ broker %v, error: %v", uri, err)
	}

	// create channel
	client.channel, err = client.connection.Channel()
	if err != nil {
		return fmt.Errorf("AMQP: failed to open a channel, error: %v", err)
	}
	// declare exchange
	err = client.channel.ExchangeDeclare(VILLAS_EXCHANGE,
		"headers",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return fmt.Errorf("AMQP: failed to declare the exchange, error: %v", err)
	}

	// add a queue for the ICs
	ICQueue, err := client.channel.QueueDeclare("infrastructure_components",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return fmt.Errorf("AMQP: failed to declare the queue, error: %v", err)
	}

	err = client.channel.QueueBind(ICQueue.Name, "", VILLAS_EXCHANGE, false, nil)
	if err != nil {
		return fmt.Errorf("AMQP: failed to bind the queue, error: %v", err)
	}

	// consume deliveries
	client.replies, err = client.channel.Consume(ICQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return fmt.Errorf("AMQP: failed to consume deliveries, error: %v", err)
	}

	// consuming queue
	go func() {
		for {
			for message := range client.replies {
				err = processMessage(message)
				if err != nil {
					log.Println(err.Error())
				}
			}
			time.Sleep(2) // sleep for 2 sek
		}
	}()

	log.Printf(" AMQP: Waiting for messages... ")

	return nil
}

func sendActionAMQP(action Action) error {

	payload, err := json.Marshal(action)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		DeliveryMode:    2,
		Timestamp:       time.Now(),
		ContentType:     "application/json",
		ContentEncoding: "utf-8",
		Priority:        0,
		Body:            payload,
	}

	// set message headers
	var headers map[string]interface{}
	headers = make(map[string]interface{}) // empty map
	if action.Properties.UUID != nil {
		headers["uuid"] = *action.Properties.UUID
	}
	if action.Properties.Type != nil {
		headers["type"] = *action.Properties.Type
	}
	if action.Properties.Category != nil {
		headers["category"] = *action.Properties.Category
	}
	msg.Headers = headers

	err = CheckConnection()
	if err != nil {
		return err
	}

	//log.Println("AMQP: Sending message", string(msg.Body))
	err = client.channel.Publish(VILLAS_EXCHANGE,
		"",
		false,
		false,
		msg)
	return err

}

//func PingAMQP() error {
//	log.Println("AMQP: sending ping command to all ICs")
//
//	var a Action
//	a.Act = "ping"
//	*a.Properties.UUID = ""
//
//	err := sendActionAMQP(a)
//	return err
//}

func CheckConnection() error {

	if client.connection != nil {
		if client.connection.IsClosed() {
			return fmt.Errorf("connection to broker is closed")
		}
	} else {
		return fmt.Errorf("connection is nil")
	}

	return nil
}

func StartAMQP(AMQPurl string, api *gin.RouterGroup) error {
	if AMQPurl != "" {
		log.Println("Starting AMQP client")

		err := ConnectAMQP(AMQPurl)
		if err != nil {
			return err
		}

		// register IC action endpoint only if AMQP client is used
		RegisterAMQPEndpoint(api.Group("/ic"))

		// Periodically call the Ping function to check which ICs are still there
		ticker := time.NewTicker(10 * time.Second)
		go func() {

			for {
				select {
				case <-ticker.C:
					//TODO Add a useful regular event here
					/*
						err = PingAMQP()
						if err != nil {
							log.Println("AMQP Error: ", err.Error())
						}
					*/
				}
			}

		}()

		log.Printf("Connected AMQP client to %s", AMQPurl)
	}

	return nil
}

func processMessage(message amqp.Delivery) error {

	var payload ICUpdate
	err := json.Unmarshal(message.Body, &payload)
	if err != nil {
		return fmt.Errorf("AMQP: Could not unmarshal message to JSON: %v err: %v", string(message.Body), err)
	}

	if payload.Status != nil || payload.Properties != nil {
		//log.Println("Processing AMQP message: ", string(message.Body))
		// if a message contains a "state" field, it is an update for an IC

		headers := amqp.Table(message.Headers)
		ICUUID := fmt.Sprintf("%v", headers["uuid"])
		_, err = uuid.Parse(ICUUID)

		if err != nil {
			return fmt.Errorf("AMQP: UUID not valid: %v, message ignored: %v \n", ICUUID, string(message.Body))
		}
		var sToBeUpdated InfrastructureComponent
		err = sToBeUpdated.byUUID(ICUUID)

		if err == gorm.ErrRecordNotFound {
			// create new record
			err = createExternalIC(payload, ICUUID)
		} else if err != nil {
			// database error
			err = fmt.Errorf("AMQP: Database error for IC %v DB error message: %v", ICUUID, err)
		} else {
			// update record based on payload
			err = sToBeUpdated.updateExternalIC(payload)
		}
	} else {
		log.Println("INFO: ignoring message, payload neither contains status nor properties", message)
	}

	return err
}

func createExternalIC(payload ICUpdate, ICUUID string) error {

	if payload.Properties == nil {
		return fmt.Errorf("AMQP: Cannot create new IC, Propertie field missing")
	}

	var newICReq AddICRequest
	newICReq.InfrastructureComponent.UUID = ICUUID
	if payload.Properties.Name == nil ||
		payload.Properties.Category == nil ||
		payload.Properties.Type == nil {
		// cannot create new IC because required information (name, type, and/or category missing)
		return fmt.Errorf("AMQP: Cannot create new IC, required field(s) is/are missing: name, type, category")
	}
	newICReq.InfrastructureComponent.Name = *payload.Properties.Name
	newICReq.InfrastructureComponent.Category = *payload.Properties.Category
	newICReq.InfrastructureComponent.Type = *payload.Properties.Type

	// add optional params
	if payload.Status != nil {
		if payload.Status.State != nil {
			newICReq.InfrastructureComponent.State = *payload.Status.State
		} else {
			newICReq.InfrastructureComponent.State = "unknown"
		}
		if newICReq.InfrastructureComponent.State == "gone" {
			// Check if state is "gone" and abort creation of IC in this case
			log.Println("AMQP: Aborting creation of IC with state gone")
			return nil
		}

		if payload.Status.Uptime != nil {
			newICReq.InfrastructureComponent.Uptime = *payload.Status.Uptime
		}
	}

	if payload.Properties.WS_url != nil {
		newICReq.InfrastructureComponent.WebsocketURL = *payload.Properties.WS_url
	}
	if payload.Properties.API_url != nil {
		newICReq.InfrastructureComponent.APIURL = *payload.Properties.API_url
	}
	if payload.Properties.Location != nil {
		newICReq.InfrastructureComponent.Location = *payload.Properties.Location
	}
	if payload.Properties.Description != nil {
		newICReq.InfrastructureComponent.Description = *payload.Properties.Description
	}

	// TODO add JSON start parameter scheme

	// set managed externally to true because this IC is created via AMQP
	newICReq.InfrastructureComponent.ManagedExternally = newTrue()

	// set raw status update if IC
	payloadRaw, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("AMQP: failed to marshal raw payload: %v", err)
	}
	newICReq.InfrastructureComponent.StatusUpdateRaw = postgres.Jsonb{RawMessage: payloadRaw}

	// Validate the new IC
	err = newICReq.validate()
	if err != nil {
		return fmt.Errorf("AMQP: Validation of new IC failed: %v", err)
	}

	// Create the new IC
	newIC, err := newICReq.createIC(true)
	if err != nil {
		return fmt.Errorf("AMQP: Creating new IC failed: %v", err)
	}

	// save IC
	err = newIC.save()
	if err != nil {
		return fmt.Errorf("AMQP: Saving new IC to DB failed: %v", err)
	}

	log.Println("AMQP: Created IC with UUID ", newIC.UUID)
	return nil
}

func (s *InfrastructureComponent) updateExternalIC(payload ICUpdate) error {

	var updatedICReq UpdateICRequest

	if payload.Status != nil {
		if payload.Status.State != nil {
			updatedICReq.InfrastructureComponent.State = *payload.Status.State

			if *payload.Status.State == "gone" {
				// remove IC from DB
				log.Println("AMQP: Deleting IC with state gone")
				err := s.delete(true)
				if err != nil {
					// if component could not be deleted there are still configurations using it in the DB
					// continue with the update to save the new state of the component and get back to the deletion later
					log.Println("AMQP: Deletion of IC postponed (config(s) associated to it)")
				}

			}
		}

		if payload.Status.Uptime != nil {
			updatedICReq.InfrastructureComponent.Uptime = *payload.Status.Uptime
		}
	}

	if payload.Properties != nil {
		if payload.Properties.Type != nil {
			updatedICReq.InfrastructureComponent.Type = *payload.Properties.Type
		}
		if payload.Properties.Category != nil {
			updatedICReq.InfrastructureComponent.Category = *payload.Properties.Category
		}
		if payload.Properties.Name != nil {
			updatedICReq.InfrastructureComponent.Name = *payload.Properties.Name
		}
		if payload.Properties.WS_url != nil {
			updatedICReq.InfrastructureComponent.WebsocketURL = *payload.Properties.WS_url
		}
		if payload.Properties.API_url != nil {
			updatedICReq.InfrastructureComponent.APIURL = *payload.Properties.API_url
		}
		if payload.Properties.Location != nil {
			updatedICReq.InfrastructureComponent.Location = *payload.Properties.Location
		}
		if payload.Properties.Description != nil {
			updatedICReq.InfrastructureComponent.Description = *payload.Properties.Description
		}
	}

	// set raw status update if IC
	payloadRaw, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("AMQP: failed to marshal raw payload: %v", err)
	}
	updatedICReq.InfrastructureComponent.StatusUpdateRaw = postgres.Jsonb{RawMessage: payloadRaw}

	// TODO add JSON start parameter scheme

	// Validate the updated IC
	err = updatedICReq.validate()
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
