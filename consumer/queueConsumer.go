package consumer

import (
	"../lib"
	"../model"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type QueueConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func NewQueueConsumer(queueName string, bindingKey string, ctag string) (*QueueConsumer, error) {
	c := &QueueConsumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}

	var err error
	uri := lib.GetAmqpUri()
	c.conn, err = amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	go func() {
		log.Println("closing:", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	c.channel, err = lib.GetChannel()
	if err != nil {
		return nil, err
	}

	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}

	err = c.channel.QueueBind(
		queue.Name, // name of the queue
		bindingKey,
		lib.GetExchange(), // sourceExchange
		false,             // noWait
		nil,               // arguments
	)
	if err != nil {
		return nil, err
	}

	deliveries, err := c.channel.Consume(
		queue.Name, // name
		c.tag,      // consumerTag,
		true,       // noAck
		true,       // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, err
	}

	go handle(deliveries, c.done)

	return c, nil
}

func (c *QueueConsumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		var message model.Message
		err := json.Unmarshal(d.Body, &message)
		if err != nil {
			log.Println("Invalid data received:", err)
		} else {
			log.Println("Invite message", message)
			message.ProcessAndDeliver()
		}
	}
	log.Println("handle: deliveries channel closed")
	done <- nil
}
