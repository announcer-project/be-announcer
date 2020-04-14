package database

import (
	"be_nms/models"

	"github.com/jinzhu/gorm"
)

func Open() *gorm.DB {
	db, err := gorm.Open("mysql", "root@tcp(localhost:3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(
		&models.User{},
		&models.System{},
		&models.Admin{},
		&models.News{},
		&models.NewsType{},
		&models.TypeOfNews{})
	return db
}

func SetData(db *gorm.DB) {
	user := models.User{}
	user.CreateUser("Panupong", "Joknoi", "panupong.jkn@gmail.com", "Ufc12c85816992da6381aa3405b9e8083", "", "")
	db.Create(&user)
	db.First(&user)
	system := models.System{}
	system.CreateSystem("NMS", user.UserID)
	// system2 := models.System{}
	// system2.CreateSystem("NMS2", user.UserID)
	db.Create(&system)
	// db.Create(&system2)
	db.First(&system)
	admin := models.Admin{}
	admin.CreateAdmin(system.SystemID, user.UserID, "Admin")
	db.Create(&admin)
	// db.First(&admin)
	// news := models.News{}
	// news.CreateNews("Opening NMS", "Today have a new system so god!!!", time.Now(), admin.AdminID)
	// db.Create(&news)
}
