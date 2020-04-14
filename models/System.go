package models

import "time"

type System struct {
	SystemID   uint `gorm:"primary_key"`
	SystemName string
	OwnerID    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
	Admin      []Admin    `gorm:"foreignkey:SystemID"`
	NewsType   []NewsType `gorm:"foreignkey:NewsTypeID"`
}

func (s *System) CreateSystem(SystemName, OwnerID string) {
	s.SystemName = SystemName
	s.OwnerID = OwnerID
}
