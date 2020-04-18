package models

import (
	"be_nms/models/modelsNews"

	"github.com/jinzhu/gorm"
)

type System struct {
	gorm.Model
	SystemName	string
	OwnerID		string

	Admin	[]Admin	`gorm:"foreignkey:SystemID"`
	News []modelsNews.News `gorm:"foreignkey:SystemID"`
}