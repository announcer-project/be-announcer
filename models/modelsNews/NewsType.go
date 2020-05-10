package modelsNews

import (
	"github.com/jinzhu/gorm"
)

type NewsType struct {
	gorm.Model
	NewTypeName string
	SystemID    uint

	TypeOfNews []TypeOfNews `gorm:"foreignkey:NewsTypeID"`
}

func (NewsType) TableName() string {
	return "NewsTypes"
}
