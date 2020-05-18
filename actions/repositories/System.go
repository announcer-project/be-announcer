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

func GetSystemByID(c echo.Context, id string) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	authorization := c.Request().Header.Get("Authorization")
	log.Print(authorization)
	jwt := string([]rune(authorization)[7:])
	tokens, _ := DecodeJWT(jwt)
	admin := models.Admin{}
	db.Where("user_id = ? AND system_id = ?", tokens["user_id"], id).Find(&admin)
	if admin.ID == 0 {
		return nil, errors.New("You not admin.")
	}
	system := models.System{}
	db.Where("id = ?", id).Find(&system)
	if system.ID == 0 {
		return nil, errors.New("System not found.")
	}
	return system, nil
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
	system.AddAdmin(models.Admin{UserID: user.ID, Position: "admin"})
	for _, lineoa := range systemReq.LineOA {
		system.AddLineOA(models.LineOA{
			ChannelID:     lineoa.Channelid,
			ChannelName:   lineoa.Lineoaname,
			ChannelSecret: lineoa.Channeltoken,
		})
	}
	db.Create(&system)
	db.Save(&system)
	if system.ID == 0 {
		return nil, errors.New("Create fail.")
	}
	for _, lineoa := range system.LineOA {
		richmenuid, err := CreateRichmenu(lineoa.ChannelID, lineoa.ChannelSecret, system.SystemName, system.ID)
		if err != nil {
			return nil, err
		}
		SetImageToRichMenu(richmenuid.(string), lineoa.ChannelID, lineoa.ChannelSecret)
		SetDefaultRichMenu(richmenuid.(string), lineoa.ChannelID, lineoa.ChannelSecret)
	}
	return "systemReq", nil
}
