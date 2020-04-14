package database

import (
	"github.com/jinzhu/gorm"
)

func Open() *gorm.DB {
	db, err := gorm.Open("mysql", "root@tcp(localhost:3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
