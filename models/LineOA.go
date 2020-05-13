package models

import (
	"github.com/jinzhu/gorm"
)

type LineOA struct {
	gorm.Model
	ChannelName   string
	ChannelID     string
	ChannelSecret string
	SystemID      uint
}

func (LineOA) TableName() string {
	return "lineoa"
}
