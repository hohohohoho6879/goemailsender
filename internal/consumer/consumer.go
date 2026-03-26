package consumer

import (
	"context"
	"emailsender/internal/config"
	"emailsender/internal/mailer"
	"emailsender/internal/template"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	config config.Config
	mailer *mailer.Mailer
}

func New(configuration config.Config, mailSender *mailer.Mailer) *Consumer {
	return &Consumer{config: configuration, mailer: mailSender}
}

func (consumer *Consumer) Run(context context.Context) {
	for {
		select {
		case <-context.Done():
			log.Println("Consumer stopped")
			return
		default:
			if error := consumer.consume(context); error != nil {
				log.Printf("Consumer error: %s, reconnecting in 5s...", error)
				select {
				case <-time.After(5 * time.Second):
				case <-context.Done():
					return
				}
			}
		}
	}
}

func (consumer *Consumer) consume(context context.Context) error {
	connection, error := amqp.Dial(consumer.config.RabbitMQURL)
	if error != nil {
		return error
	}
	defer connection.Close()

	channel, error := connection.Channel()
	if error != nil {
		return error
	}
	defer channel.Close()

	_, error = channel.QueueDeclare("mail", true, false, false, false, nil)
	if error != nil {
		return error
	}

	channel.Qos(1, 0, false)

	deliveries, error := channel.Consume("mail", "", false, false, false, false, nil)
	if error != nil {
		return error
	}

	connectionClosed := make(chan *amqp.Error, 1)
	connection.NotifyClose(connectionClosed)

	log.Println("Waiting for messages...")

	for {
		select {
		case <-context.Done():
			return nil
		case amqpError := <-connectionClosed:
			return fmt.Errorf("connection closed: %w", amqpError)
		case delivery, ok := <-deliveries:
			if !ok {
				return fmt.Errorf("channel closed")
			}
			consumer.handleMessage(delivery)
		}
	}
}

func (consumer *Consumer) handleMessage(delivery amqp.Delivery) {
	var emailData template.EmailData
	if error := json.Unmarshal(delivery.Body, &emailData); error != nil {
		log.Printf("Failed to parse message: %s, body: %s", error, string(delivery.Body))
		delivery.Nack(false, false)
		return
	}

	if emailData.To == "" {
		log.Printf("Skipping message with empty 'to' field, body: %s", string(delivery.Body))
		delivery.Ack(false)
		return
	}

	html, error := template.Render(emailData)
	if error != nil {
		log.Printf("Failed to render template: %s", error)
		delivery.Nack(false, false)
		return
	}

	if error := consumer.mailer.Send(emailData.To, emailData.Subject, html); error != nil {
		log.Printf("Failed to send email to %s: %s", emailData.To, error)
		delivery.Nack(false, true)
		return
	}

	log.Printf("Email sent to %s", emailData.To)
	delivery.Ack(false)
}
