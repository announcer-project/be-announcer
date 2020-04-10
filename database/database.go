package database

import (
	"be_nms/models"

	"github.com/jinzhu/gorm"
)

func Open() *gorm.DB {
	db, err := gorm.Open("mysql", "sql12331330:mynBjQxz6q@tcp(sql12.freemysqlhosting.net:3306)/sql12331330?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&models.User{})
	user := models.User{FName: "Panupong", LName: "Joknoi", Email: "panupong.jkn@gmail.com", LineID: "Ufc12c85816992da6381aa3405b9e8083"}
	db.Create(&user)
	return db
}
