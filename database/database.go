package database

import (
	"be_nms/models"
	"log"

	"github.com/jinzhu/gorm"
)

func Open() *gorm.DB {
	db, err := gorm.Open("mysql", "root@tcp(localhost:3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func Migration(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
		&models.System{},
		&models.Admin{},
		&models.News{},
		&models.NewsType{},
		&models.TypeOfNews{})
}

func SetData(db *gorm.DB) {
	user := models.User{}
	user.CreateUser("Panupong", "Joknoi", "panupong.jkn@gmail.com", "Ufc12c85816992da6381aa3405b9e8083", "", "")
	db.Create(&user)
	db.First(&user)
	log.Print(user)
	system := models.System{}
	system.CreateSystem("NMS", user.UserID)
	db.Create(&system)
	db.First(&system)
	admin := models.Admin{}
	admin.CreateAdmin(system.SystemID, user.UserID, "Admin")
	db.Create(&admin)
}
