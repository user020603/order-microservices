package model

type Order struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    UserID    uint      `json:"user_id"`
    Status    string    `json:"status"`
    Total     float64   `json:"total"`
    Items     []Item    `json:"items" gorm:"foreignKey:OrderID"`
}

type Item struct {
    ID       uint    `json:"id" gorm:"primaryKey"`
    OrderID  uint    `json:"order_id"`
    Name     string  `json:"name"`
    Price    float64 `json:"price"`
    Quantity int     `json:"quantity"`
}