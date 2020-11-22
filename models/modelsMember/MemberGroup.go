package modelsMember

import (
	"github.com/jinzhu/gorm"
)

type MemberGroup struct {
	gorm.Model
	MemberID      string `gorm:"not null"`
	TargetGroupID uint   `gorm:"not null"`
}

func (MemberGroup) TableName() string {
	return "membergroups"
}
