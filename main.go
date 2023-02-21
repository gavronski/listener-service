package main

import (
	"fmt"
	"listener-service/event"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// connect to rabbitmq
	rabbitConn, err := connect()

	if err != nil {
		log.Panic(err)
	}

	defer rabbitConn.Close()
	// start listening for messages
	log.Println("Listening an consuming messages")

	// create consumer
	consumer, err := event.NewConsumer(rabbitConn)

	if err != nil {
		panic(err)
	}

	err = consumer.Listen([]string{"note-message"})
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	var attempts int64
	var delay = 1 * time.Second
	var connection *amqp.Connection

	for {
		// create new connection
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq")

		if err != nil {
			fmt.Println("RabbitMQ not yet ready")
			attempts++
		} else {
			// if connection is ready break the loop
			connection = conn
			break
		}

		// if number of connection attempts is bigger than 5 return error
		if attempts > 5 {
			fmt.Println(err)
			return nil, err
		}

		delay = time.Duration(math.Pow(float64(attempts), 2)) * time.Second
		log.Println("Delay attempt...")

		time.Sleep(delay)
		continue
	}

	return connection, nil
}
