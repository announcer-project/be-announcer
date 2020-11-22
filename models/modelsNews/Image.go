package modelsNews

import "github.com/jinzhu/gorm"

type Image struct {
	gorm.Model
	ImageName string `gorm:"unique;not null"`
	NewsID    uint   `gorm:"not null"`
}
