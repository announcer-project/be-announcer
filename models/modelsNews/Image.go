package modelsNews

import "github.com/jinzhu/gorm"

type Image struct {
	gorm.Model
	ImageName string
	NewsID    uint
}
