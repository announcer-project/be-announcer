package modelsMember

import (
	"github.com/jinzhu/gorm"
)

type MemberGroup struct {
	gorm.Model
	MemberID      uint
	TargetGroupID uint

	// Admin    []Admin               `gorm:"foreignkey:SystemID"`
	// News     []modelsNews.News     `gorm:"foreignkey:SystemID"`
	// NewsType []modelsNews.NewsType `gorm:"foreignkey:SystemID"`
}
