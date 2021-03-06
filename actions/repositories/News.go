package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsMember"
	"be_nms/models/modelsNews"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

//News
func CreateNews(
	cover,
	systemid,
	user_id string,
	checkexpiredate bool,
	expiredate,
	title,
	body,
	status string,
	newstypes []struct {
		ID   int
		Name string
	},
	images []string) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ?", systemid).First(&system)
	if system.ID == "" {
		return nil, errors.New("not have this system.")
	}
	admin := models.Admin{}
	db.Where("user_id = ? AND system_id = ?", user_id, system.ID).Find(&admin)
	if admin.ID == 0 {
		return nil, errors.New("you not admin.")
	}
	input := ""
	layout := "02-01-2006"
	if checkexpiredate {
		input = expiredate
	} else {
		input = "10-02-2000"
	}
	expiredateParse, _ := time.Parse(layout, input)
	news := modelsNews.News{Title: title, Body: body, ExpireDate: expiredateParse, SystemID: system.ID, AuthorID: admin.ID, Status: status}
	for _, newstype := range newstypes {
		newstypedb := modelsNews.NewsType{}
		db.Where("id = ?", newstype.ID).Find(&newstypedb)
		if newstypedb.ID == 0 {
			return nil, errors.New("create fail.")
		}
		typeofnews := modelsNews.TypeOfNews{NewsID: news.ID, NewsTypeID: newstypedb.ID}
		news.AddTypeOfNews(typeofnews)
	}
	tx := db.Begin()
	tx.Create(&news)
	sess := ConnectFileStorage()
	if cover != "" {
		imageByte := Base64toByte(cover, "image")
		imagename := system.ID + "-" + fmt.Sprint(news.ID) + "-cover.png"
		if err := CreateFile(sess, imageByte, imagename, "/news"); err != nil {
			tx.Rollback()
			return nil, errors.New("Register fail.")
		}
	}
	for i, image := range images {
		imageByte := Base64toByte(image, "image")
		imagename := system.ID + "-" + fmt.Sprint(news.ID) + "-" + strconv.Itoa(i) + `.png`
		imagedb := modelsNews.Image{ImageName: imagename, NewsID: news.ID}
		tx.Create(&imagedb)
		if err := CreateFile(sess, imageByte, imagename, "/news"); err != nil {
			tx.Rollback()
			return nil, errors.New("Register fail.")
		}
	}
	tx.Commit()
	return news.ID, nil
}

// func UploadImages(images []string, newsid string, system models.System, news *modelsNews.News) error {
// 	sess := ConnectFileStorage()
// 	for i, image := range images {
// 		imageByte := Base64toByte(image)
// 		imagename := system.SystemName + "-" + system.ID + "-" + newsid + "-" + strconv.Itoa(i) + `.png`
// 		if err := CreateFile(sess, imageByte, imagename, "/news"); err != nil {
// 			return errors.New("Register fail.")
// 		}
// 	}
// 	return nil
// }

func DeleteNewsByID(id string) error {
	db := database.Open()
	defer db.Close()
	tx := db.Begin()
	news := modelsNews.News{}
	db.Where("id = ? and deleted_at is null", id).First(&news)
	if news.ID == 0 {
		return errors.New("news not found")
	}
	tx.Where("news_id = ? and deleted_at is null", news.ID).Delete(&modelsNews.Image{})
	tx.Where("news_id = ? and deleted_at is null", news.ID).Delete(&modelsNews.TypeOfNews{})
	tx.Where("id = ? and deleted_at is null", news.ID).Delete(&modelsNews.News{})
	tx.Commit()
	return nil
}

func GetNewsByID(status, id string) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	news := modelsNews.News{}
	db.Where("id = ? and status = ?", id, status).Preload("TypeOfNews").Preload("Image").Find(&news)
	if news.ID == 0 {
		return nil, errors.New("news not found")
	}
	for i := 0; i < len(news.TypeOfNews); i++ {
		newstype := modelsNews.NewsType{}
		db.Where("id = ?", news.TypeOfNews[i].NewsTypeID).Find(&newstype)
		news.TypeOfNews[i].NewsTypeName = newstype.NewsTypeName
	}
	return news, nil
}

func GetAllNews(userid, systemid string, status string) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	if userid != "publish" {
		admin := models.Admin{}
		db.Where("user_id = ? AND system_id = ?", userid, systemid).Find(&admin)
		if admin.ID == 0 {
			return nil, errors.New("You not admin in this system.")
		}
	}
	news := []modelsNews.News{}
	db.Where("system_id = ? AND status = ?", systemid, status).Preload("TypeOfNews").Find(&news)
	return news, nil
}

func DeleteNews(userid, systemid string, newsid int) error {
	db := database.Open()
	defer db.Close()
	admin := models.Admin{}
	system := models.System{}
	db.Where("id = ?", systemid).Find(&system)
	if system.ID == "" {
		return errors.New("not have system.")
	}
	db.Where("user_id = ? AND system_id = ?", userid, systemid).Find(&admin)
	if admin.ID == 0 {
		return errors.New("You not admin in this system.")
	}
	news := modelsNews.News{}
	db.Where("id = ?", newsid).First(&news)
	db.Where("news_id = ?", news.ID).Delete(&modelsNews.TypeOfNews{})
	db.Delete(&news)
	return nil
}

