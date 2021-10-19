/** helper package, AMQP client.
*
* @author Sonja Happ <sonja.happ@eonerc.rwth-aachen.de>
* @copyright 2014-2021, Institute for Automation of Complex Power Systems, EONERC
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

package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type AMQPclient struct {
	connection *amqp.Connection
	sendCh     *amqp.Channel
	recvCh     *amqp.Channel
}

type Action struct {
	Act        string          `json:"action"`
	When       int64           `json:"when"`
	Parameters json.RawMessage `json:"parameters,omitempty"`
	Model      json.RawMessage `json:"model,omitempty"`
	Results    json.RawMessage `json:"results,omitempty"`
}

type Container struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type TemplateSpec struct {
	Containers []Container `json:"containers"`
}

type JobTemplate struct {
	Spec TemplateSpec `json:"spec"`
}

type JobSpec struct {
	Active   string      `json:"activeDeadlineSeconds"`
	Template JobTemplate `json:"template"`
}

type JobMetaData struct {
	JobName string `json:"name"`
}

type KubernetesJob struct {
	Spec     JobSpec     `json:"spec"`
	MetaData JobMetaData `json:"metadata"`
}

type ICPropertiesToCopy struct {
	Job         KubernetesJob `json:"job"`
	UUID        string        `json:"uuid"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Location    string        `json:"location"`
	Owner       string        `json:"owner"`
	Category    string        `json:"category"`
	Type        string        `json:"type"`
}

type ICUpdateToCopy struct {
	Properties ICPropertiesToCopy `json:"properties"`
	Status     json.RawMessage    `json:"status"`
	Schema     json.RawMessage    `json:"schema"`
}

var client AMQPclient

const VILLAS_EXCHANGE = "villas"

type callback func(amqp.Delivery) error

func ConnectAMQP(uri string, cb callback) error {

	var err error

	// connect to broker
	client.connection, err = amqp.Dial(uri)
	if err != nil {
		return fmt.Errorf("AMQP: failed to connect to RabbitMQ broker %v, error: %v", uri, err)
	}

	// create sendCh
	client.sendCh, err = client.connection.Channel()
	if err != nil {
		return fmt.Errorf("AMQP: failed to open a sendCh, error: %v", err)
	}
	// declare exchange
	err = client.sendCh.ExchangeDeclare(VILLAS_EXCHANGE,
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
	ICQueue, err := client.sendCh.QueueDeclare("",
		false,
		true,
		true,
		false,
		amqp.Table{
			"x-max-length-bytes": int32(32 << 20),
			"x-message-ttl":      int32(10 * 60),
		})
	if err != nil {
		return fmt.Errorf("AMQP: failed to declare the queue, error: %v", err)
	}

	err = client.sendCh.QueueBind(ICQueue.Name, "", VILLAS_EXCHANGE, false, nil)
	if err != nil {
		return fmt.Errorf("AMQP: failed to bind the queue, error: %v", err)
	}

	// create receive channel
	client.recvCh, err = client.connection.Channel()
	if err != nil {
		return fmt.Errorf("AMQP: failed to open a recvCh, error: %v", err)
	}

	// start deliveries
	messages, err := client.recvCh.Consume(ICQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return fmt.Errorf("AMQP: failed to start deliveries: %v", err)
	}

	// consume deliveries
	go func() {
		for {
			for message := range messages {
				err = cb(message)
				if err != nil {
					log.Println("AMQP: Error processing message: ", err.Error())
				}
			}
		}
	}()

	log.Printf(" AMQP: Waiting for messages... ")

	return nil
}

func SendActionAMQP(action Action, destinationUUID string) error {

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
	msg.Headers = make(map[string]interface{}) // empty map
	msg.Headers["uuid"] = destinationUUID

	err = CheckConnection()
	if err != nil {
		return err
	}

	//log.Println("AMQP: Sending message", string(msg.Body))
	err = client.sendCh.Publish(VILLAS_EXCHANGE,
		"",
		false,
		false,
		msg)
	if err != nil {
		return err
	}

	return PublishAMQP(msg)
}

func PublishAMQP(msg amqp.Publishing) error {
	err := client.sendCh.Publish(VILLAS_EXCHANGE,
		"",
		false,
		false,
		msg)
	return err
}

func SendPing(uuid string) error {
	var ping Action
	ping.Act = "ping"

	payload, err := json.Marshal(ping)
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
	msg.Headers = make(map[string]interface{}) // empty map
	msg.Headers["uuid"] = uuid                 // leave uuid empty if ping should go to all ICs

	err = CheckConnection()
	if err != nil {
		return err
	}

	err = client.sendCh.Publish(VILLAS_EXCHANGE,
		"",
		false,
		false,
		msg)
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

// WARNING: this only works with the kubernetes-simple manager of VILLAScontroller
func RequestICcreateAMQPsimpleManager(ic *database.InfrastructureComponent, managerUUID string, userName string) (string, error) {
	newUUID := uuid.New().String()
	log.Printf("New IC UUID: %s", newUUID)

	var lastUpdate ICUpdateToCopy
	log.Println(ic.StatusUpdateRaw.RawMessage)
	err := json.Unmarshal(ic.StatusUpdateRaw.RawMessage, &lastUpdate)
	if err != nil {
		return newUUID, err
	}

	msg := `{"name": "` + lastUpdate.Properties.Name + ` ` + userName + `",` +
		`"location": "` + lastUpdate.Properties.Location + `",` +
		`"category": "` + lastUpdate.Properties.Category + `",` +
		`"type": "` + lastUpdate.Properties.Type + `",` +
		`"uuid": "` + newUUID + `",` +
		`"jobname": "` + lastUpdate.Properties.Job.MetaData.JobName + `-` + userName + `",` +
		`"activeDeadlineSeconds": "` + lastUpdate.Properties.Job.Spec.Active + `",` +
		`"containername": "` + lastUpdate.Properties.Job.Spec.Template.Spec.Containers[0].Name + `-` + userName + `",` +
		`"image": "` + lastUpdate.Properties.Job.Spec.Template.Spec.Containers[0].Image + `",` +
		`"uuid": "` + newUUID + `"}`

	actionCreate := Action{
		Act:        "create",
		When:       time.Now().Unix(),
		Parameters: json.RawMessage(msg),
	}

	err = SendActionAMQP(actionCreate, managerUUID)

	return newUUID, err
}
