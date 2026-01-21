package services

import (
	"encoding/json"
	"fmt"
	"os"

	"api-budgeting.smartcodex.cloud/models"

	"github.com/rabbitmq/amqp091-go"
)

func CreateQueue(queue models.QueuePush) map[string]string {

	username := os.Getenv("AMQP_USER")
	password := os.Getenv("AMQP_PASSWORD")
	host := os.Getenv("AMQP_HOST")
	port := os.Getenv("AMQP_PORT")

	connectAmqp := fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)

	conn, err := amqp091.Dial(connectAmqp)
	if err != nil {
		return map[string]string{
			"error": "Failed to connect to message broker",
		}
	}

	ch, errCh := conn.Channel()
	if errCh != nil {
		return map[string]string{
			"error": "Failed to open channel.",
		}
	}

	defer ch.Close()

	q, _ := ch.QueueDeclare("go.webhook", true, false, false, false, nil)

	jsonEncode, errJsonEncode := json.Marshal(&queue)
	if errJsonEncode != nil {
		return map[string]string{
			"error": "Something went wrong (E500C)",
		}
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        jsonEncode,
		},
	)

	if err != nil {
		return map[string]string{
			"error": "Failed publish to message broker.",
		}
	}

	return map[string]string{
		"error": "",
	}
}
