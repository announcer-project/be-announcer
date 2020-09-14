package modelsNews

import (
	"be_nms/models/modelsMember"
	"time"
)

type NewsType struct {
	ID           uint       `gorm:"primary_key"`
	CreatedAt    time.Time  `json:"-"`
	UpdatedAt    time.Time  `json:"-"`
	DeletedAt    *time.Time `sql:"index" json:"-"`
	NewsTypeName string     `json:"newstype_name"`
	SystemID     string     `json:"system_id"`
	NumberNews   int        `gorm:"-" json:"number_news"`

	TypeOfNews       []TypeOfNews                    `gorm:"foreignkey:NewsTypeID" json:"-"`
	MemberInterested []modelsMember.MemberInterested `gorm:"foreignkey:NewsTypeID" json:"-"`
}

func (NewsType) TableName() string {
	return "newstypes"
}
