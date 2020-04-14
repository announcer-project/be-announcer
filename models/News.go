package models

import (
	"log"
	"time"
)

type News struct {
	NewsID     uint `gorm:"primary_key"`
	Title      string
	Content    string
	PostDate   time.Time
	ExpireDate time.Time
	AuthorID   uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
	TypeOfNews []TypeOfNews `gorm:"foreignkey:NewsID"`
}

func (n *News) CreateNews(Title, Content string, ExpireDate interface{}, AuthorID uint) {
	n.Title = Title
	n.Content = Content
	n.PostDate = time.Now()
	switch t := ExpireDate.(type) {
	case time.Time:
		log.Print(t)
		n.ExpireDate = ExpireDate.(time.Time)
	}
	n.AuthorID = AuthorID
}
