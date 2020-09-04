package repositories

import (
	"be_nms/database"
	"be_nms/models"
)

func GetUserByID(userid string) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	user := models.User{}
	db.Where("id = ?", userid).Find(&user)
	return user, nil
}
