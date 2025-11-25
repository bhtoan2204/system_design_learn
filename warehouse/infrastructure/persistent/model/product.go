package model

import "time"

type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Quantity    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Product) TableName() string {
	return "products"
}
