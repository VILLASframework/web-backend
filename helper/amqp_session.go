package helper

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type AMQPsession struct {
	connection          *amqp.Connection
	sendCh              *amqp.Channel
	recvCh              *amqp.Channel
	notifyConnClose     chan *amqp.Error
	notifySendChanClose chan *amqp.Error
	notifySendConfirm   chan amqp.Confirmation
	notifyRecvChanClose chan *amqp.Error
	notifyRecvConfirm   chan amqp.Confirmation
	IsReady             bool
	done                chan bool
	name                string
	exchange            string
	processMessage      func(delivery amqp.Delivery)
}

const (
	// When reconnecting to the server after connection failure
	reconnectDelay = 5 * time.Second

	// When setting up the channel after a channel exception
	reInitDelay = 2 * time.Second

	// When resending messages the server didn't confirm
	resendDelay = 5 * time.Second
)

//var client AMQPsession

// NewAMQPSession creates a new consumer state instance, and automatically
// attempts to connect to the server.
func NewAMQPSession(name string, AMQPurl string, exchange string, processMessage func(delivery amqp.Delivery)) *AMQPsession {
	if AMQPurl != "" {
		log.Println("Starting AMQP client")

		session := AMQPsession{
			name:           name,
			exchange:       exchange,
			done:           make(chan bool),
			processMessage: processMessage,
		}
		go session.handleReconnect(AMQPurl)

		return &session
	}

	return nil
}

// handleReconnect will wait for a connection error on
// notifyConnClose, and then continuously attempt to reconnect.
func (session *AMQPsession) handleReconnect(addr string) {
	for {
		session.IsReady = false
		log.Println("Attempting to connect to AMQP broker ", addr)

		conn, err := session.connect(addr)

		if err != nil {
			log.Println("Failed to connect. Retrying...")

			select {
			case <-session.done:
				return
			case <-time.After(reconnectDelay):
			}
			continue
		}

		if done := session.handleReInit(conn); done {
			break
		}
	}
}

// connect will create a new AMQP connection
func (session *AMQPsession) connect(addr string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr)

	if err != nil {
		return nil, err
	}

	// take a new connection to the queue, and updates the close listener to reflect this.
	session.connection = conn
	session.notifyConnClose = make(chan *amqp.Error)
	session.connection.NotifyClose(session.notifyConnClose)

	log.Println("Connected!")
	return conn, nil
}

// handleReInit will wait for a channel error
// and then continuously attempt to re-initialize both channels
func (session *AMQPsession) handleReInit(conn *amqp.Connection) bool {
	for {
		session.IsReady = false

		err := session.init(conn)

		if err != nil {
			log.Println("Failed to initialize channel. Retrying...")

			select {
			case <-session.done:
				return true
			case <-time.After(reInitDelay):
			}
			continue
		}

		select {
		case <-session.done:
			return true
		case <-session.notifyConnClose:
			log.Println("Connection closed. Reconnecting...")
			return false
		case <-session.notifySendChanClose:
			log.Println("Send channel closed. Re-running init...")
		case <-session.notifyRecvChanClose:
			log.Println("Receive channel closed. Re-running init...")
		}
	}
}

// init will initialize channel & declare queue
func (session *AMQPsession) init(conn *amqp.Connection) error {

	// create sendCh
	sendCh, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("AMQP: failed to open a sendCh, error: %v", err)
	}
	// declare exchange
	err = sendCh.ExchangeDeclare(session.exchange,
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
	ICQueue, err := sendCh.QueueDeclare("",
		false,
		true,
		true,
		false,
		nil)
	if err != nil {
		return fmt.Errorf("AMQP: failed to declare the queue, error: %v", err)
	}

	err = sendCh.QueueBind(ICQueue.Name, "", session.exchange, false, nil)
	if err != nil {
		return fmt.Errorf("AMQP: failed to bind the queue, error: %v", err)
	}

	session.sendCh = sendCh
	session.notifySendChanClose = make(chan *amqp.Error)
	session.notifySendConfirm = make(chan amqp.Confirmation, 1)
	session.sendCh.NotifyClose(session.notifySendChanClose)
	session.sendCh.NotifyPublish(session.notifySendConfirm)

	// create receive channel
	recvCh, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("AMQP: failed to open a recvCh, error: %v", err)
	}

	session.recvCh = recvCh
	session.notifyRecvChanClose = make(chan *amqp.Error)
	session.notifyRecvConfirm = make(chan amqp.Confirmation, 1)
	session.recvCh.NotifyClose(session.notifyRecvChanClose)
	session.recvCh.NotifyPublish(session.notifyRecvConfirm)

	// start deliveries
	messages, err := session.recvCh.Consume(ICQueue.Name,
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
				session.processMessage(message)
			}
		}
	}()
	session.IsReady = true
	log.Println("AMQP channels setup! Waiting for messages...")

	return nil
}

func (session *AMQPsession) CheckConnection() error {

	if session.connection != nil {
		if session.connection.IsClosed() {
			return fmt.Errorf("connection to broker is closed")
		}
	} else {
		return fmt.Errorf("connection is nil")
	}

	return nil
}

func (session *AMQPsession) Send(payload []byte, destinationUuid string) error {

	message := amqp.Publishing{
		DeliveryMode:    2,
		Timestamp:       time.Now(),
		ContentType:     "application/json",
		ContentEncoding: "utf-8",
		Priority:        0,
		Body:            payload,
	}

	// set message headers
	message.Headers = make(map[string]interface{}) // empty map
	message.Headers["uuid"] = destinationUuid      // leave uuid empty if message should go to all ICs

	err := session.CheckConnection()
	if err != nil {
		return err
	}

	err = session.sendCh.Publish(session.exchange,
		"",
		false,
		false,
		message)
	return err
}
