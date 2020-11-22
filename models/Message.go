package models

import (
	"github.com/jinzhu/gorm"
)

type Message struct {
	gorm.Model
	IntentName   string `gorm:"not null"`
	TypeMessage  string `gorm:"not null"`
	JSONMessage  string
	DialogflowID uint `gorm:"not null"`
}

func (Message) TableName() string {
	return "messages"
}
