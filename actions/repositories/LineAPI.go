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
	"math/rand"
	"net/http"
	"os"
	"unicode/utf8"

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

func CreateRichmenu(channelid, channelaccesstoken, richname string, richMenu linebot.RichMenu) (interface{}, error) {
	bot, err := linebot.New(channelid, channelaccesstoken)
	if err != nil {
		return nil, err
	}
	res, err := bot.CreateRichMenu(richMenu).Do()
	if err != nil {
		return nil, err
	}
	log.Print(res.RichMenuID)
	return res.RichMenuID, nil
}

func DeleteRichmenu(channelid, channelaccesstoken, richmenuid string) error {
	bot, err := linebot.New(channelid, channelaccesstoken)
	if err != nil {
		log.Print(err)
		return err
	}
	_, err = bot.DeleteRichMenu(richmenuid).Do()
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func SetImageToRichMenu(richmenu, channelid, channelaccesstoken, image string) error {
	bot, err := linebot.New(channelid, channelaccesstoken)
	if err != nil {
		return err
	}
	imagePath, err := GetFile("/richmenu", image)
	if err != nil {
		return err
	}
	if _, err := bot.UploadRichMenuImage(richmenu, imagePath).Do(); err != nil {
		return err
	}
	os.Remove(imagePath)
	return nil
}

func SetDefaultRichMenu(richmenuid, channelid, channelaccesstoken string) error {
	bot, err := linebot.New(channelid, channelaccesstoken)
	if err != nil {
		return err
	}
	if _, err = bot.SetDefaultRichMenu(richmenuid).Do(); err != nil {
		return err
	}
	return nil
}
func SetLinkRichMenu(richmenuid, channelid, channelaccesstoken, lineuserid string) error {
	bot, err := linebot.New(channelid, channelaccesstoken)
	if err != nil {
		return err
	}
	if _, err = bot.LinkUserRichMenu(lineuserid, richmenuid).Do(); err != nil {
		return err
	}
	return nil
}

func CreateNewsCardLine(AltText, CoverUrl, Title, Body, NewsUrl string) linebot.TemplateMessage {
	if utf8.RuneCountInString(Title) > 37 {
		Title = string([]rune(Title)[0:37]) + "..."
	}
	if utf8.RuneCountInString(Body) > 57 {
		Body = string([]rune(Body)[0:57]) + "..."
	}
	container := linebot.TemplateMessage{
		AltText: AltText,
		Template: linebot.NewButtonsTemplate(
			CoverUrl,
			Title,
			Body,
			linebot.NewURIAction("More Detail", NewsUrl)),
	}
	return container
}

func CreateTextLine(text string) linebot.TextMessage {
	container := linebot.TextMessage{
		Text: text,
	}
	return container
}

func CreateImageLine(base64img string) (linebot.ImageMessage, string) {
	session := ConnectFileStorage()
	imageByte := Base64toByte(base64img, "image")
	Filename := "BC-"
	for i := 0; i < 6; i++ {
		ranType := rand.Intn(2)
		switch ranType {
		case 0:
			Filename += string(rand.Intn(57-48) + 48)
		case 1:
			Filename += string(rand.Intn(90-65) + 65)
		case 2:
			Filename += string(rand.Intn(122-97) + 97)
		}
	}
	CreateFile(session, imageByte, Filename+".png", "/broadcast")
	container := linebot.ImageMessage{
		OriginalContentURL: "https://announcer-project.s3-ap-southeast-1.amazonaws.com/broadcast/" + Filename + ".png",
		PreviewImageURL:    "https://announcer-project.s3-ap-southeast-1.amazonaws.com/broadcast/" + Filename + ".png",
	}
	return container, Filename + ".png"
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
	titleNews := news.Title
	log.Print("t:", utf8.RuneCountInString(titleNews))
	if utf8.RuneCountInString(titleNews) > 40 {
		log.Print("t:", utf8.RuneCountInString(titleNews))
		titleNews = string([]rune(news.Title)[0:37]) + "..."
	}
	bodyNews := strip.StripTags(news.Body)
	log.Print("b:", utf8.RuneCountInString(bodyNews))
	if utf8.RuneCountInString(bodyNews) > 40 {
		log.Print("b:", utf8.RuneCountInString(bodyNews))
		bodyNews = string([]rune(bodyNews)[0:57]) + "..."
	}
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
