package models

import (
	"log"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	ID         string     `gorm:"primary_key"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt  time.Time  `json:"-"`
	DeletedAt  *time.Time `sql:"index" json:"-"`
	FName      string     `gorm:"not null" json:"fname"`
	LName      string     `gorm:"not null" json:"lname"`
	Email      string     `gorm:"unique;not null" json:"email"`
	LineID     string     `json:"line_id"`
	FacebookID string     `json:"facebook_id"`

	System []System `gorm:"foreignKey:OwnerID;references:ID" json:"-"`
	Admin  []Admin  `gorm:"foreignKey:UserID;references:ID" json:"-"`
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
}
