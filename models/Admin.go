package models

import (
	"be_nms/models/modelsNews"
	"time"
)

type Admin struct {
	ID        uint       `gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
	Position  string     `json:"position"`
	System    System     `gorm:"-" json:"system"`
	UserID    string     `json:"user_id"`
	SystemID  string     `json:"system_id"`

	News []modelsNews.News `gorm:"foreignkey:AuthorID" json:"-"`
}
