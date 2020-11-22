package modelsMember

import (
	"github.com/jinzhu/gorm"
)

type TargetGroup struct {
	gorm.Model
	TargetGroupName string `gorm:"not null" json:"targetgroup_name"`
	NumberOfMembers int    `gorm:"not null" json:"number_members"`
	SystemID        string `gorm:"not null" json:"system_id"`

	MemberGroup []MemberGroup `gorm:"foreignkey:TargetGroupID" json:"member_group"`
}

func (TargetGroup) TableName() string {
	return "targetgroups"
}

func (group *TargetGroup) AddMemberGroup(member MemberGroup) {
	group.MemberGroup = append(group.MemberGroup, member)
}
