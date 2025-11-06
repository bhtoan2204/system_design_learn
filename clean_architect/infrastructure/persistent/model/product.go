package model

import "time"

type Product struct {
	ID          string    `json:"id"`
	Sku         string    `json:"sku"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	SellerID    string    `json:"seller_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
