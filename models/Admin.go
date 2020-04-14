package models

import "time"

type Admin struct {
	AdminID   uint `gorm:"primary_key"`
	UserID    string
	SystemID  uint
	Position  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	News      []News `gorm:"foreignkey:AuthorID"`
}

func (a *Admin) CreateAdmin(SystemID uint, UserID, Position string) {
	a.UserID = UserID
	a.SystemID = SystemID
	a.Position = Position
}
