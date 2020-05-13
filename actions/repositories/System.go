package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"errors"
	"log"

	"github.com/labstack/echo/v4"
)

func GetAllsystems(c echo.Context) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	authorization := c.Request().Header.Get("Authorization")
	log.Print(authorization)
	jwt := string([]rune(authorization)[7:])
	tokens, _ := DecodeJWT(jwt)
	admins := []models.Admin{}
	db.Where("user_id = ?", tokens["user_id"]).Find(&admins)
	systems := []models.System{}
	for _, admin := range admins {
		system := models.System{}
		db.Where("id = ?", admin.SystemID).First(&system)
		systems = append(systems, system)
	}
	return systems, nil
}

type System struct {
	Systemname string
	LineOA     []LineOA
}
type LineOA struct {
	Lineoaname   string
	Channelid    string
	Channeltoken string
}

func CreateSystem(c echo.Context) (interface{}, error) {
	systemReq := System{}
	if err := c.Bind(&systemReq); err != nil {
		return nil, err
	}
	log.Print(systemReq)
	db := database.Open()
	defer db.Close()
	authorization := c.Request().Header.Get("Authorization")
	jwt := string([]rune(authorization)[7:])
	tokens, _ := DecodeJWT(jwt)
	user := models.User{}
	db.Where("id = ?", tokens["user_id"]).Find(&user)
	if user.ID == "" {
		return nil, errors.New("Create fail.")
	}
	system := models.System{SystemName: systemReq.Systemname, OwnerID: user.ID}
	db.Create(&system)
	db.First(&system)
	if system.ID == 0 {
		return nil, errors.New("Create fail.")
	}

	admin := models.Admin{UserID: user.ID, SystemID: system.ID, Position: "admin"}
	db.Create(&admin)
	db.First(&admin)
	if system.ID == 0 {
		return nil, errors.New("Create fail.")
	}
	for _, lineoa := range systemReq.LineOA {
		lineoadb := models.LineOA{ChannelName: lineoa.Lineoaname, ChannelID: lineoa.Channelid, ChannelSecret: lineoa.Channeltoken, SystemID: system.ID}
		db.Create(&lineoadb)
		if lineoadb.ID == 0 {
			return nil, errors.New("Create lineoa fail.")
		}
		richmenuid, err := CreateRichmenu(lineoadb.ChannelID, lineoadb.ChannelSecret)
		if err != nil {
			return nil, err
		}
		SetImageToRichMenu(richmenuid.(string), lineoadb.ChannelID, lineoadb.ChannelSecret)
		SetDefaultRichMenu(richmenuid.(string), lineoadb.ChannelID, lineoadb.ChannelSecret)
	}
	return systemReq, nil
}
