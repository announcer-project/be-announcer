package modelsNews

import (
	"github.com/jinzhu/gorm"
)

type TypeOfNews struct {
	gorm.Model
	NewsID     uint
	NewsTypeID uint

	NewsTypeName string `gorm:"-"`
}

func (TypeOfNews) TableName() string {
	return "TypeOfNews"
}
