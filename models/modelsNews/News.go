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

	Announcement []Announcement `gorm:"foreignkey:NewsID`
	Image        []Image        `gorm:"foreignkey:NewsID"`
}