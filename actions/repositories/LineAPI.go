package repositories

import (
	"be_nms/models/modelsLineAPI"
	"be_nms/models/modelsNews"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/linebot"
)

type LineBroadcastMessage struct {
	Messages []modelsLineAPI.CardLine `json:"messages"`
}

type RichMenuID struct {
	Richmenuid string `json:"richMenuId"`
}

func CreateRichmenu(c echo.Context) (interface{}, error) {
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
					URI:  "https://www.google.com",
					Text: "click me",
				},
			},
		},
	}
	bot, err := linebot.New("5af2298b7e1947a1c459a57a8b04c24c", "DvrYhtBgw8CnugSi1Y+dMMjhwK1gsSWXGv4xTd/kE32zycLc3CjMdhUq63MyR0Gfk0CM202eaiyLMWAQ7EFIzCH/nXY8wg1Gp14+7jq1afGpvyP9dKXN91bZIYvfDusm54b6bDAFfORaWFwRphbN+AdB04t89/1O/w1cDnyilFU=")
	if err != nil {
		return nil, err
	}
	log.Print(bot)
	res, err := bot.CreateRichMenu(richMenu).Do()
	if err != nil {
		return nil, err
	}
	log.Print(res.RichMenuID)
	return res.RichMenuID, nil
}

func SetImageToRichMenu(c echo.Context, richmenu string) error {
	bot, err := linebot.New("5af2298b7e1947a1c459a57a8b04c24c", "DvrYhtBgw8CnugSi1Y+dMMjhwK1gsSWXGv4xTd/kE32zycLc3CjMdhUq63MyR0Gfk0CM202eaiyLMWAQ7EFIzCH/nXY8wg1Gp14+7jq1afGpvyP9dKXN91bZIYvfDusm54b6bDAFfORaWFwRphbN+AdB04t89/1O/w1cDnyilFU=")
	if err != nil {
		return err
	}
	if _, err := bot.UploadRichMenuImage(richmenu, `D:\Downloads\Browser\rich menu #5.png`).Do(); err != nil {
		return err
	}
	return nil
}

func SetDefaultRichMenu(c echo.Context, richmenuid string) error {
	bot, err := linebot.New("5af2298b7e1947a1c459a57a8b04c24c", "DvrYhtBgw8CnugSi1Y+dMMjhwK1gsSWXGv4xTd/kE32zycLc3CjMdhUq63MyR0Gfk0CM202eaiyLMWAQ7EFIzCH/nXY8wg1Gp14+7jq1afGpvyP9dKXN91bZIYvfDusm54b6bDAFfORaWFwRphbN+AdB04t89/1O/w1cDnyilFU=")
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

func BroadcastNewsLine(c echo.Context, news modelsNews.News) (bool, error) {
	newsCard := modelsLineAPI.CardLine{}
	newsCard.CreateCardLine("https://www.google.com", news.Title, news.Body)
	cards := []modelsLineAPI.CardLine{newsCard}
	messages := LineBroadcastMessage{cards}
	messagesJSON, _ := json.Marshal(messages)
	log.Print(messagesJSON)
	jsonStr := []byte(messagesJSON)
	log.Print(string(jsonStr))
	client := &http.Client{}
	request, _ := http.NewRequest("POST", "https://api.line.me/v2/bot/message/broadcast", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+getEnv("CHANNEL_ACCESS_TOKEN", ""))
	res, err := client.Do(request)
	defer res.Body.Close()
	if err != nil {
		return false, err
	}
	body, _ := ioutil.ReadAll(res.Body)
	log.Print(string(body))
	return true, nil
}
