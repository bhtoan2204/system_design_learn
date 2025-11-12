package model

import "time"

type User struct {
	ID        string    `gorm:"primaryKey;type:varchar(36);not null" json:"id"`
	Username  string    `gorm:"type:varchar(50);not null" json:"username"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"password"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null" json:"updated_at"`
}

func (u *User) TableName() string {
	return "users"
}
