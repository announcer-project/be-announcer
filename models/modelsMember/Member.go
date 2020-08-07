package modelsMember

import (
	"github.com/jinzhu/gorm"
)

type Member struct {
	gorm.Model
	UserID   string
	SystemID string
	RoleID   uint

	MemberGroup []MemberGroup `gorm:"foreignkey:MemberID"`
	// News     []modelsNews.News     `gorm:"foreignkey:SystemID"`
	// NewsType []modelsNews.NewsType `gorm:"foreignkey:SystemID"`
}
