package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsLineAPI"
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
	if system.ID == 0 {
		return nil, errors.New("Create fail.")
	}
	for _, lineoa := range system.LineOA {
		richMenu := linebot.RichMenu{
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
		richmenuid, err := CreateRichmenu(lineoa.ChannelID, lineoa.ChannelSecret, "Register", richMenu)
		if err != nil {
			return nil, err
		}
		richmenu := modelsLineAPI.RichMenu{RichID: richmenuid.(string), Status: "defalut", LineOAID: lineoa.ID}
		db.Create(&richmenu)
		if richmenu.ID == 0 {
			return nil, errors.New("Create rich menu errorr.")
		}
		SetImageToRichMenu(richmenu.RichID, lineoa.ChannelID, lineoa.ChannelSecret, "rich-menu.png")
		SetDefaultRichMenu(richmenu.RichID, lineoa.ChannelID, lineoa.ChannelSecret)
		richMenu2 := linebot.RichMenu{
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
		richmenuid2, err := CreateRichmenu(lineoa.ChannelID, lineoa.ChannelSecret, "Menu", richMenu2)
		if err != nil {
			return nil, err
		}
		richmenu2 := modelsLineAPI.RichMenu{RichID: richmenuid2.(string), Status: "afterregister", LineOAID: lineoa.ID}
		db.Create(&richmenu2)
		if richmenu2.ID == 0 {
			return nil, errors.New("Create rich menu errorr.")
		}
		SetImageToRichMenu(richmenu2.RichID, lineoa.ChannelID, lineoa.ChannelSecret, "rich-afterregister.png")
		// SetAfterRegisterRichMenu(richmenu2.RichID, lineoa.ChannelID, lineoa.ChannelSecret, user.LineID)
	}
	db.Save(&system)
	return "systemReq", nil
}
