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

func CreateSystem(c echo.Context) (interface{}, error) {
	//check jwt
	db := database.Open()
	defer db.Close()
	authorization := c.Request().Header.Get("Authorization")
	jwt := string([]rune(authorization)[7:])
	tokens, _ := DecodeJWT(jwt)
	user := models.User{}
	db.Where("id = ?", tokens["user_id"]).Find(&user)
	if user.ID == "" {
		return nil, errors.New("You not admin.")
	}
	//checkjsonreq
	systemReq := SystemJson{}
	if err := c.Bind(&systemReq); err != nil {
		log.Print("err", err)
		return nil, err
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
		return nil, errors.New("Create Fail.")
	}
	imageByte := Base64toByte(systemReq.SystemProfile)
	sess := ConnectFileStorage()
	if err := CreateFile(sess, imageByte, system.ID+".jpg", "/systems"); err != nil {
		tx.Rollback()
		return nil, errors.New("Upload profile system fail.")
	}
	log.Print("system Req", systemReq)
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
			log.Print("1. richmenu 1 error (Create rich menu pre register fail.) ", err)
			return nil, err
		}
		richmenuPreRegister := modelsLineAPI.RichMenu{RichID: richmenuidPreRegister.(string), Status: "preregister"}
		if err = SetImageToRichMenu(richmenuPreRegister.RichID, systemReq.LineOA.ChannelID, systemReq.LineOA.ChannelAccessToken, "richmenu-register.png"); err != nil {
			tx.Rollback()
			log.Print("2. set image richmenu 1 error ", err)
			return nil, err
		}
		if err = SetDefaultRichMenu(richmenuPreRegister.RichID, systemReq.LineOA.ChannelID, systemReq.LineOA.ChannelAccessToken); err != nil {
			tx.Rollback()
			log.Print("3. set richmenu 1 error ", err)
			return nil, err
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
			log.Print("4. richmenu 2 error Create rich menu after register fail.", err)
			return nil, err
		}
		richmenuAfterRegister := modelsLineAPI.RichMenu{RichID: richmenuidAfterRegister.(string), Status: "afterregister"}
		if err = SetImageToRichMenu(richmenuAfterRegister.RichID, systemReq.LineOA.ChannelID, systemReq.LineOA.ChannelAccessToken, "richmenu-afterregister.png"); err != nil {
			tx.Rollback()
			log.Print("5. set image richmenu 2 error ", err)
			return nil, err
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
			return nil, err
		}
	}
	tx.Commit()

	return system, nil
}
