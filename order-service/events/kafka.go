package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
)

type KafkaClient struct {
    producer sarama.SyncProducer
    consumer sarama.Consumer
}

type OrderEvent struct {
    EventType  string  `json:"eventType"`
    OrderID    uint    `json:"orderId"`
    UserID     uint    `json:"userId"`
    TotalPrice float64 `json:"totalPrice"`
    Status     string  `json:"status"`
}

type UserEvent struct {
    EventType string `json:"eventType"`
    UserID    int64  `json:"userId"`
    Username  string `json:"username"`
    Email     string `json:"email"`
}

func NewKafkaClient() (*KafkaClient, error) {
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    
    // Get Kafka brokers from environment variable
    brokerList := getKafkaBrokers()
    
    producer, err := sarama.NewSyncProducer(brokerList, config)
    if err != nil {
        return nil, err
    }
    
    consumer, err := sarama.NewConsumer(brokerList, config)
    if err != nil {
        return nil, err
    }
    
    return &KafkaClient{
        producer: producer,
        consumer: consumer,
    }, nil
}

func getKafkaBrokers() []string {
    brokersEnv := os.Getenv("KAFKA_BROKERS")
    if brokersEnv == "" {
        // Default to localhost if not set
        return []string{"localhost:9092"}
    }
    return strings.Split(brokersEnv, ",")
}

func (k *KafkaClient) SendOrderEvent(event OrderEvent) error {
    eventJSON, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    msg := &sarama.ProducerMessage{
        Topic: "order-events",
        Value: sarama.StringEncoder(eventJSON),
    }
    
    _, _, err = k.producer.SendMessage(msg)
    return err
}

func (k *KafkaClient) ConsumeUserEvents(ctx context.Context) {
    consumer, err := k.consumer.ConsumePartition("user-events", 0, sarama.OffsetNewest)
    if err != nil {
        log.Printf("Error creating consumer: %v", err)
        return
    }
    
    go func() {
        for {
            select {
            case msg := <-consumer.Messages():
                var userEvent UserEvent
                if err := json.Unmarshal(msg.Value, &userEvent); err != nil {
                    log.Printf("Error unmarshaling user event: %v", err)
                    continue
                }
                
                // Handle user event
                handleUserEvent(userEvent)
                
            case <-ctx.Done():
                return
            }
        }
    }()
}

func handleUserEvent(event UserEvent) {
    // Handle different user events
    switch event.EventType {
    case "CREATED":
        // Maybe create order for this user
        log.Printf("User %d created, handling their orders...", event.UserID)
		fmt.Printf("User %d created, handling their orders...", event.UserID)
    case "UPDATED":
        // Update user information in orders
        log.Printf("User %d updated, updating order records...", event.UserID)
    }
}