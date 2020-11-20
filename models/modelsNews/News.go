package modelsNews

import (
	"time"
)

type News struct {
	ID         uint       `gorm:"primary_key"`
	CreatedAt  time.Time  `json:"create_date"`
	UpdatedAt  time.Time  `json:"-"`
	DeletedAt  *time.Time `sql:"index" json:"-"`
	Title      string     `json:"title"`
	Body       string     `sql:"type:text;" json:"body"`
	ExpireDate time.Time  `json:"expire_date"`
	Status     string     `json:"status"` //Draft or Publish
	AuthorID   uint       `json:"author_id"`
	SystemID   string     `json:"system_id"`

	Image      []Image      `gorm:"foreignkey:NewsID"`
	TypeOfNews []TypeOfNews `gorm:"foreignkey:NewsID" json:"type_of_news"`
}

func (n *News) CreateNews(Title, Body string, ExpireDate time.Time, SystemID, AuthorID uint) {

}

func (news *News) AddTypeOfNews(typeofnews TypeOfNews) {
	news.TypeOfNews = append(news.TypeOfNews, typeofnews)
}

func (news *News) AddImage(img Image) {
	news.Image = append(news.Image, img)
}