//NewsType
func CreateNewsType(userid, systemid, newstypename string) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	admin := models.Admin{}
	system := models.System{}
	db.Where("id = ?", systemid).Find(&system)
	if system.ID == "" {
		return nil, errors.New("not have system.")
	}
	db.Where("user_id = ? AND system_id = ?", userid, systemid).Find(&admin)
	if admin.ID == 0 {
		return nil, errors.New("you not admin in this system.")
	}
	newstypeDB := modelsNews.NewsType{}
	db.Where("news_type_name = ? and system_id = ? and deleted_at is null", newstypename, system.ID).First(&newstypeDB)
	if newstypeDB.ID != 0 {
		return nil, errors.New("Duplicate news type name " + newstypename)
	}
	newsType := modelsNews.NewsType{NewsTypeName: newstypename, SystemID: system.ID}
	tx := db.Begin()
	tx.Create(&newsType)
	if newsType.ID == 0 {
		return nil, errors.New("create fail.")
	}
	df := models.DialogflowProcessor{}
	db.Where("system_id = ? and deleted_at is null", system.ID).First(&df)
	if df.ID != 0 {
		msg := models.Message{IntentName: newsType.NewsTypeName, TypeMessage: "news", DialogflowID: df.ID}
		tx.Create(&msg)
		err := DowloadFileJSON(df.AuthJSONFilePath, df.ProjectID+".json")
		if err != nil {
			return nil, err
		}
		defer os.Remove("dialogflow/" + df.ProjectID + ".json")
		df.AuthJSONFilePath = "dialogflow/" + df.ProjectID + ".json"
		err = df.Init()
		err = CreateIntentNewstype(newsType.NewsTypeName, []string{newsType.NewsTypeName}, []string{}, df)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return newsType, nil
}

func DeleteNewsType(userid, systemid string, newstypeid int) error {
	db := database.Open()
	defer db.Close()
	admin := models.Admin{}
	system := models.System{}
	db.Where("id = ?", systemid).Find(&system)
	if system.ID == "" {
		return errors.New("not have system.")
	}
	db.Where("user_id = ? AND system_id = ?", userid, systemid).Find(&admin)
	if admin.ID == 0 {
		return errors.New("You not admin in this system.")
	}
	newstype := modelsNews.NewsType{}
	db.Where("id = ?", newstypeid).First(&newstype)
	typeofnews := []modelsNews.TypeOfNews{}
	db.Where("news_type_id = ? and deleted_at is null", newstype.ID).Find(&typeofnews)
	tx := db.Begin()
	for _, typenews := range typeofnews {
		tx.Where("id = ? and deleted_at is null", typenews.NewsID).Delete(&modelsNews.News{})
	}
	tx.Where("news_type_id = ? and deleted_at is null", newstype.ID).Delete(&modelsMember.MemberInterested{})
	tx.Where("news_type_id = ? and deleted_at is null", newstype.ID).Delete(&modelsNews.TypeOfNews{})
	// DeleteIntent()
	df := models.DialogflowProcessor{}
	db.Where("system_id = ? and deleted_at is null", system.ID).First(&df)
	if df.ID != 0 {
		intent, err := GetIntent(admin.UserID, system.ID, newstype.NewsTypeName)
		if err != nil {
			return err
		}
		err = DeleteIntent(admin.UserID, system.ID, intent.Name, intent.DisplayName)
		if err != nil {
			return err
		}
		tx.Where("intent_name = ? and dialogflow_id = ? and deleted_at is null", newstype.NewsTypeName, df.ID).Delete(&models.Message{})
	}
	tx.Delete(&newstype)
	tx.Commit()
	return nil
}

func GetAllNewsType(systemid string, getnumberofnews bool) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	newsTypes := []modelsNews.NewsType{}
	db.Where("system_id = ? and deleted_at is null", systemid).Find(&newsTypes)
	if getnumberofnews {
		typeofnews := []modelsNews.TypeOfNews{}
		for i, newstype := range newsTypes {
			db.Where("news_type_id = ?", newstype.ID).Find(&typeofnews)
			newsTypes[i].NumberNews = len(typeofnews)
		}
		return newsTypes, nil
	}
	return newsTypes, nil
}

func GetNewsTypeMember(systemid, memberid string) (interface{}, error) {
	var interested []struct {
		NewsType   modelsNews.NewsType `json:'newstype'`
		Interested bool                `json:'interested'`
	}
	db := database.Open()
	defer db.Close()
	newstypes := []modelsNews.NewsType{}
	db.Where("system_id = ? and deleted_at is null", systemid).Find(&newstypes)
	if len(newstypes) == 0 {
		return nil, errors.New("not found news type.")
	}
	for _, newstype := range newstypes {
		newstype_interested := modelsMember.MemberInterested{}
		db.Where("member_id = ? and news_type_id = ? and deleted_at is null", memberid, newstype.ID).First(&newstype_interested)
		var interest struct {
			NewsType   modelsNews.NewsType `json:'newstype'`
			Interested bool                `json:'interested'`
		}
		interest.NewsType = newstype
		interest.Interested = true
		if newstype_interested.ID == 0 {
			interest.Interested = false
		}
		interested = append(interested, interest)
	}
	return interested, nil
}

func SearchNews(systemid, search, newstypeid string) (interface{}, error) {
	search = "%" + search + "%"
	if newstypeid == "0" {
		newstypeid = "%"
	}
	db := database.Open()
	news := []modelsNews.News{}
	db.Joins("JOIN typeofnews ON typeofnews.news_id = news.id AND news.deleted_at is null AND typeofnews.news_type_id LIKE ? AND news.system_id = ?", newstypeid, systemid).Where("news.title LIKE ?", search).Or("news.body LIKE ?", search).Find(&news)
	return news, nil
}
