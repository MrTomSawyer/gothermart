package entity

type Withdrawal struct {
	ID          int
	UserID      int
	OrderID     string
	Sum         float32
	ProcessedAt string
}
