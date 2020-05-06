package database

import (
	"be_nms/models"
	"be_nms/models/modelsNews"

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
		&modelsNews.News{},
		&modelsNews.Image{},
		&modelsNews.Announcement{},
	)
}

func SetData(db *gorm.DB) {
	// user := models.User{FName: "Panupong", LName: "Joknoi", Email: "panpong.jkn@gmail.com", LineID: "Ufc12c85816992da6381aa3405b9e8083", FacebookID: "", GoogleID: ""}
	// db.Create(&user)
	// db.Create(&user2)
	// db.First(&user)
	// // log.Print(user)
	// system := models.System{SystemName: "NMS", OwnerID: user.ID}
	// db.Create(&system)
	// db.First(&system)
	// admin := models.Admin{UserID: user.ID, SystemID: system.ID, Position: "admin"}
	// db.Create(&admin)
}
