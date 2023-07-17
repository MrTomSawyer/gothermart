package models

type User struct {
	Id        int     `json:"id"`
	Login     string  `json:"login"`
	Password  string  `json:"password"`
	Balance   float32 `json:"balance"`
	Withdrawn float32 `json:"withdrawn"`
}
