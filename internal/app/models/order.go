package models

type Order struct {
	ID        int     `json:"id"`
	UserID    string  `json:"user_id"`
	OrderID   string  `json:"number"`
	Accrual   float32 `json:"accrual"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
}
