package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"errors"
	"io/ioutil"
)

func Register(email, fname, lname, line, facebook, google string, imagesocial bool, imageUrl, imageProfile string) (interface{}, error) {
	db := database.Open()
	user := models.User{}
	db.Where("email = ?", email).First(&user)
	if user.ID != "" {
		return user, errors.New("You have account.")
	}
	tx := db.Begin()
	user.CreateUser(fname, lname, email, line, facebook, google)
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("Register fail.")
	}
	sess := ConnectFileStorage()
	if imagesocial == true {
		fileName := user.ID + ".jpg"
		URL := imageUrl
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
		imageByte := Base64toByte(imageProfile)
		if err := CreateFile(sess, imageByte, user.ID+".jpg", "/profile"); err != nil {
			tx.Rollback()
			return nil, errors.New("Register fail.")
		}
	}
	tx.Commit()
	return user, nil
}

func CheckUserByLineID(id string) bool {
	db := database.Open()
	defer db.Close()
	user := models.User{}
	db.Where("line_id = ?", id).First(&user)
	if user.ID != "" {
		return true
	}
	return false
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
