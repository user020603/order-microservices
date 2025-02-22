package handler

import (
	"context"
	"encoding/json"
	"order/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type OrderHandler struct {
    db    *gorm.DB
    redis *redis.Client
}

func NewOrderHandler(db *gorm.DB, redis *redis.Client) *OrderHandler {
    return &OrderHandler{db: db, redis: redis}
}

func (h *OrderHandler) Create(c *gin.Context) {
    var order model.Order
    if err := c.ShouldBindJSON(&order); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    if err := h.db.Create(&order).Error; err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, order)
}

func (h *OrderHandler) Get(c *gin.Context) {
    id := c.Param("id")
    
    // Try to get from cache
    val, err := h.redis.Get(context.Background(), "order:"+id).Result()
    if err == nil {
        var order model.Order
        json.Unmarshal([]byte(val), &order)
        c.JSON(200, order)
        return
    }
    
    // Get from database
    var order model.Order
    if err := h.db.Preload("Items").First(&order, id).Error; err != nil {
        c.JSON(404, gin.H{"error": "Order not found"})
        return
    }
    
    // Store in cache
    orderJSON, _ := json.Marshal(order)
    h.redis.Set(context.Background(), "order:"+id, orderJSON, time.Hour)
    
    c.JSON(200, order)
}

func (h *OrderHandler) Update(c *gin.Context) {
    id := c.Param("id")
    var order model.Order
    
    if err := h.db.First(&order, id).Error; err != nil {
        c.JSON(404, gin.H{"error": "Order not found"})
        return
    }
    
    if err := c.ShouldBindJSON(&order); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    h.db.Save(&order)
    
    // Invalidate cache
    h.redis.Del(context.Background(), "order:"+id)
    
    c.JSON(200, order)
}

func (h *OrderHandler) Delete(c *gin.Context) {
    id := c.Param("id")
    
    if err := h.db.Delete(&model.Order{}, id).Error; err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    // Invalidate cache
    h.redis.Del(context.Background(), "order:"+id)
    
    c.JSON(200, gin.H{"message": "Order deleted"})
}