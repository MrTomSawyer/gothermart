package models

type Order struct {
	Number    int     `json:"number"`
	UserID    string  `json:"user_id"`
	OrderID   string  `json:"order_id"`
	Accrual   float32 `json:"accrual"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
}
