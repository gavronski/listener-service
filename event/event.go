package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"notes_topic", // name
		"topic",       // type
		true,          // durable?
		false,         // auto-delated?
		false,         // internal? no through micorservices
		false,         // no-wait?
		nil,           // arguments
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name?
		false, // durable?
		false, // delete when unused?
		true,  // exclusive?
		false,
		nil,
	)
}
