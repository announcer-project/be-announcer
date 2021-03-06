package modelsMember

import (
	"log"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
)

type Member struct {
	ID        string     `gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
	Approve   bool       `gorm:"not null" json:"approve"`
	LineID    string     `gorm:"not null" json:"line_id"`
	FName     string     `gorm:"not null" json:"f_name"`
	LName     string     `gorm:"not null" json:"l_name"`
	SystemID  string     `gorm:"not null" json:"system_id"`
	RoleID    uint       `gorm:"not null" json:"role_id"`

	MemberGroup      []MemberGroup      `gorm:"foreignkey:MemberID"`
	MemberInterested []MemberInterested `gorm:"foreignkey:MemberID"`
}

func (m *Member) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", GenUserID())
	return nil
}

func GenUserID() string {
	UserID := "MEMBER-"
	for i := 0; i < 16; i++ {
		ranType := rand.Intn(2)
		switch ranType {
		case 0:
			UserID += string(rand.Intn(57-48) + 48)
		case 1:
			UserID += string(rand.Intn(90-65) + 65)
		case 2:
			UserID += string(rand.Intn(122-97) + 97)
		}
	}
	log.Print(UserID)
	return UserID
}

func (m *Member) AddNewsTypeInterested(memberInterested MemberInterested) {
	m.MemberInterested = append(m.MemberInterested, memberInterested)
}
