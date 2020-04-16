package repositories

import (
	"be_nms/models"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LineBroadcastMessage struct {
	Messages []models.CardLine `json:"messages"`
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

func BroadcastNewsLine(c echo.Context, news models.News) (bool, error) {
	newsCard := models.CardLine{}
	newsCard.CreateCardLine("https://www.google.com", news.Title, news.Content)
	cards := []models.CardLine{newsCard}
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
