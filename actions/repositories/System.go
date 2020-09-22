package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsLineAPI"
	"be_nms/models/modelsNews"
	"errors"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/linebot"
)

func GetAllsystems(user_id string) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	admins := []models.Admin{}
	db.Where("user_id = ?", user_id).Find(&admins)
	for i, admin := range admins {
		system := models.System{}
		db.Where("id = ?", admin.SystemID).First(&system)
		admins[i].System = system
	}
	return admins, nil
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
	if system.ID == "" {
		return nil, errors.New("System not found.")
	}
	return system, nil
}

type SystemJson struct {
	SystemProfile string
	Systemname    string
	NewsTypes     []string
	LineOA        struct {
		ChannelID          string
		ChannelAccessToken string
		RoleUsers          []struct {
			RoleName string
			Require  bool
		}
	}
}

func CreateSystem(user_id string, data interface{}) (interface{}, error) {
	systemReq := data.(struct {
		SystemProfile string
		Systemname    string
		NewsTypes     []string
		LineOA        struct {
			ChannelID          string
			ChannelAccessToken string
			RoleUsers          []struct {
				RoleName string
				Require  bool
			}
		}
	})

	db := database.Open()
	defer db.Close()
	user := models.User{}
	db.Where("id = ?", user_id).First(&user)
	if user.ID == "" {
		return nil, errors.New("you not user.")
	}
	system := models.System{SystemName: systemReq.Systemname, OwnerID: user.ID}
	system.AddAdmin(models.Admin{UserID: user.ID, Position: "admin"})
	for _, newstype := range systemReq.NewsTypes {
		system.AddNewsTypes(modelsNews.NewsType{NewsTypeName: newstype})
	}
	tx := db.Begin()
	tx.Create(&system)
	if system.ID == "" {
		tx.Rollback()
		return nil, errors.New("create fail.")
	}
	imageByte := Base64toByte(systemReq.SystemProfile)
	sess := ConnectFileStorage()
	if err := CreateFile(sess, imageByte, system.ID+".png", "/systems"); err != nil {
		tx.Rollback()
		return nil, errors.New("upload profile system fail.")
	}
	if systemReq.LineOA.ChannelID != "" {
		richMenuPreRegister := linebot.RichMenu{
			Size:        linebot.RichMenuSize{Width: 2500, Height: 1686},
			Selected:    true,
			Name:        "Register",
			ChatBarText: "Register",
			Areas: []linebot.AreaDetail{
				{
					Bounds: linebot.RichMenuBounds{X: 0, Y: 0, Width: 2500, Height: 1686},
					Action: linebot.RichMenuAction{
						Type: linebot.RichMenuActionTypeURI,
						URI:  getEnv("LINE_LIFF", "") + "/line/" + system.SystemName + "/" + fmt.Sprint(system.ID) + "/register",
						Text: "click me",
					},
				},
			},
		}
		richmenuidPreRegister, err := CreateRichmenu(systemReq.LineOA.ChannelID, systemReq.LineOA.ChannelAccessToken, "Register", richMenuPreRegister)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("Channel ID or Channel Access Token invalid")
		}
		richmenuPreRegister := modelsLineAPI.RichMenu{RichID: richmenuidPreRegister.(string), Status: "preregister"}
		if err = SetImageToRichMenu(richmenuPreRegister.RichID, systemReq.LineOA.ChannelID, systemReq.LineOA.ChannelAccessToken, "richmenu-register.png"); err != nil {
			tx.Rollback()
			return nil, errors.New("set image richmenu 1 error.")
		}
		if err = SetDefaultRichMenu(richmenuPreRegister.RichID, systemReq.LineOA.ChannelID, systemReq.LineOA.ChannelAccessToken); err != nil {
			tx.Rollback()
			return nil, errors.New("set richmenu 1 error.")
		}
		richMenuAfterRegister := linebot.RichMenu{
			Size:        linebot.RichMenuSize{Width: 2500, Height: 1686},
			Selected:    true,
			Name:        "Menu",
			ChatBarText: "Menu",
			Areas: []linebot.AreaDetail{
				{
					Bounds: linebot.RichMenuBounds{X: 0, Y: 0, Width: 1683, Height: 839},
					Action: linebot.RichMenuAction{
						Type: linebot.RichMenuActionTypeURI,
						URI:  "https://www.sit.kmutt.ac.th/",
						Text: "click me",
					},
				},
				{
					Bounds: linebot.RichMenuBounds{X: 1683, Y: 0, Width: 817, Height: 839},
					Action: linebot.RichMenuAction{
						Type: linebot.RichMenuActionTypeMessage,
						Text: "โปรไฟล์ของฉัน",
					},
				},
				{
					Bounds: linebot.RichMenuBounds{X: 0, Y: 834, Width: 830, Height: 852},
					Action: linebot.RichMenuAction{
						Type: linebot.RichMenuActionTypeMessage,
						Text: "ทุนการศึกษา",
					},
				},
				{
					Bounds: linebot.RichMenuBounds{X: 830, Y: 839, Width: 853, Height: 847},
					Action: linebot.RichMenuAction{
						Type: linebot.RichMenuActionTypeMessage,
						Text: "ผลงานและกิจกรรม",
					},
				},
				{
					Bounds: linebot.RichMenuBounds{X: 1682, Y: 839, Width: 818, Height: 847},
					Action: linebot.RichMenuAction{
						Type: linebot.RichMenuActionTypeMessage,
						Text: "อยากคุยกับน้องบอท",
					},
				},
			},
		}
		richmenuidAfterRegister, err := CreateRichmenu(systemReq.LineOA.ChannelID, systemReq.LineOA.ChannelAccessToken, "Menu", richMenuAfterRegister)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("richmenu 2 error (create rich menu after register fail.)")
		}
		richmenuAfterRegister := modelsLineAPI.RichMenu{RichID: richmenuidAfterRegister.(string), Status: "afterregister"}
		if err = SetImageToRichMenu(richmenuAfterRegister.RichID, systemReq.LineOA.ChannelID, systemReq.LineOA.ChannelAccessToken, "richmenu-afterregister.png"); err != nil {
			tx.Rollback()
			return nil, errors.New("set image richmenu 2 error.")
		}
		for _, role := range systemReq.LineOA.RoleUsers {
			system.AddRole(models.Role{RoleName: role.RoleName, Require: role.Require})
		}
		lineoa := models.LineOA{
			ChannelID:     systemReq.LineOA.ChannelID,
			ChannelSecret: systemReq.LineOA.ChannelAccessToken,
		}
		lineoa.AddRichMenu(richmenuPreRegister)
		lineoa.AddRichMenu(richmenuAfterRegister)
		system.AddLineOA(lineoa)
		if err = tx.Save(&system).Error; err != nil {
			return nil, errors.New("server error.")
		}
	}
	tx.Commit()

	return system, nil
}
