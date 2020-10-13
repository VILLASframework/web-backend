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
package amqp

import (
	"encoding/json"
	"fmt"
	infrastructure_component "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/infrastructure-component"
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
	Act        string   `json:"action"`
	When       float32  `json:"when"`
	Parameters struct{} `json:"parameters"`
	UUID       *string  `json:"uuid"`
	//Model      struct{} `json:"model"`
	//Results    struct{} `json:"results"`
}

type ICUpdate struct {
	State      *string `json:"state"`
	Properties struct {
		UUID     string  `json:"uuid"`
		Name     *string `json:"name"`
		Category *string `json:"category"`
		Type     *string `json:"type"`
		Location *string `json:"location"`
	} `json:"properties"`
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
				processMessage(message)
			}
			time.Sleep(2) // sleep for 2 sek
		}
	}()

	log.Printf(" AMQP: Waiting for messages... ")

	return nil
}

func SendActionAMQP(action Action) error {

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

	err = CheckConnection()
	if err != nil {
		return err
	}

	log.Println("AMQP: Sending message", string(msg.Body))
	err = client.channel.Publish(VILLAS_EXCHANGE,
		"",
		false,
		false,
		msg)
	return err

}

func PingAMQP() error {
	log.Println("AMQP: sending ping command to all ICs")

	var a Action
	a.Act = "ping"
	*a.UUID = ""

	err := SendActionAMQP(a)
	return err
}

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

func processMessage(message amqp.Delivery) {

	log.Println("Processing AMQP message: ", string(message.Body))

	var payload ICUpdate
	err := json.Unmarshal(message.Body, &payload)
	if err != nil {
		log.Println("AMQP: Could not unmarshal message to JSON:", string(message.Body), "err: ", err)
		return
	}

	ICUUID := payload.Properties.UUID
	_, err = uuid.Parse(ICUUID)

	if err != nil {
		log.Printf("AMQP: UUID not valid: %v, message ignored: %v \n", ICUUID, string(message.Body))
		return
	}

	var sToBeUpdated infrastructure_component.InfrastructureComponent
	err = sToBeUpdated.ByUUID(ICUUID)

	if err == gorm.ErrRecordNotFound {
		// create new record
		var newICReq infrastructure_component.AddICRequest
		newICReq.InfrastructureComponent.UUID = payload.Properties.UUID
		if payload.Properties.Name == nil ||
			payload.Properties.Category == nil ||
			payload.Properties.Type == nil {
			// cannot create new IC because required information (name, type, and/or category missing)
			log.Println("AMQP: Cannot create new IC, required field(s) is/are missing: name, type, category")
			return
		}
		newICReq.InfrastructureComponent.Name = *payload.Properties.Name
		newICReq.InfrastructureComponent.Category = *payload.Properties.Category
		newICReq.InfrastructureComponent.Type = *payload.Properties.Type

		// add optional params
		if payload.State != nil {
			newICReq.InfrastructureComponent.State = *payload.State
		} else {
			newICReq.InfrastructureComponent.State = "unknown"
		}
		if payload.Properties.Location != nil {
			newICReq.InfrastructureComponent.Properties = postgres.Jsonb{json.RawMessage(`{"location" : " ` + *payload.Properties.Location + `"}`)}
		}

		// Validate the new IC
		err = newICReq.Validate()
		if err != nil {
			log.Println("AMQP: Validation of new IC failed:", err)
			return
		}

		// Create the new IC
		newIC := newICReq.CreateIC()

		// save IC
		err = newIC.Save()
		if err != nil {
			log.Println("AMQP: Saving new IC to DB failed:", err)
			return
		}

		log.Println("AMQP: Created IC ", newIC.Name)

	} else if err != nil {
		log.Println("AMQP: Database error for IC", ICUUID, " DB error message: ", err)
		return
	} else {

		var updatedICReq infrastructure_component.UpdateICRequest
		if payload.State != nil {
			updatedICReq.InfrastructureComponent.State = *payload.State
		}
		if payload.Properties.Type != nil {
			updatedICReq.InfrastructureComponent.Type = *payload.Properties.Type
		}
		if payload.Properties.Category != nil {
			updatedICReq.InfrastructureComponent.Category = *payload.Properties.Category
		}
		if payload.Properties.Name != nil {
			updatedICReq.InfrastructureComponent.Name = *payload.Properties.Name
		}
		if payload.Properties.Location != nil {
			updatedICReq.InfrastructureComponent.Properties = postgres.Jsonb{json.RawMessage(`{"location" : " ` + *payload.Properties.Location + `"}`)}
		}

		// Validate the updated IC
		if err = updatedICReq.Validate(); err != nil {
			log.Println("AMQP: Validation of updated IC failed:", err)
			return
		}

		// Create the updated IC from old IC
		updatedIC := updatedICReq.UpdatedIC(sToBeUpdated)

		// Finally update the IC in the DB
		err = sToBeUpdated.Update(updatedIC)
		if err != nil {
			log.Println("AMQP: Unable to update IC", sToBeUpdated.Name, "in DB: ", err)
			return
		}

		log.Println("AMQP: Updated IC ", sToBeUpdated.Name)
	}
}
