package modelsMember

import (
	"github.com/jinzhu/gorm"
)

type MemberGroup struct {
	gorm.Model
	MemberID      uint
	TargetGroupID uint
}

func (MemberGroup) TableName() string {
	return "membergroups"
}
