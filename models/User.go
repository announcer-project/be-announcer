package models

import (
	"math/rand"
	"strconv"

	"github.com/jinzhu/gorm"
)

type User struct {
	ID         string `gorm:"primary_key"`
	FName      string
	LName      string
	Email      string `gorm:"unique"`
	LineID     string
	FacebookID string
	GoogleID   string

	System []System `gorm:"foreignkey:OwnerID`
	Admin  []Admin  `gorm:"foreignkey:UserID"`
}

func (u *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", GenUserID())
	return nil
}

func GenUserID() string {
	UserID := "NMS"
	for i := 0; i < 16; i++ {
		switch i {
		case 2, 3, 8, 9, 14, 15:
			UserID += string(rand.Intn(122-97) + 97)
		default:
			UserID += strconv.Itoa(rand.Intn(9))
		}
	}
	return UserID
}

func (u *User) CreateUser(Fname, LName, Email, LineID, FacebookID, GoogleID string) {
	u.FName = Fname
	u.LName = LName
	u.Email = Email
	u.LineID = LineID
	u.FacebookID = FacebookID
	u.GoogleID = GoogleID
}
