package modelsNews

import "github.com/jinzhu/gorm"

type Announcement struct {
	gorm.Model
	NewsID  uint
	AdminID uint
	Social  string // line or facebook;
}
