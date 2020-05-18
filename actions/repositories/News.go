package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsNews"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type NewsType struct {
	ID       int
	Newstype string
	Selected bool
}
type News struct {
	Title           string
	Body            string
	Checkexpiredate bool
	Expiredate      string
	Images          []string
	Newstypes       []NewsType
	System          string
	SystemID        string
	Status          string
}

//News
func CreateNews(c echo.Context) (interface{}, error) {
	data := News{}
	if err := c.Bind(&data); err != nil {
		return nil, err
	}
	authorization := c.Request().Header.Get("Authorization")
	jwt := string([]rune(authorization)[7:])
	tokens, _ := DecodeJWT(jwt)
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ?", data.SystemID).First(&system)
	if system.ID == 0 {
		return nil, errors.New("Have not this system.")
	}
	admin := models.Admin{}
	db.Where("user_id = ? AND system_id = ?", tokens["user_id"], system.ID).Find(&admin)
	if admin.ID == 0 {
		return nil, errors.New("You not admin.")
	}
	input := ""
	layout := "02-01-2006"
	if data.Expiredate == "Invalid date" {
		input = "01-01-2000"
	} else {
		input = data.Expiredate
	}
	expiredate, _ := time.Parse(layout, input)
	news := modelsNews.News{Title: data.Title, Body: data.Body, ExpireDate: expiredate, SystemID: system.ID, AuthorID: admin.ID, Status: data.Status}
	for _, newstype := range data.Newstypes {
		newstypedb := modelsNews.NewsType{}
		db.Where("id = ?", newstype.ID).Find(&newstypedb)
		if newstypedb.ID == 0 {
			return nil, errors.New("Create fail.")
		}
		typeofnews := modelsNews.TypeOfNews{NewsID: news.ID, NewsTypeID: newstypedb.ID}
		news.AddTypeOfNews(typeofnews)
	}
	lastNews := modelsNews.News{}
	db.Last(&lastNews)
	if lastNews.ID == 0 {
		UploadImages(data.Images, "1", system, &news)
	} else {
		id := fmt.Sprint(lastNews.ID + 1)
		UploadImages(data.Images, id, system, &news)
	}
	db.Create(&news)
	db.Save(&news)
	if news.ID == 0 {
		return nil, errors.New("Create fail.")
	}
	return news.ID, nil
}

func UploadImages(images []string, newsid string, system models.System, news *modelsNews.News) error {
	db := database.Open()
	defer db.Close()
	for i, image := range images {
		checkbase64 := string([]rune(image)[16:22])
		file := ""
		if checkbase64 == "base64" {
			file = string([]rune(image)[23:])
		} else {
			file = string([]rune(image)[22:])
		}
		dec, err := base64.StdEncoding.DecodeString(file)
		if err != nil {
			panic(err)
		}
		imagename := system.SystemName + "-" + fmt.Sprint(system.ID) + "-" + newsid + "-" + strconv.Itoa(i) + `.jpg`
		err = CreateFile(dec, imagename)
		if err != nil {
			return err
		}
		img := modelsNews.Image{ImageName: imagename}
		news.AddImage(img)
		os.Remove(imagename)
	}
	return nil
}

func GetNewsByID(c echo.Context, id string) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	news := modelsNews.News{}
	db.Where("id = ?", id).Preload("TypeOfNews").Preload("Image").Find(&news)
	if news.ID == 0 {
		return nil, errors.New("Get news error.")
	}
	for i := 0; i < len(news.TypeOfNews); i++ {
		newstype := modelsNews.NewsType{}
		db.Where("id = ?", news.TypeOfNews[i].NewsTypeID).Find(&newstype)
		news.TypeOfNews[i].NewsTypeName = newstype.NewsTypeName
	}
	return news, nil
}

// func GetAllNews(c echo.Context) (interface{}, error) {
// 	authorization := c.Request().Header.Get("Authorization")
// 	jwt := string([]rune(authorization)[7:])
// 	tokens, _ := DecodeJWT(jwt)
// 	db := database.Open()
// 	defer db.Close()
// 	admin := models.Admin{}
// 	db.Where("user_id = ? AND system_id = ?", tokens["user_id"], c.FormValue("systemid")).Find(&admin)
// 	if admin.ID == 0 {
// 		return nil, errors.New("You not admin in this system.")
// 	}
// 	news := []modelsNews.News{}
// 	db.Where("system_id = ?", c.FormValue("systemid")).Find(&news)

// 	return news, nil
// }

func GetAllNews(c echo.Context, status string) (interface{}, error) {
	authorization := c.Request().Header.Get("Authorization")
	jwt := string([]rune(authorization)[7:])
	tokens, _ := DecodeJWT(jwt)
	log.Print("test1")
	db := database.Open()
	defer db.Close()
	log.Print("test2")
	admin := models.Admin{}
	db.Where("user_id = ? AND system_id = ?", tokens["user_id"], c.QueryParam("systemid")).Find(&admin)
	log.Print("test3")
	if admin.ID == 0 {
		return nil, errors.New("You not admin in this system.")
	}
	log.Print("test4")
	log.Print(admin)
	news := []modelsNews.News{}
	db.Where("system_id = ? AND status = ?", c.FormValue("systemid"), status).Preload("TypeOfNews").Find(&news)
	return news, nil
}

//NewsType
func CreateNewsType(c echo.Context) (interface{}, error) {
	authorization := c.Request().Header.Get("Authorization")
	jwt := string([]rune(authorization)[7:])
	tokens, _ := DecodeJWT(jwt)
	db := database.Open()
	defer db.Close()
	admin := models.Admin{}
	db.Where("user_id = ? AND system_id = ?", tokens["user_id"], c.FormValue("systemid")).Find(&admin)
	if admin.ID == 0 {
		return nil, errors.New("You not admin in this system.")
	}
	system := models.System{}
	db.Where("id = ?", c.FormValue("systemid")).Find(&system)
	if system.ID == 0 {
		return nil, errors.New("Not have system.")
	}
	newsType := modelsNews.NewsType{NewsTypeName: c.FormValue("newstypename"), SystemID: system.ID}
	db.Create(&newsType)
	if newsType.ID == 0 {
		return nil, errors.New("Create fail.")
	}
	return newsType, nil
}

func GetAllNewsType(c echo.Context) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ? AND system_name = ?", c.QueryParam("systemid"), c.QueryParam("systemname")).Find(&system)
	if system.ID == 0 {
		return nil, errors.New("Not have this system.")
	}
	newsTypes := []modelsNews.NewsType{}
	db.Where("system_id = ?", c.QueryParam("systemid")).Find(&newsTypes)
	return newsTypes, nil
}
