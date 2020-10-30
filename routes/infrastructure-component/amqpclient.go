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

type ICUpdate struct {
	Status *struct {
		UUID        string   `json:"uuid"`
		State       *string  `json:"state"`
		Name        *string  `json:"name"`
		Category    *string  `json:"category"`
		Type        *string  `json:"type"`
		Location    *string  `json:"location"`
		WS_url      *string  `json:"ws_url"`
		API_url     *string  `json:"api_url"`
		Description *string  `json:"description"`
		Uptime      *float64 `json:"uptime"` // TODO check if data type of uptime is float64 or int
	} `json:"status"`
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
	*a.Properties.UUID = ""

	err := sendActionAMQP(a)
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

func processMessage(message amqp.Delivery) error {

	log.Println("Processing AMQP message: ", string(message.Body))

	var payload ICUpdate
	err := json.Unmarshal(message.Body, &payload)
	if err != nil {
		return fmt.Errorf("AMQP: Could not unmarshal message to JSON: %v err: %v", string(message.Body), err)
	}

	if payload.Status != nil {
		// if a message contains a "state" field, it is an update for an IC
		ICUUID := payload.Status.UUID
		_, err = uuid.Parse(ICUUID)

		if err != nil {
			return fmt.Errorf("AMQP: UUID not valid: %v, message ignored: %v \n", ICUUID, string(message.Body))
		}
		var sToBeUpdated InfrastructureComponent
		err = sToBeUpdated.ByUUID(ICUUID)

		if err == gorm.ErrRecordNotFound {
			// create new record
			err = createNewICviaAMQP(payload)
		} else if err != nil {
			// database error
			err = fmt.Errorf("AMQP: Database error for IC %v DB error message: %v", ICUUID, err)
		} else {
			// update record based on payload
			err = sToBeUpdated.updateICviaAMQP(payload)
		}
	}
	return err
}
