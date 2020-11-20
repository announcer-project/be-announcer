package models

import (
	"github.com/jinzhu/gorm"
)

type Message struct {
	gorm.Model
	IntentName   string
	TypeMessage  string
	JSONMessage  string
	DialogflowID uint
}

func (Message) TableName() string {
	return "messages"
}
