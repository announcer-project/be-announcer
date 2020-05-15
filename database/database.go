package database

import (
	"be_nms/models"
	"be_nms/models/modelsMember"
	"be_nms/models/modelsNews"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func Open() *gorm.DB {
	// db, err := gorm.Open("mysql", "root@tcp(localhost:3306)/test?charset=utf8&parseTime=True&loc=Local")
	db, err := gorm.Open(getEnv("DB_TYPE", ""), getEnv("DB_USERNAME", "")+`:`+getEnv("DB_PASSWORD", "")+`@tcp(`+getEnv("DB_HOST", "")+`:`+getEnv("DB_PORT", "")+`)/`+getEnv("DB_NAME", "")+`?charset=utf8&parseTime=True&loc=Local`)
	if err != nil {
		log.Print(err)
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
		&modelsNews.NewsType{},
		&modelsNews.TypeOfNews{},
		&modelsMember.TargetGroup{},
		&models.LineOA{},
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
