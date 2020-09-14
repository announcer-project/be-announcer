package modelsNews

import "time"

type TypeOfNews struct {
	ID         uint       `gorm:"primary_key"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt  time.Time  `json:"-"`
	DeletedAt  *time.Time `sql:"index" json:"-"`
	NewsID     uint
	NewsTypeID uint

	NewsTypeName string `gorm:"-"`
}

func (TypeOfNews) TableName() string {
	return "typeofnews"
}
