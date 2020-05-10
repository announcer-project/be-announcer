package modelsMember

import (
	"github.com/jinzhu/gorm"
)

type TargetGroup struct {
	gorm.Model
	TargetGroupName string
	NumberOfMembers int
	SystemID        uint

	MemberGroup []MemberGroup `gorm:"foreignkey:TargetGroupID"`
	// News     []modelsNews.News     `gorm:"foreignkey:SystemID"`
	// NewsType []modelsNews.NewsType `gorm:"foreignkey:SystemID"`
}