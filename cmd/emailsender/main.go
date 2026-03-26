package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"emailsender/internal/config"
	"emailsender/internal/consumer"
	"emailsender/internal/mailer"
)

func main() {
	configuration := config.Load()
	mailSender := mailer.New(configuration)
	messageConsumer := consumer.New(configuration, mailSender)

	signalContext, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Println("Starting email sender service...")
	messageConsumer.Run(signalContext)
}
