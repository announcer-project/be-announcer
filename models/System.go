package models

import (
	"be_nms/models/modelsMember"
	"be_nms/models/modelsNews"
	"log"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
)

type System struct {
	ID         string     `gorm:"primary_key"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt  time.Time  `json:"-"`
	DeletedAt  *time.Time `sql:"index" json:"-"`
	SystemName string     `gorm:"not null" json:"system_name"`
	OwnerID    string     `gorm:"not null" json:"owner_id"`

	Admin       []Admin                    `gorm:"foreignkey:SystemID" json:"-"`
	News        []modelsNews.News          `gorm:"foreignkey:SystemID" json:"-"`
	NewsType    []modelsNews.NewsType      `gorm:"foreignkey:SystemID" json:"-"`
	Member      []modelsMember.Member      `gorm:"foreignkey:SystemID" json:"-"`
	TargetGroup []modelsMember.TargetGroup `gorm:"foreignkey:SystemID" json:"-"`
	LineOA      LineOA                     `gorm:"foreignkey:SystemID" json:"-"`
	Role        []Role                     `gorm:"foreignkey:SystemID" json:"-"`
	Dialogflow  DialogflowProcessor        `gorm:"foreignkey:SystemID" json:"-"`
}

func (s *System) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", GenSystemID())
	return nil
}

func GenSystemID() string {
	SystemID := "AC-"
	for i := 0; i < 10; i++ {
		ranType := rand.Intn(2)
		switch ranType {
		case 0:
			SystemID += string(rand.Intn(57-48) + 48)
		case 1:
			SystemID += string(rand.Intn(90-65) + 65)
		case 2:
			SystemID += string(rand.Intn(122-97) + 97)
		}
	}
	log.Print(SystemID)
	return SystemID
}

func (system *System) AddAdmin(admin Admin) {
	system.Admin = append(system.Admin, admin)
}
func (system *System) AddLineOA(lineoa LineOA) {
	system.LineOA = lineoa
}
func (system *System) AddNewsTypes(newstype modelsNews.NewsType) {
	system.NewsType = append(system.NewsType, newstype)
}
func (system *System) AddRole(role Role) {
	system.Role = append(system.Role, role)
}
func (system *System) AddTargerGroup(targetgroup modelsMember.TargetGroup) {
	system.TargetGroup = append(system.TargetGroup, targetgroup)
}
