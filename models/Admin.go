package models

import (
	"be_nms/models/modelsNews"

	"github.com/jinzhu/gorm"
)

type Admin struct {
	gorm.Model
	UserID   string
	SystemID string
	Position string

	News         []modelsNews.News         `gorm:"foreignkey:AuthorID"`
	Announcement []modelsNews.Announcement `gorm:"foreignkey:AdminID"`
}
