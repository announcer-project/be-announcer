package models

import "time"

type TypeOfNews struct {
	TypeOfNewsID uint `gorm:"primary_key"`
	NewsID       uint
	NewsTypeID   uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

func (t *TypeOfNews) CreateTypeOfNews(NewsID, NewsTypeID uint) {
	t.NewsID = NewsID
	t.NewsTypeID = NewsTypeID
}
