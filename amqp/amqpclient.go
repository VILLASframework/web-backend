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
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"github.com/streadway/amqp"
	"github.com/tidwall/gjson"
	"log"
	"strings"
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
	Model      struct{} `json:"model"`
	Results    struct{} `json:"results"`
}

var client AMQPclient

func ConnectAMQP(uri string) error {

	var err error

	// connect to broker
	client.connection, err = amqp.Dial(uri)
	if err != nil {
		return fmt.Errorf("AMQP: failed to connect to RabbitMQ broker %v", uri)
	}

	// create channel
	client.channel, err = client.connection.Channel()
	if err != nil {
		return fmt.Errorf("AMQP: failed to open a channel")
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
		return fmt.Errorf("AMQP: failed to declare the exchange")
	}

	// add a queue for the ICs
	ICQueue, err := client.channel.QueueDeclare("infrastructure_components",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return fmt.Errorf("AMQP: failed to declare the queue")
	}

	err = client.channel.QueueBind(ICQueue.Name, "", VILLAS_EXCHANGE, false, nil)
	if err != nil {
		return fmt.Errorf("AMQP: failed to bind the queue")
	}

	// consume deliveries
	client.replies, err = client.channel.Consume(ICQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		return fmt.Errorf("AMQP: failed to consume deliveries")
	}

	// consuming queue
	go func() {
		for message := range client.replies {
			//err = message.Ack(false)
			//if err != nil {
			//	fmt.Println("AMQP: Unable to ack message:", err)
			//}

			content := string(message.Body)

			// any action message sent by the VILLAScontroller should be ignored by the web backend
			if strings.Contains(content, "action") {
				continue
			}

			var sToBeUpdated database.InfrastructureComponent
			db := database.GetDB()
			ICUUID := gjson.Get(content, "properties.uuid").String()
			if ICUUID == "" {
				log.Println("AMQP: Could not extract UUID of IC from content of received message, COMPONENT NOT UPDATED")
			} else {
				err = db.Where("UUID = ?", ICUUID).Find(sToBeUpdated).Error
				if err != nil {
					log.Println("AMQP: Unable to find IC with UUID: ", gjson.Get(content, "properties.uuid"), " DB error message: ", err)
				}

				err = db.Model(&sToBeUpdated).Updates(map[string]interface{}{
					"Host":          gjson.Get(content, "host"),
					"Modeltype":     gjson.Get(content, "model"),
					"Uptime":        gjson.Get(content, "uptime"),
					"State":         gjson.Get(content, "state"),
					"StateUpdateAt": time.Now().String(),
					"RawProperties": gjson.Get(content, "properties"),
				}).Error
				if err != nil {
					log.Println("AMQP: Unable to update IC in DB: ", err)
				}

				log.Println("AMQP: Updated IC with UUID ", gjson.Get(content, "properties.uuid"))
			}
		}
	}()

	return nil
}

func SendActionAMQP(action Action, uuid string) error {

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

	if uuid != "" {
		msg.Headers["uuid"] = uuid
		msg.Headers["action"] = "ping"
	}

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

	err := SendActionAMQP(a, "")
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
