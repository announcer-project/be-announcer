package modelsMember

import (
	"github.com/jinzhu/gorm"
)

type Member struct {
	gorm.Model
	UserID   string
	SystemID string
	RoleID   uint

	MemberGroup      []MemberGroup      `gorm:"foreignkey:MemberID"`
	MemberInterested []MemberInterested `gorm:"foreignkey:MemberID"`
	// News     []modelsNews.News     `gorm:"foreignkey:SystemID"`
	// NewsType []modelsNews.NewsType `gorm:"foreignkey:SystemID"`
}

func (m *Member) AddNewsTypeInterested(memberInterested MemberInterested) {
	m.MemberInterested = append(m.MemberInterested, memberInterested)
}
