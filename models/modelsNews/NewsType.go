package modelsNews

import (
	"be_nms/models/modelsMember"

	"github.com/jinzhu/gorm"
)

type NewsType struct {
	gorm.Model
	NewsTypeName string `gorm:"not null" json:"newstype_name"`
	SystemID     string `gorm:"not null" json:"system_id"`
	NumberNews   int    `gorm:"-" json:"number_news"`

	TypeOfNews       []TypeOfNews                    `gorm:"foreignkey:NewsTypeID" json:"-"`
	MemberInterested []modelsMember.MemberInterested `gorm:"foreignkey:NewsTypeID" json:"-"`
}

func (NewsType) TableName() string {
	return "newstypes"
}
