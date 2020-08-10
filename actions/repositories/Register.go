package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"errors"
	"io/ioutil"

	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) (interface{}, error) {
	db := database.Open()
	user := models.User{}
	db.Where("email = ?", c.FormValue("email")).First(&user)
	if user.ID != "" {
		return user, errors.New("You have account.")
	}
	tx := db.Begin()
	user.CreateUser(c.FormValue("fname"), c.FormValue("lname"), c.FormValue("email"), c.FormValue("line"), c.FormValue("facebook"), c.FormValue("google"))
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("Register fail.")
	}
	sess := ConnectFileStorage()
	if c.FormValue("imagesocial") == "true" {
		fileName := user.ID + ".jpg"
		URL := c.FormValue("imageUrl")
		err := DownloadFile(URL, fileName)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("Register fail.")
		}
		imageByte, err := ioutil.ReadFile(fileName)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("Register fail.")
		}
		if err := CreateFile(sess, imageByte, user.ID+".jpg", "/profile"); err != nil {
			tx.Rollback()
			return nil, errors.New("Register fail.")
		}
	} else {
		imageByte := Base64toByte(c.FormValue("imageprofile"))
		if err := CreateFile(sess, imageByte, user.ID+".jpg", "/profile"); err != nil {
			tx.Rollback()
			return nil, errors.New("Register fail.")
		}
	}
	tx.Commit()
	return user, nil
}

func CheckUserByEmail(email string) (interface{}, error) {
	db := database.Open()
	user := models.User{}
	db.Where("email = ?", email).First(&user)
	if user.ID != "" {
		return user, errors.New("You have account.")
	}
	return nil, nil
}

func ConnectSocialWithAccount(social, socialid, userid string) (interface{}, error) {
	db := database.Open()
	user := models.User{}
	column := social + "_id"
	if err := db.Where("id = ?", userid).First(&user).Update(column, socialid).Error; err != nil {
		return nil, err
	}
	return user, nil
}
