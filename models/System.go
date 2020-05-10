package models

import (
	"be_nms/models/modelsMember"
	"be_nms/models/modelsNews"

	"github.com/jinzhu/gorm"
)

type System struct {
	gorm.Model
	SystemName string
	OwnerID    string

	Admin       []Admin                    `gorm:"foreignkey:SystemID"`
	News        []modelsNews.News          `gorm:"foreignkey:SystemID"`
	NewsType    []modelsNews.NewsType      `gorm:"foreignkey:SystemID"`
	Member      []modelsMember.Member      `gorm:"foreignkey:SystemID"`
	TargetGroup []modelsMember.TargetGroup `gorm:"foreignkey:SystemID"`
}
