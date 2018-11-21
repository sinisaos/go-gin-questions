package models

import (
	"time"
)

type User struct {
	Id        int    `gorm:"primary_key" json:"id"`
	Username  string `gorm:"size:100;unique"`
	Email     string `gorm:"size:100"`
	CreatedAt time.Time
	UpdatedAt time.Time
	IsAdmin   bool
	Password  []byte
}
