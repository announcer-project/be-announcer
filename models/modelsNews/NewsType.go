package modelsNews

import (
	"github.com/jinzhu/gorm"
)

type NewsType struct {
	gorm.Model
	NewsTypeName string `json:"newstype_name"`
	SystemID     string
	NumberNews   int `gorm:"-" json:"number_news"`

	TypeOfNews []TypeOfNews `gorm:"foreignkey:NewsTypeID"`
}

func (NewsType) TableName() string {
	return "newstypes"
}
