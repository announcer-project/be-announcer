package models

import (
	"be_nms/models/modelsLineAPI"

	"github.com/jinzhu/gorm"
)

type LineOA struct {
	gorm.Model
	ChannelName   string
	ChannelID     string
	ChannelSecret string
	SystemID      uint

	RichMenu []modelsLineAPI.RichMenu `gorm:"foreignkey:LineOAID"`
}

func (LineOA) TableName() string {
	return "lineoa"
}
