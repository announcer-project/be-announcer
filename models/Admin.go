package models

import (
	"be_nms/models/modelsNews"

	"github.com/jinzhu/gorm"
)

type Admin struct {
	gorm.Model
	UserID   string `json:"user_id"`
	SystemID string `json:"system_id"`
	Position string `json:"position"`
	System   System `gorm:"-" json:"system"`

	News         []modelsNews.News         `gorm:"foreignkey:AuthorID"`
	Announcement []modelsNews.Announcement `gorm:"foreignkey:AdminID"`
}
