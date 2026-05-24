package main

import (
	"encoding/json"
	"log"
	"notification-service/client"
	"notification-service/config"
	"notification-service/mail"
	"notification-service/message"
)

type Message struct {
	Message    string `json:"message"`
	CustomerId string `json:"customer_id"`
}

func main() {

	config := config.InitConfig()

	emailSender := mail.NewEmailSender(config)

	sqsReader, err := message.NewSQSClient(config)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	for {
		messages, err := sqsReader.ReceiveMessages()
		if err != nil {
			log.Printf("error %v", err)
		}

		for _, message := range messages {

			log.Printf("RAW SQS BODY: %s", *message.Body)

			var messageBody Message

			err := json.Unmarshal([]byte(*message.Body), &messageBody)
			if err != nil {
				log.Printf("message body parse edilemedi: %v", err)
			}
			customer, err := client.CustomerServiceRequest(messageBody.CustomerId)
			if err != nil {
				log.Printf("customer istegi atilamadi: %v", err)
				continue
			}
			err = emailSender.SendEMail(customer.Email, "You have a new order", *message.Body)

			if err != nil {
				log.Printf("error: %v", err)
				continue
			}

			err = sqsReader.DeleteMessage(message.ReceiptHandle)

			if err != nil {
				log.Printf("error: %v", err)
				continue
			}

		}
	}

}
