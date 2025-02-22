package config

import (
    "fmt"
    "os"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func SetupDB() (*gorm.DB, error) {
    host := os.Getenv("DB_HOST")
    if host == "" {
        host = "localhost"
    }
    
    user := os.Getenv("DB_USER")
    if user == "" {
        user = "root"
    }
    
    password := os.Getenv("DB_PASSWORD")
    if password == "" {
        password = "password"
    }
    
    dbname := os.Getenv("DB_NAME")
    if dbname == "" {
        dbname = "ordersdb"
    }

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
        host,
        user,
        password,
        dbname,
    )
    
    return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// curl -X POST http://localhost:8081/api/orders -H "Content-Type: application/json" -d '{"user_id":1,"status":"pending","total":99.99,"items":[{"name":"Item 1","price":49.99,"quantity":2}]}'
// curl http://localhost:8081/api/orders/1