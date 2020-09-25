package modelsMember

import (
	"time"
)

type Member struct {
	ID        uint       `gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
	UserID    string     `json:"user_id"`
	SystemID  string     `json:"system_id"`
	RoleID    uint       `json:"role_id"`
	Approve   bool       `json:"approve"`

	MemberGroup      []MemberGroup      `gorm:"foreignkey:MemberID"`
	MemberInterested []MemberInterested `gorm:"foreignkey:MemberID"`
	// News     []modelsNews.News     `gorm:"foreignkey:SystemID"`
	// NewsType []modelsNews.NewsType `gorm:"foreignkey:SystemID"`
}

func (m *Member) AddNewsTypeInterested(memberInterested MemberInterested) {
	m.MemberInterested = append(m.MemberInterested, memberInterested)
}
