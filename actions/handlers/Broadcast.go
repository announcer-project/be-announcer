package handlers

import (
	"be_nms/actions/repositories"
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsMember"
	"be_nms/models/modelsNews"
	"fmt"
	"log"
	"net/http"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/linebot"
)

type AboutLineBroadcast struct {
	NewsTypes    []modelsNews.NewsType      `json:"newstypes"`
	TargetGroups []modelsMember.TargetGroup `json:"targetgroups"`
	Users        []models.User              `json:"users"`
	News         []modelsNews.News          `json:"news"`
}

func GetAboutLineBroadcast(c echo.Context) error {
	var message struct {
		Message string `json:"message"`
	}
	authorization := c.Request().Header.Get("Authorization")
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	newstypes, err := repositories.GetAllNewsType(c.QueryParam("systemid"), false)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	targetgroups, err := repositories.GetAllTargetGroup(tokens["user_id"].(string), c.QueryParam("systemid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	users := []models.User{}
	members, err := repositories.GetAllMember(tokens["user_id"].(string), c.QueryParam("systemid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	for _, member := range members.([]modelsMember.Member) {
		user, err := repositories.GetUserByID(member.UserID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		users = append(users, user.(models.User))
	}
	news, err := repositories.GetAllNews(tokens["user_id"].(string), c.QueryParam("systemid"), "Publish")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	aboutLineBroadcast := AboutLineBroadcast{
		NewsTypes:    newstypes.([]modelsNews.NewsType),
		TargetGroups: targetgroups.([]modelsMember.TargetGroup),
		Users:        users,
		News:         news.([]modelsNews.News),
	}
	return c.JSON(http.StatusOK, aboutLineBroadcast)
}

func BroadcastNewsToLine(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	admin := models.Admin{}
	system := models.System{}
	var data struct {
		SystemID          string
		Everyone          bool
		CheckNewsTypes    bool
		NewsTypes         []modelsNews.NewsType
		CheckTargetGroups bool
		TargetGroups      []modelsMember.TargetGroup
		CheckUsers        bool
		Users             []models.User
		Messages          []struct {
			Type string
			Data string
			News modelsNews.News
		}
	}
	if err := c.Bind(&data); err != nil {
		message.Message = "server error."
		return c.JSON(500, message)
	}
	log.Print(data)
	db := database.Open()
	defer db.Close()
	db.Where("id = ?", data.SystemID).Find(&system)
	if system.ID == "" {
		message.Message = "not have system."
		return c.JSON(401, message)
	}
	db.Where("user_id = ? AND system_id = ?", tokens["user_id"].(string), system.ID).Find(&admin)
	if admin.ID == 0 {
		message.Message = "you not admin in this system."
		return c.JSON(401, message)
	}
	var messages []linebot.SendingMessage
	var images []string
	for _, message := range data.Messages {
		switch message.Type {
		case "text":
			container := repositories.CreateTextLine(strip.StripTags(message.Data))
			messages = append(messages, &container)
		case "image":
			container, filename := repositories.CreateImageLine(message.Data)
			messages = append(messages, &container)
			images = append(images, filename)
		case "news":
			body := strip.StripTags(message.News.Body)
			container := repositories.CreateNewsCardLine(
				message.News.Title,
				"https://announcer-project.s3-ap-southeast-1.amazonaws.com/news/"+system.ID+"-"+fmt.Sprint(message.News.ID)+"-"+"cover.png",
				message.News.Title,
				body,
				"https://announcer-system.com/news/"+system.SystemName+"/"+system.ID+"/"+fmt.Sprint(message.News.ID),
			)
			messages = append(messages, &container)
		}
	}
	lineoa := models.LineOA{}
	db.Where("system_id = ?", system.ID).First(&lineoa)
	bot, err := linebot.New(lineoa.ChannelID, lineoa.ChannelSecret)
	if err != nil {
		message.Message = "connect line api fail."
		return c.JSON(500, message)
	}
	if data.Everyone {
		go repositories.BroadcastToEveryone(messages, bot, system.ID)
	} else {
		go repositories.BroadcastToSelected(
			messages,
			bot,
			data.CheckNewsTypes,
			data.CheckTargetGroups,
			data.CheckUsers,
			data.NewsTypes,
			data.TargetGroups,
			data.Users,
		)
	}
	message.Message = "announce success."
	return c.JSON(http.StatusOK, message)
}
