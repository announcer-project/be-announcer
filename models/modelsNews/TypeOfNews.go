package modelsNews

import (
	"github.com/jinzhu/gorm"
)

type TypeOfNews struct {
	gorm.Model
	NewsID     uint
	NewsTypeID uint
}
