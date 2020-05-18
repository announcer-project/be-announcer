package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsLineAPI"
	"be_nms/models/modelsNews"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/linebot"
)

type LineBroadcastMessage struct {
	Messages []modelsLineAPI.CardLine `json:"messages"`
}

type RichMenuID struct {
	Richmenuid string `json:"richMenuId"`
}

func CreateRichmenu(channelid, channeltoken, system string, systemid uint) (interface{}, error) {
	richMenu := linebot.RichMenu{
		Size:        linebot.RichMenuSize{Width: 2500, Height: 1686},
		Selected:    true,
		Name:        "Default rich menu register",
		ChatBarText: "Register",
		Areas: []linebot.AreaDetail{
			{
				Bounds: linebot.RichMenuBounds{X: 0, Y: 0, Width: 2500, Height: 1686},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypeURI,
					URI:  getEnv("LINE_LIFF", "") + "/line/" + system + "/" + fmt.Sprint(systemid) + "/register",
					Text: "click me",
				},
			},
		},
	}
	log.Print(getEnv("LINE_LIFF", "") + "/line/" + system + "/" + fmt.Sprint(systemid) + "/register")
	bot, err := linebot.New(channelid, channeltoken)
	if err != nil {
		return nil, err
	}
	log.Print(richMenu)
	res, err := bot.CreateRichMenu(richMenu).Do()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	log.Print(res.RichMenuID)
	return res.RichMenuID, nil
}

func SetImageToRichMenu(richmenu, channelid, channeltoken string) error {
	bot, err := linebot.New(channelid, channeltoken)
	if err != nil {
		return err
	}
	imagePath, err := GetFile("rich-menu.png")
	if err != nil {
		return err
	}
	if _, err := bot.UploadRichMenuImage(richmenu, imagePath).Do(); err != nil {
		return err
	}
	os.Remove(imagePath)
	return nil
}

func SetDefaultRichMenu(richmenuid, channelid, channeltoken string) error {
	bot, err := linebot.New(channelid, channeltoken)
	if err != nil {
		return err
	}
	if _, err = bot.SetDefaultRichMenu(richmenuid).Do(); err != nil {
		return err
	}
	return nil
}

// func BroadMessageLine(c echo.Context, news models.News) bool {
// 	text := models.TextLine{}
// 	text.CreateLineMessageText("Test System")
// 	messages := LineBroadcastMessage{}
// 	messages.Messages = []models.TextLine{text}
// 	messagesJSON, _ := json.Marshal(messages)
// 	jsonStr := []byte(messagesJSON)
// 	client := &http.Client{}
// 	request, _ := http.NewRequest("POST", "https://api.line.me/v2/bot/message/broadcast", bytes.NewBuffer(jsonStr))
// 	request.Header.Set("Content-Type", "application/json")
// 	request.Header.Set("Authorization", "Bearer "+getEnv("CHANNEL_ACCESS_TOKEN", ""))
// 	res, err := client.Do(request)
// 	defer res.Body.Close()
// 	if err != nil {
// 		return false
// 	}
// 	return true
// }

func BroadcastNewsLine(c echo.Context, news modelsNews.News, system models.System) (bool, error) {
	lineoa := models.LineOA{}
	db := database.Open()
	defer db.Close()
	db.Where("system_id = ?", system.ID).Find(&lineoa)
	newsCard := modelsLineAPI.CardLine{}
	link := getEnv("LINE_LIFF", "") + "/line/news/" + fmt.Sprint(news.ID)
	titleNews := string([]rune(news.Title)[0:37]) + "..."
	bodyNews := strip.StripTags(news.Body)
	bodyNews = string([]rune(bodyNews)[0:57]) + "..."
	newsCard.CreateCardLine(link, titleNews, bodyNews)
	cards := []modelsLineAPI.CardLine{newsCard}
	messages := LineBroadcastMessage{cards}
	messagesJSON, _ := json.Marshal(messages)
	log.Print(messagesJSON)
	jsonStr := []byte(messagesJSON)
	log.Print(string(jsonStr))
	client := &http.Client{}
	request, _ := http.NewRequest("POST", "https://api.line.me/v2/bot/message/broadcast", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+lineoa.ChannelSecret)
	res, err := client.Do(request)
	defer res.Body.Close()
	if err != nil {
		return false, err
	}
	body, _ := ioutil.ReadAll(res.Body)
	log.Print(string(body))
	return true, nil
}
