package modelsNews

import (
	"time"

	"github.com/jinzhu/gorm"
)

type News struct {
	gorm.Model
	Title      string    `json:"title"`
	Body       string    `sql:"type:text;" json:"body"`
	ExpireDate time.Time `json:"expire_date"`
	SystemID   string    `json:"system_id"`
	AuthorID   uint      `json:"authr_id"`
	Status     string    `json:"status"` //Draft or Publish

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
