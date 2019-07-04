package common

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/tidwall/gjson"
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
	Act  string  `json:"action"`
	When float32 `json:"when"`

	// TODO add more fields here
}

var client AMQPclient

func ConnectAMQP(uri string) error {

	var err error

	// connect to broker
	client.connection, err = amqp.Dial(uri)
	if err != nil {
		return fmt.Errorf("AMQP: failed to connect to RabbitMQ broker")
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

	// add a queue for the simulators
	simulatorQueue, err := client.channel.QueueDeclare("simulators",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return fmt.Errorf("AMQP: failed to declare the queue")
	}

	err = client.channel.QueueBind(simulatorQueue.Name, "", VILLAS_EXCHANGE, false, nil)
	if err != nil {
		return fmt.Errorf("AMQP: failed to bind the queue")
	}

	// consume deliveries
	client.replies, err = client.channel.Consume(simulatorQueue.Name,
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
			err = message.Ack(false)
			if err != nil {
				fmt.Println("AMQP: Unable to ack message:", err)
			}

			content := string(message.Body)

			if strings.Contains(content, "action") {
				continue
			}

			var sToBeUpdated Simulator
			db := GetDB()
			err = db.Where("UUID = ?", gjson.Get(content, "properties.uuid")).Find(sToBeUpdated).Error
			if err != nil {
				fmt.Println("AMQP: Unable to find simulator with UUID: ", gjson.Get(content, "properties.uuid"), " DB error message: ", err)
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
				fmt.Println("AMQP: Unable to update simulator in DB: ", err)
			}

			fmt.Println("AMQP: Updated simulator with UUID ", gjson.Get(content, "properties.uuid"))
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
	fmt.Println("AMQP: sending ping command to all simulators")

	var a Action
	a.Act = "ping"

	err := SendActionAMQP(a, "")
	return err
}
