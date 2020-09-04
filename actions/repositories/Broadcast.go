package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsMember"
	"be_nms/models/modelsNews"
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

func BroadcastToEveryone(messages []linebot.SendingMessage, bot *linebot.Client) {
	db := database.Open()
	line_id := []string{}
	members := []modelsMember.Member{}
	db.Find(&members)
	for _, member := range members {
		user := models.User{}
		db.Where("id = ?", member.UserID).First(&user)
		line_id = append(line_id, user.LineID)
	}
	bot.Multicast(line_id, messages...).Do()
}

func BroadcastToSelected(
	messages []linebot.SendingMessage,
	bot *linebot.Client,
	CheckNewsTypes,
	CheckTargetGroups,
	CheckUsers bool,
	newstypes []modelsNews.NewsType,
	targetgroups []modelsMember.TargetGroup,
	users []models.User,
) {
	line_id_hash := make(map[string]string)
	if CheckNewsTypes {
		line_id := GetLineIDMemberInterested(newstypes)
		for _, lineid := range line_id {
			line_id_hash[lineid] = lineid
		}
	}
	if CheckTargetGroups {

	}
	if CheckUsers {
		for _, user := range users {
			line_id_hash[user.LineID] = user.LineID
		}
	}
	line_id := []string{}
	for key, _ := range line_id_hash {
		line_id = append(line_id, key)
	}
	log.Print("line_id : ", line_id)
	bot.Multicast(line_id, messages...).Do()
}

func GetLineIDMemberInterested(newstypes []modelsNews.NewsType) []string {
	db := database.Open()
	line_id_hash := make(map[string]string)
	for _, newstype := range newstypes {
		members_interested := []modelsMember.MemberInterested{}
		db.Where("news_type_id = ?", newstype.ID).Find(&members_interested)
		for _, member_interested := range members_interested {
			member := modelsMember.Member{}
			db.Where("id = ?", member_interested.MemberID).First(&member)
			user := models.User{}
			db.Where("id = ?", member.UserID).First(&user)
			line_id_hash[user.LineID] = user.LineID
		}
	}
	line_id := []string{}
	for key, _ := range line_id_hash {
		line_id = append(line_id, key)
	}
	log.Print(line_id)
	return line_id
}
