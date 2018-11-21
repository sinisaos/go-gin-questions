package models

import (
	"time"
)

type Question struct {
	Id             int `gorm:"primary_key" json:"id"`
	Title          string
	Body           string `sql:"type:text;" json:"body"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Views          int
	Likes          int
	AnswerCount    int
	AcceptedAnswer bool
	UserID         int `gorm:"size:10"`
	User           User
	Tags           []Tag `gorm:"many2many:taggings;"`
}
