package database

import (
	"be_nms/models"
	"be_nms/models/modelsLineAPI"
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
		&modelsNews.NewsType{},
		&modelsNews.TypeOfNews{},
		&modelsMember.TargetGroup{},
		&models.LineOA{},
		&models.Role{},
		&modelsMember.Member{},
		&modelsMember.MemberInterested{},
		&modelsMember.MemberGroup{},
		&modelsLineAPI.RichMenu{},
		&models.DialogflowProcessor{},
		&models.Message{},
	)
	db.Model(&models.System{}).AddForeignKey("owner_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Admin{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Admin{}).AddForeignKey("system_id", "systems(id)", "RESTRICT", "RESTRICT")
	db.Model(&modelsNews.News{}).AddForeignKey("system_id", "systems(id)", "RESTRICT", "RESTRICT")
	db.Model(&modelsNews.NewsType{}).AddForeignKey("system_id", "systems(id)", "RESTRICT", "RESTRICT")
	db.Model(&modelsMember.Member{}).AddForeignKey("system_id", "systems(id)", "RESTRICT", "RESTRICT")
	db.Model(&modelsMember.TargetGroup{}).AddForeignKey("system_id", "systems(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.LineOA{}).AddForeignKey("system_id", "systems(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Role{}).AddForeignKey("system_id", "systems(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.DialogflowProcessor{}).AddForeignKey("system_id", "systems(id)", "RESTRICT", "RESTRICT")
	db.Model(&modelsNews.News{}).AddForeignKey("author_id", "admins(id)", "RESTRICT", "RESTRICT")
	db.Model(&modelsNews.News{}).AddForeignKey("author_id", "admins(id)", "RESTRICT", "RESTRICT")
	db.Model(&modelsNews.Image{}).AddForeignKey("news_id", "news(id)", "RESTRICT", "RESTRICT")
	db.Model(&modelsNews.TypeOfNews{}).AddForeignKey("news_id", "news(id)", "RESTRICT", "RESTRICT")
	db.Model(&modelsNews.TypeOfNews{}).AddForeignKey("news_type_id", "newstypes(id)", "RESTRICT", "RESTRICT")
	db.Model(&modelsMember.MemberInterested{}).AddForeignKey("news_type_id", "newstypes(id)", "RESTRICT", "RESTRICT")
	db.Model(&modelsMember.MemberGroup{}).AddForeignKey("target_group_id", "targetgroups(id)", "RESTRICT", "RESTRICT")
	db.Model(&modelsLineAPI.RichMenu{}).AddForeignKey("line_oa_id", "lineoas(id)", "RESTRICT", "RESTRICT")
	db.Model(&modelsMember.Member{}).AddForeignKey("role_id", "roles(id)", "RESTRICT", "RESTRICT")
	db.Model(&modelsMember.MemberGroup{}).AddForeignKey("member_id", "members(id)", "RESTRICT", "RESTRICT")
	db.Model(&modelsMember.MemberInterested{}).AddForeignKey("member_id", "members(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Message{}).AddForeignKey("dialogflow_id", "dialogflows(id)", "RESTRICT", "RESTRICT")
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
