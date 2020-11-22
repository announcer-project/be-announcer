package models

import (
	"be_nms/models/modelsLineAPI"

	"github.com/jinzhu/gorm"
)

type LineOA struct {
	gorm.Model
	ChannelID     string `gorm:"unique;not null"`
	ChannelSecret string `gorm:"unique;not null"`
	LiffID        string `gorm:"unique;not null"`
	SystemID      string `gorm:"not null"`

	RichMenu []modelsLineAPI.RichMenu `gorm:"foreignkey:LineOAID"`
}

func (LineOA) TableName() string {
	return "lineoas"
}

func (lineoa *LineOA) AddRichMenu(richmenu modelsLineAPI.RichMenu) {
	lineoa.RichMenu = append(lineoa.RichMenu, richmenu)
}
