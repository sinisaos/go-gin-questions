package models

type Tag struct {
	Id        int `gorm:"primary_key" json:"id"`
	Name      string
	Questions []Question `gorm:"many2many:taggings;"`
}
