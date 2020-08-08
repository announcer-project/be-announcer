package modelsMember

import (
	"github.com/jinzhu/gorm"
)

type TargetGroup struct {
	gorm.Model
	TargetGroupName string `json:"targetgroup_name"`
	NumberOfMembers int    `json:"number_members"`
	SystemID        string

	MemberGroup []MemberGroup `gorm:"foreignkey:TargetGroupID"`
	// News     []modelsNews.News     `gorm:"foreignkey:SystemID"`
	// NewsType []modelsNews.NewsType `gorm:"foreignkey:SystemID"`
}

func (TargetGroup) TableName() string {
	return "targetgroups"
}
