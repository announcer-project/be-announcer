package models

type User struct {
	UserID     int64 `gorm:"primary_key;"`
	FName      string
	LName      string
	Email      string `gorm:"unique"`
	LineID     string
	FacebookID string
	GoogleID   string
}

func (u *User) CreateUser(Fname, LName, Email, LineID, FacebookID, GoogleID string) bool {
	u.FName = Fname
	u.LName = LName
	u.Email = Email
	u.LineID = LineID
	u.FacebookID = FacebookID
	u.GoogleID = GoogleID
	return true
}
