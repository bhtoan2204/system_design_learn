package persistent

import "time"

type Tenant struct {
	ID        int       `gorm:"primaryKey,autoIncrement"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Tenant) TableName() string {
	return "tenants"
}
