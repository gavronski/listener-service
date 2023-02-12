package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	Conn      *amqp.Connection
	queueName string
}

type Note struct {
	ID              int       `json:"note_id,string,omitempty"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	TextColor       string    `json:"text_color"`
	BackgroundColor string    `json:"background_color"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		Conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// open up channel and declare an exchange
func (consumer *Consumer) setup() error {
	channel, err := consumer.Conn.Channel()

	if err != nil {
		return err
	}

	return declareExchange(channel)
}

// pushing event to rabbitmq
type Payload struct {
	Action string `json:"action"`
	Data   Note   `json:"data,omitempty"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"notes_topic",
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)

	if err != nil {
		return err
	}

	// consume messages from rabbit 4ever

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()
	fmt.Printf("Waiting for message [Exchange, Queue] [notes_topic, %s]\n", q.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload) {

	switch payload.Action {
	case "add-note":
		err := sendToNoteService(payload, "POST", "http://notes-service/add")

		if err != nil {
			if err != nil {
				log.Println(err)
			}
		}

	case "update-note":
		err := sendToNoteService(payload, "PATCH", "http://notes-service/update")

		if err != nil {
			if err != nil {
				log.Println(err)
			}
		}

	case "delete-note":
		err := sendToNoteService(payload, "DELETE", "http://notes-service/delete")

		if err != nil {
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func sendToNoteService(entry Payload, method string, url string) error {
	jsonData, _ := json.Marshal(entry)
	log.Println(entry.Data)
	request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	log.Println(response)
	return nil
}
