package models

import (
	"log"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	ID         string `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
	FName      string `json:"fname"`
	LName      string `json:"lname"`
	Email      string `gorm:"unique" json:"email"`
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
	UserID := "ANNION"
	for i := 0; i < 16; i++ {
		ranType := rand.Intn(2)
		switch ranType {
		case 0:
			UserID += string(rand.Intn(57-48) + 48)
		case 1:
			UserID += string(rand.Intn(90-65) + 65)
		case 2:
			UserID += string(rand.Intn(122-97) + 97)
		}
	}
	log.Print(UserID)
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
