package main

import (
	"order/config"
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
    
    // Auto migrate the schema
    db.AutoMigrate(&model.Order{}, &model.Item{})
    
    router := gin.Default()
    orderHandler := handler.NewOrderHandler(db, redis)
    
    orders := router.Group("/api/orders")
    {
        orders.POST("/", orderHandler.Create)
        orders.GET("/:id", orderHandler.Get)
        orders.PUT("/:id", orderHandler.Update)
        orders.DELETE("/:id", orderHandler.Delete)
    }
    
    router.Run(":8081")
}