package models

import (
	"be_nms/models/modelsLineAPI"

	"github.com/jinzhu/gorm"
)

type LineOA struct {
	gorm.Model
	ChannelID     string
	ChannelSecret string
	SystemID      string

	RichMenu []modelsLineAPI.RichMenu `gorm:"foreignkey:LineOAID"`
}

func (LineOA) TableName() string {
	return "lineoas"
}

func (lineoa *LineOA) AddRichMenu(richmenu modelsLineAPI.RichMenu) {
	lineoa.RichMenu = append(lineoa.RichMenu, richmenu)
}
