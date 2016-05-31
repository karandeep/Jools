package lib

import (
	"../config"
	"github.com/streadway/amqp"
	"log"
	"strconv"
)

var (
	queueConnection amqp.Connection
	queueChannel    *amqp.Channel
	channelCreated  bool   = false
	exchangeName    string = "message-exchange"
	exchangeType    string = "direct"
)

func GetAmqpUri() string {
	configData := config.GetConfig()
	uri := "amqp://" + configData.RABBITMQ_USER + ":" + configData.RABBITMQ_PASSWORD + "@" + configData.RABBITMQ_HOST + ":5672"
	return uri
}

func GetExchange() string {
	return exchangeName
}

func GetChannel() (*amqp.Channel, error) {
	if !channelCreated {
		uri := GetAmqpUri()
		queueConnection, err := amqp.Dial(uri)
		if err != nil {
			log.Println("Rabbit MQ Connection failure:", err)
			return queueChannel, err
		}
		queueChannel, err = queueConnection.Channel()
		if err != nil {
			log.Println("Rabbit MQ Connection failure:", err)
			return queueChannel, err
		}

		if err := queueChannel.ExchangeDeclare(
			exchangeName, // name
			exchangeType, // type
			true,         // durable
			false,        // auto-deleted
			false,        // internal
			false,        // noWait
			nil,          // arguments
		); err != nil {
			return queueChannel, err
		}

		channelCreated = true
	}
	return queueChannel, nil
}

func Enqueue(queueType int, message []byte) {
	channel, err := GetChannel()
	if err != nil {
		return
	}
	routingKey := strconv.Itoa(queueType)
	if err = channel.Publish(
		exchangeName, // publish to an exchange
		routingKey,   // routing to 0 or more queues
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            message,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		log.Println("Message publish failed:", queueType, err)
		return
	}
}
