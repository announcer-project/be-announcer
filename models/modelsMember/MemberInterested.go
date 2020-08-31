package modelsMember

import (
	"github.com/jinzhu/gorm"
)

type MemberInterested struct {
	gorm.Model
	MemberID   uint
	NewsTypeID uint

	// Admin    []Admin               `gorm:"foreignkey:SystemID"`
	// News     []modelsNews.News     `gorm:"foreignkey:SystemID"`
	// NewsType []modelsNews.NewsType `gorm:"foreignkey:SystemID"`
}

func (MemberInterested) TableName() string {
	return "memberinteresteds"
}
