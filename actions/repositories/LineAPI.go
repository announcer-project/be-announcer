package repositories

import (
	"be_nms/models"
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LineBroadcastMessage struct {
	Messages []models.Text `json:"messages"`
}

func BroadMessageLine(c echo.Context, news models.News) bool {
	text := models.Text{}
	text.CreateLineMessageText("Test System")
	messages := LineBroadcastMessage{}
	messages.Messages = []models.Text{text}
	messagesJSON, _ := json.Marshal(messages)
	jsonStr := []byte(messagesJSON)
	client := &http.Client{}
	request, _ := http.NewRequest("POST", "https://api.line.me/v2/bot/message/broadcast", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+getEnv("CHANNEL_ACCESS_TOKEN", ""))
	res, err := client.Do(request)
	defer res.Body.Close()
	if err != nil {
		return false
	}
	return true
}
