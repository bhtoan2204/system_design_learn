package model

import "time"

type Stock struct {
	ID        string
	ProductID string
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Stock) TableName() string {
	return "stocks"
}
