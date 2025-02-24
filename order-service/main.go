package main

import (
	"context"
	"log"
	"order/config"
	"order/events"
	"order/handler"
	"order/model"

	"github.com/gin-gonic/gin"
)

func main() {
    db, err := config.SetupDB()
    if err != nil {
        panic(err)
    }
    
    redis := config.SetupRedis()
    kafkaClient, err := events.NewKafkaClient()
    if err != nil {
        log.Fatalf("Failed to create Kafka client: %v", err)
    }

    ctx := context.Background()
    kafkaClient.ConsumeUserEvents(ctx)
    
    // Auto migrate the schema
    db.AutoMigrate(&model.Order{}, &model.Item{})
    
    router := gin.Default()
    orderHandler := handler.NewOrderHandler(db, redis, kafkaClient)
    
    orders := router.Group("/api/orders")
    {
        orders.POST("/", orderHandler.Create)
        orders.GET("/:id", orderHandler.Get)
        orders.PUT("/:id", orderHandler.Update)
        orders.DELETE("/:id", orderHandler.Delete)
    }
    
    router.Run(":8081")
}