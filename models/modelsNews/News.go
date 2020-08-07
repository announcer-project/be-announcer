package modelsNews

import (
	"time"

	"github.com/jinzhu/gorm"
)

type News struct {
	gorm.Model
	Title      string
	Body       string `sql:"type:text;"`
	ExpireDate time.Time
	SystemID   string
	AuthorID   uint
	Status     string //Draft or Publish

	Announcement []Announcement `gorm:"foreignkey:NewsID"`
	Image        []Image        `gorm:"foreignkey:NewsID"`
	TypeOfNews   []TypeOfNews   `gorm:"foreignkey:NewsID"`
}

func (n *News) CreateNews(Title, Body string, ExpireDate time.Time, SystemID, AuthorID uint) {

}

func (news *News) AddTypeOfNews(typeofnews TypeOfNews) {
	news.TypeOfNews = append(news.TypeOfNews, typeofnews)
}

func (news *News) AddImage(img Image) {
	news.Image = append(news.Image, img)
}
