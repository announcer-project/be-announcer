package models

import (
	"be_nms/models/modelsNews"

	"github.com/jinzhu/gorm"
)

type Admin struct {
	gorm.Model
	Position string `gorm:"not null" json:"position"`
	SystemID string `gorm:"not null" json:"system_id"`
	UserID   string `gorm:"not null" json:"user_id"`
	System   System `gorm:"-" json:"system"`

	News []modelsNews.News `gorm:"foreignkey:AuthorID" json:"-"`
}
