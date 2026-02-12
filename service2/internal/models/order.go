package models

type Order struct {
	ID     uint    `json:"id"`
	UserID uint    `json:"user_id"`
	Item   string  `json:"item"`
	Price  float64 `json:"price"`
	User   *User   `json:"-"`
}

type User struct {
	ID      uint    `json:"id"`
	Name    string  `json:"name"`
	Age     uint    `json:"age"`
	Balance float64 `json:"balance"`
}
