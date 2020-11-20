package modelsMember

import "time"

type TargetGroup struct {
	ID              uint       `gorm:"primary_key"`
	CreatedAt       time.Time  `json:"-"`
	UpdatedAt       time.Time  `json:"-"`
	DeletedAt       *time.Time `sql:"index" json:"-"`
	TargetGroupName string     `json:"targetgroup_name"`
	NumberOfMembers int        `json:"number_members"`
	SystemID        string     `json:"system_id"`

	MemberGroup []MemberGroup `gorm:"foreignkey:TargetGroupID" json:"member_group"`
}

func (TargetGroup) TableName() string {
	return "targetgroups"
}

func (group *TargetGroup) AddMemberGroup(member MemberGroup) {
	group.MemberGroup = append(group.MemberGroup, member)
}
