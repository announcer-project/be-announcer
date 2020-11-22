package modelsLineAPI

import "github.com/jinzhu/gorm"

type RichMenu struct {
	gorm.Model
	RichID   string `gorm:"unique;not null"`
	Status   string `gorm:"not null"`
	LineOAID uint   `gorm:"not null"`
}

func (RichMenu) TableName() string {
	return "richmenus"
}
