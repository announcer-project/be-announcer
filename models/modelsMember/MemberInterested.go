package modelsMember

import (
	"github.com/jinzhu/gorm"
)

type MemberInterested struct {
	gorm.Model
	MemberID   string `gorm:"not null"`
	NewsTypeID uint   `gorm:"not null"`
}

func (MemberInterested) TableName() string {
	return "memberinteresteds"
}
