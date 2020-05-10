package modelsNews

import (
	"time"

	"github.com/jinzhu/gorm"
)

type News struct {
	gorm.Model
	Title      string
	Body       string
	ExpireDate time.Time
	SystemID   uint
	AuthorID   uint
	Status     string //Draft or Publish

	Announcement []Announcement `gorm:"foreignkey:NewsID"`
	Image        []Image        `gorm:"foreignkey:NewsID"`
	TypeOfNews   []TypeOfNews   `gorm:"foreignkey:NewsID"`
}

func (n *News) CreateNews(Title, Body string, ExpireDate time.Time, SystemID, AuthorID uint) {

}
