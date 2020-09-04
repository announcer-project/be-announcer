package handlers

import (
	"be_nms/actions/repositories"
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsMember"
	"be_nms/models/modelsNews"
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
	newstypes, err := repositories.GetAllNewsType(c.QueryParam("systemid"), false)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	targetgroups, err := repositories.GetAllTargetGroup(c.QueryParam("systemid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	users := []models.User{}
	members, err := repositories.GetAllMember(c.QueryParam("systemid"))
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
	authorization := c.Request().Header.Get("Authorization")
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	log.Print(tokens["user_id"].(string))
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

type LineMessageBroadcast struct {
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

func BroadcastNewsToLine(c echo.Context) error {
	data := LineMessageBroadcast{}
	if err := c.Bind(&data); err != nil {
		log.Print("error ", err)
		return err
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
				"https://announcer-project.s3-ap-southeast-1.amazonaws.com/news/full-team.png",
				message.News.Title,
				body,
				"https://www.google.com",
			)
			messages = append(messages, &container)
		}
	}
	db := database.Open()
	lineoa := models.LineOA{}
	db.Where("system_id = ?", data.SystemID).First(&lineoa)
	bot, err := linebot.New(lineoa.ChannelID, lineoa.ChannelSecret)
	if err != nil {
		return err
	}
	if data.Everyone {
		go repositories.BroadcastToEveryone(messages, bot)
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
	return c.JSON(http.StatusOK, "Announce Success!")
}
