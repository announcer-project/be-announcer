package modelsLineAPI

import "github.com/jinzhu/gorm"

type RichMenu struct {
	gorm.Model
	RichID   string
	Status   string
	LineOAID uint
}

func (RichMenu) TableName() string {
	return "richmenus"
}
