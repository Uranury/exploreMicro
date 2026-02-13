package models

type User struct {
	ID      uint    `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
	Age     uint    `json:"age"`
}
