package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

type QueueMessage struct {
	OrderID string `json:"order_id"`
}

func main() {
	// Determine environment
	env := os.Getenv("ENV")
	if env != "prod" {
		log.Fatal("ENV is not prod")
	}

	// Load configuration
	viper.SetConfigFile("consumer.toml")
	viper.SetConfigType("toml")
	viper.SetEnvPrefix("puzzles")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	time.Sleep(10 * time.Second)

	// Connect to rabbitmq
	fmt.Printf("Connecting to rabbitmq with: %s:%s@%s:%d\n", viper.GetString("rabbitmq.user"), viper.GetString("rabbitmq.pass"), viper.GetString("rabbitmq.host"), viper.GetInt("rabbitmq.port"))
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", viper.GetString("rabbitmq.user"), viper.GetString("rabbitmq.pass"), viper.GetString("rabbitmq.host"), viper.GetInt("rabbitmq.port")))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		viper.GetString("rabbitmq.queue"),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for msg := range msgs {
			var qm QueueMessage
			if err := json.Unmarshal(msg.Body, &qm); err != nil {
				fmt.Printf("Error decoding message: %v", err)
				continue
			}

			log.Printf("Received order: %s", qm.OrderID)
		}
	}()
}
