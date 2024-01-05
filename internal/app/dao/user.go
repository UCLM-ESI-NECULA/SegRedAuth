package dao

import (
	"time"
)

type BaseModel struct {
	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
}

type User struct {
	ID       int    `gorm:"column:id; primary_key; not null" json:"-"`
	Username string `gorm:"unique;column:username; not null" json:"username"`
	Password string `gorm:"column:password; not null" json:"password,omitempty"`
	BaseModel
}

type Token struct {
	Token string `json:"access_token"`
}
