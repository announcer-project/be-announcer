package modelsMember

import (
	"github.com/jinzhu/gorm"
)

type MemberInterested struct {
	gorm.Model
	MemberID   string
	NewsTypeID uint
}

func (MemberInterested) TableName() string {
	return "memberinteresteds"
}
