package models

type User struct {
	ID      uint    `json:"ID"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
	Age     uint    `json:"age"`
}
