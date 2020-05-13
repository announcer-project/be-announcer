package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsNews"
	"encoding/base64"
	"errors"
	"fmt"
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
func CreateNews(c echo.Context) error {
	data := News{}
	if err := c.Bind(&data); err != nil {
		return err
	}
	authorization := c.Request().Header.Get("Authorization")
	jwt := string([]rune(authorization)[7:])
	tokens, _ := DecodeJWT(jwt)
	db := database.Open()
	defer db.Close()
	admin := models.Admin{}
	db.Where("user_id = ? AND system_id = ?", tokens["user_id"], data.SystemID).Find(&admin)
	if admin.ID == 0 {
		return errors.New("You not admin.")
	}
	system := models.System{}
	db.Where("id = ?", data.SystemID).First(&system)
	if system.ID == 0 {
		return errors.New("Have not this system.")
	}
	expiredate, _ := time.Parse("dd-mm-yy", data.Expiredate)
	news := modelsNews.News{Title: data.Title, Body: data.Body, ExpireDate: expiredate, SystemID: system.ID, AuthorID: admin.ID, Status: data.Status}
	db.Create(&news)
	if news.ID == 0 {
		return errors.New("Create fail.")
	}
	for _, newstype := range data.Newstypes {
		newstypedb := modelsNews.NewsType{}
		db.Where("id = ?", newstype.ID).Find(&newstypedb)
		if newstypedb.ID == 0 {
			return errors.New("Create fail.")
		}
		typeofnews := modelsNews.TypeOfNews{NewsID: news.ID, NewsTypeID: newstypedb.ID}
		db.Create(&typeofnews)
		if newstypedb.ID == 0 {
			return errors.New("Create fail.")
		}
	}
	UploadImages(data.Images, news.ID, system)
	return nil
}

func UploadImages(images []string, newsID uint, system models.System) error {
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
		imagename := system.SystemName + "-" + fmt.Sprint(system.ID) + "-" + fmt.Sprint(newsID) + "-" + strconv.Itoa(i) + `.jpg`
		path := getEnv("FE_PATH", "") + `\public\image\News\` + imagename
		f, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if _, err := f.Write(dec); err != nil {
			panic(err)
		}
		if err := f.Sync(); err != nil {
			panic(err)
		}
		img := modelsNews.Image{NewsID: newsID, ImageName: imagename}
		db.Create(&img)
		if img.ID == 0 {
			return errors.New("Upload error.")
		}
	}
	return nil
}

func GetNewsByID(c echo.Context) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	news := modelsNews.News{}
	if c.Param("id") != "" {
		db.First(&news, c.Param("id"))
	} else {
		db.First(&news, c.FormValue("id"))
	}
	return news, nil
}

func GetAllNews(c echo.Context) (interface{}, error) {
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
	news := []modelsNews.News{}
	db.Where("system_id = ?", c.FormValue("systemid")).Find(&news)
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
	authorization := c.Request().Header.Get("Authorization")
	jwt := string([]rune(authorization)[7:])
	tokens, _ := DecodeJWT(jwt)
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ? AND system_name = ?", c.QueryParam("systemid"), c.QueryParam("systemname")).Find(&system)
	if system.ID == 0 {
		return nil, errors.New("Not have this system.")
	}
	admin := models.Admin{}
	db.Where("user_id = ? AND system_id = ?", tokens["user_id"], c.QueryParam("systemid")).Find(&admin)
	if admin.ID == 0 {
		return nil, errors.New("You not admin in this system.")
	}
	newsTypes := []modelsNews.NewsType{}
	db.Where("system_id = ?", c.QueryParam("systemid")).Find(&newsTypes)
	return newsTypes, nil
}
