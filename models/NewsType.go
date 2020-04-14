package models

import "time"

type NewsType struct {
	NewsTypeID uint	`gorm:"primary_key;"`
	NewsTypeName  string
	SystemID uint 
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
	TypeOfNews []TypeOfNews `gorm:"foreignkey:NewsTypeID"`
}

func (n *NewsType) CreateNewsType(TypeName string, SystemID uint) {
	n.NewsTypeName = TypeName
	n.SystemID = SystemID
}