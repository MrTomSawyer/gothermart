package entity

type Order struct {
	ID        int
	UserID    string
	OrderID   string
	Accrual   float32
	Status    string
	CreatedAt string
}
