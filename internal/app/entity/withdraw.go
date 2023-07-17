package entity

type Withdrawal struct {
	ID          int
	UserId      int
	OrderID     string
	Sum         float32
	ProcessedAt string
}
