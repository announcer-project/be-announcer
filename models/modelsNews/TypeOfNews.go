package modelsNews

import (
	"github.com/jinzhu/gorm"
)

type TypeOfNews struct {
	gorm.Model
	NewsID     uint `gorm:"not null"`
	NewsTypeID uint `gorm:"not null"`

	NewsTypeName string `gorm:"-"`
}

func (TypeOfNews) TableName() string {
	return "typeofnews"
}
