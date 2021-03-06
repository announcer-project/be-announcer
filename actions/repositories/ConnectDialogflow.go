package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsNews"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"unicode/utf8"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/line/line-bot-sdk-go/linebot"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

func CheckConnectDialogflow(userID, systemID string) (bool, error) {
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ? and deleted_at is null", systemID).First(&system)
	if system.ID == "" {
		return false, errors.New("system not found.")
	}
	admin := models.Admin{}
	db.Where("user_id = ? and system_id = ?", userID, system.ID).First(&admin)
	if admin.ID == 0 {
		return false, errors.New("you not admin.")
	}
	df := models.DialogflowProcessor{}
	db.Where("system_id = ? and deleted_at is null", system.ID).First(&df)
	if df.ID == 0 {
		return false, nil
	}
	return true, nil
}

func ConnectDialogflow(
	userID,
	systemID,
	projectID,
	authJSONFileBase64,
	lang,
	timeZone string) error {
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ? and deleted_at is null", systemID).First(&system)
	if system.ID == "" {
		return errors.New("system not found.")
	}
	admin := models.Admin{}
	db.Where("user_id = ? and system_id = ?", userID, system.ID).First(&admin)
	if admin.ID == 0 {
		return errors.New("you not admin.")
	}
	fileByte := Base64toByte(authJSONFileBase64, "json")
	err := ioutil.WriteFile("dialogflow/"+projectID+".json", fileByte, 0700)
	if err != nil {
		return err
	}
	defer os.Remove("dialogflow/" + projectID + ".json")

	sess := ConnectFileStorage()
	if err := CreateFile(sess, fileByte, projectID+".json", "/dialogflow"); err != nil {
		return errors.New("upload auth json file fail.")
	}
	ctx := context.Background()
	_, clientErr := dialogflow.NewIntentsClient(ctx, option.WithCredentialsFile("dialogflow/"+projectID+".json"))
	if clientErr != nil {
		return errors.New(clientErr.Error())
	}

	newstypes := []modelsNews.NewsType{}
	db.Where("system_id = ? and deleted_at is null", system.ID).Find(&newstypes)
	messages := []models.Message{}
	for _, newstype := range newstypes {
		msg := models.Message{IntentName: newstype.NewsTypeName, TypeMessage: "news"}
		messages = append(messages, msg)
	}
	df := models.DialogflowProcessor{
		ProjectID:        projectID,
		AuthJSONFilePath: getEnv("STORAGE_PATH", "") + "/dialogflow/" + projectID + ".json",
		Lang:             lang,
		TimeZone:         timeZone,
		Message:          messages,
	}
	system.Dialogflow = df
	tx := db.Begin()
	tx.Save(&system)
	err = DowloadFileJSON(df.AuthJSONFilePath, df.ProjectID+".json")
	if err != nil {
		return err
	}
	defer os.Remove("dialogflow/" + df.ProjectID + ".json")
	df.AuthJSONFilePath = "dialogflow/" + df.ProjectID + ".json"
	err = df.Init()
	for _, newstype := range newstypes {
		err = CreateIntentNewstype(newstype.NewsTypeName, []string{newstype.NewsTypeName}, []string{}, df)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func DisconnectDialogflow(systemid string, userid string) error {
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ? and deleted_at is null", systemid).First(&system)
	if system.ID == "" {
		return errors.New("system not found.")
	}
	admin := models.Admin{}
	db.Where("user_id = ? and system_id = ?", userid, system.ID).First(&admin)
	if admin.ID == 0 {
		return errors.New("you not admin.")
	}
	tx := db.Begin()
	df := models.DialogflowProcessor{}
	db.Where("system_id = ? and deleted_at is null", system.ID).First(&df)
	if df.ID != 0 {
		log.Print("delete 1")
		err := DowloadFileJSON(df.AuthJSONFilePath, df.ProjectID+".json")
		if err != nil {
			return err
		}
		defer os.Remove("dialogflow/" + df.ProjectID + ".json")
		df.AuthJSONFilePath = "dialogflow/" + df.ProjectID + ".json"
		err = df.Init()

		ctx := context.Background()

		intentsClient, clientErr := dialogflow.NewIntentsClient(ctx, option.WithCredentialsFile(df.AuthJSONFilePath))
		if clientErr != nil {
			log.Print(clientErr)
			return clientErr
		}
		defer intentsClient.Close()
		newstypes := []modelsNews.NewsType{}
		db.Where("system_id = ? and deleted_at is null", system.ID).Find(&newstypes)
		intents, err := ListIntents(admin.UserID, system.ID)
		if err != nil {
			log.Print(err)
			return err
		}
		for _, newstype := range newstypes {
			for _, intent := range intents {
				if intent.DisplayName == newstype.NewsTypeName {
					request := dialogflowpb.DeleteIntentRequest{Name: intent.Name}

					requestErr := intentsClient.DeleteIntent(ctx, &request)
					if requestErr != nil {
						log.Print(requestErr)
						tx.Rollback()
						return requestErr
					}
					break
				}
			}

		}
		log.Print("delete 2")
		tx.Where("dialogflow_id = ? and deleted_at is null", df.ID).Delete(&models.Message{})
		tx.Where("system_id = ? and deleted_at is null", system.ID).Delete(&models.DialogflowProcessor{})
	}
	log.Print("delete 3")
	tx.Commit()
	return nil
}

// func TestGet(systemid string) (interface{}, error) {
// 	db := database.Open()
// 	defer db.Close()
// 	multiplenews := []modelsNews.News{}
// 		db.Raw("SELECT * FROM news WHERE system_id = ?", 3).Scan(&multiplenews)
// }

func Webhook(systemid, message, replytoken string) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ? and deleted_at is null", systemid).First(&system)
	if system.ID == "" {
		return nil, errors.New("system not found.")
	}
	df := models.DialogflowProcessor{}
	db.Where("system_id = ? and deleted_at is null", system.ID).First(&df)
	if df.ID == 0 {
		return nil, errors.New("not connect dialogflow.")
	}
	err := DowloadFileJSON(df.AuthJSONFilePath, df.ProjectID+".json")
	if err != nil {
		return nil, err
	}
	defer os.Remove("dialogflow/" + df.ProjectID + ".json")
	df.AuthJSONFilePath = "dialogflow/" + df.ProjectID + ".json"
	err = df.Init()
	if err != nil {
		return nil, err
	}
	response := df.ProcessNLP(message, "testUser")
	msg := models.Message{}
	db.Where("intent_name = ? and dialogflow_id = ? and deleted_at is null", response.Intent, df.ID).First(&msg)
	lineoa := models.LineOA{}
	db.Where("system_id = ? and deleted_at is null", system.ID).First(&lineoa)
	bot, _ := linebot.New(lineoa.ChannelID, lineoa.ChannelSecret)
	if msg.ID == 0 {
		log.Print("from dialogflow")
		textmessage := linebot.NewTextMessage(response.Response)
		bot.ReplyMessage(replytoken, textmessage).Do()
	} else {
		newstype := modelsNews.NewsType{}
		db.Where("news_type_name = ? and system_id = ? and deleted_at is null", msg.IntentName, system.ID).First(&newstype)
		var columns []*linebot.CarouselColumn
		typeofnews := []modelsNews.TypeOfNews{}
		db.Raw("SELECT * FROM typeofnews WHERE news_type_id = ? and deleted_at is null ORDER BY id DESC LIMIT 5", newstype.ID).Scan(&typeofnews)
		if len(typeofnews) != 0 {
			for _, typenews := range typeofnews {
				news := modelsNews.News{}
				db.Where("id = ? and deleted_at is null", typenews.NewsID).First(&news)
				title := news.Title
				body := strip.StripTags(news.Body)
				if utf8.RuneCountInString(title) > 37 {
					title = string([]rune(title)[0:37]) + "..."
				}
				if utf8.RuneCountInString(body) > 57 {
					body = string([]rune(body)[0:57]) + "..."
				}
				cardline := linebot.NewCarouselColumn(
					getEnv("STORAGE_PATH", "")+"/news/"+system.ID+"-"+fmt.Sprint(news.ID)+"-cover.png",
					title,
					body,
					linebot.NewURIAction("More Detail", "https://announcer-system.com/news/"+system.SystemName+"/"+system.ID+"/"+fmt.Sprint(news.ID)))
				columns = append(columns, cardline)
			}
			carousel := linebot.NewCarouselTemplate(columns...)
			template := linebot.NewTemplateMessage("ข่าว", carousel)
			bot.ReplyMessage(replytoken, template).Do()
		} else {
			textmessage := linebot.NewTextMessage("ไม่มีข่าวเกี่ยวกับ " + response.Intent)
			if df.Lang == "en" {
				textmessage = linebot.NewTextMessage("not have news about " + response.Intent)
			}
			bot.ReplyMessage(replytoken, textmessage).Do()
		}
		log.Print("from db")
	}
	return nil, nil
}

func DowloadFileJSON(URL, filename string) error {
	response, err := http.Get(URL)
	if err != nil {
	}
	defer response.Body.Close()

	//Create a empty file
	file, err := os.Create("dialogflow/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}

func ListIntents(userID, systemID string) ([]*dialogflowpb.Intent, error) {
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ? and deleted_at is null", systemID).First(&system)
	if system.ID == "" {
		return nil, errors.New("system not found.")
	}
	admin := models.Admin{}
	db.Where("user_id = ? and system_id = ?", userID, system.ID).First(&admin)
	if admin.ID == 0 {
		return nil, errors.New("you not admin.")
	}
	df := models.DialogflowProcessor{}
	db.Where("system_id = ? and deleted_at is null", system.ID).First(&df)
	if df.ID == 0 {
		return nil, errors.New("system not connect dialogflow.")
	}
	err := DowloadFileJSON(df.AuthJSONFilePath, df.ProjectID+".json")
	if err != nil {
		return nil, err
	}
	defer os.Remove("dialogflow/" + df.ProjectID + ".json")
	df.AuthJSONFilePath = "dialogflow/" + df.ProjectID + ".json"
	err = df.Init()
	ctx := context.Background()

	intentsClient, clientErr := dialogflow.NewIntentsClient(ctx, option.WithCredentialsFile(df.AuthJSONFilePath))
	if clientErr != nil {
		return nil, clientErr
	}
	defer intentsClient.Close()

	parent := fmt.Sprintf("projects/%s/agent", df.ProjectID)

	request := dialogflowpb.ListIntentsRequest{Parent: parent}

	intentIterator := intentsClient.ListIntents(ctx, &request)
	var intents []*dialogflowpb.Intent

	for intent, status := intentIterator.Next(); status != iterator.Done; {
		intents = append(intents, intent)
		intent, status = intentIterator.Next()
	}

	return intents, nil
}

func GetIntent(userID, systemID, IntentName string) (*dialogflowpb.Intent, error) {
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ? and deleted_at is null", systemID).First(&system)
	if system.ID == "" {
		return nil, errors.New("system not found.")
	}
	admin := models.Admin{}
	db.Where("user_id = ? and system_id = ?", userID, system.ID).First(&admin)
	if admin.ID == 0 {
		return nil, errors.New("you not admin.")
	}
	df := models.DialogflowProcessor{}
	db.Where("system_id = ? and deleted_at is null", system.ID).First(&df)
	if df.ID == 0 {
		return nil, errors.New("system not connect dialogflow.")
	}
	err := DowloadFileJSON(df.AuthJSONFilePath, df.ProjectID+".json")
	if err != nil {
		return nil, err
	}
	defer os.Remove("dialogflow/" + df.ProjectID + ".json")
	df.AuthJSONFilePath = "dialogflow/" + df.ProjectID + ".json"
	err = df.Init()

	ctx := context.Background()
	c, err := dialogflow.NewIntentsClient(ctx, option.WithCredentialsFile(df.AuthJSONFilePath))
	if err != nil {
		return nil, err
	}
	defer c.Close()
	// name:="projects/<Project ID>/agent/intents/<Intent ID>"
	req := &dialogflowpb.GetIntentRequest{Name: IntentName, IntentView: dialogflowpb.IntentView_INTENT_VIEW_FULL}
	response, err := c.GetIntent(ctx, req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func CreateIntent(userID, systemID, displayName string, trainingPhraseParts, messageTexts []string) error {
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ? and deleted_at is null", systemID).First(&system)
	if system.ID == "" {
		return errors.New("system not found.")
	}
	admin := models.Admin{}
	db.Where("user_id = ? and system_id = ?", userID, system.ID).First(&admin)
	if admin.ID == 0 {
		return errors.New("you not admin.")
	}
	df := models.DialogflowProcessor{}
	db.Where("system_id = ? and deleted_at is null", system.ID).First(&df)
	if df.ID == 0 {
		return errors.New("system not connect dialogflow.")
	}
	err := DowloadFileJSON(df.AuthJSONFilePath, df.ProjectID+".json")
	if err != nil {
		return err
	}
	defer os.Remove("dialogflow/" + df.ProjectID + ".json")
	df.AuthJSONFilePath = "dialogflow/" + df.ProjectID + ".json"
	err = df.Init()

	ctx := context.Background()

	intentsClient, clientErr := dialogflow.NewIntentsClient(ctx, option.WithCredentialsFile(df.AuthJSONFilePath))
	if clientErr != nil {
		log.Print(clientErr)
		return clientErr
	}
	defer intentsClient.Close()

	parent := fmt.Sprintf("projects/%s/agent", df.ProjectID)

	var targetTrainingPhrases []*dialogflowpb.Intent_TrainingPhrase
	var targetTrainingPhraseParts []*dialogflowpb.Intent_TrainingPhrase_Part
	for _, partString := range trainingPhraseParts {
		part := dialogflowpb.Intent_TrainingPhrase_Part{Text: partString}
		targetTrainingPhraseParts = []*dialogflowpb.Intent_TrainingPhrase_Part{&part}
		targetTrainingPhrase := dialogflowpb.Intent_TrainingPhrase{Type: dialogflowpb.Intent_TrainingPhrase_EXAMPLE, Parts: targetTrainingPhraseParts}
		targetTrainingPhrases = append(targetTrainingPhrases, &targetTrainingPhrase)
	}

	intentMessageTexts := dialogflowpb.Intent_Message_Text{Text: messageTexts}
	wrappedIntentMessageTexts := dialogflowpb.Intent_Message_Text_{Text: &intentMessageTexts}
	intentMessage := dialogflowpb.Intent_Message{Message: &wrappedIntentMessageTexts}

	target := dialogflowpb.Intent{DisplayName: displayName, WebhookState: dialogflowpb.Intent_WEBHOOK_STATE_UNSPECIFIED, TrainingPhrases: targetTrainingPhrases, Messages: []*dialogflowpb.Intent_Message{&intentMessage}}

	request := dialogflowpb.CreateIntentRequest{Parent: parent, Intent: &target}

	_, requestErr := intentsClient.CreateIntent(ctx, &request)
	if requestErr != nil {
		return requestErr
	}

	return nil
}

func CreateIntentNewstype(displayName string, trainingPhraseParts, messageTexts []string, df models.DialogflowProcessor) error {
	ctx := context.Background()

	intentsClient, clientErr := dialogflow.NewIntentsClient(ctx, option.WithCredentialsFile(df.AuthJSONFilePath))
	if clientErr != nil {
		log.Print(clientErr)
		return clientErr
	}
	defer intentsClient.Close()

	parent := fmt.Sprintf("projects/%s/agent", df.ProjectID)

	var targetTrainingPhrases []*dialogflowpb.Intent_TrainingPhrase
	var targetTrainingPhraseParts []*dialogflowpb.Intent_TrainingPhrase_Part
	for _, partString := range trainingPhraseParts {
		part := dialogflowpb.Intent_TrainingPhrase_Part{Text: partString}
		targetTrainingPhraseParts = []*dialogflowpb.Intent_TrainingPhrase_Part{&part}
		targetTrainingPhrase := dialogflowpb.Intent_TrainingPhrase{Type: dialogflowpb.Intent_TrainingPhrase_EXAMPLE, Parts: targetTrainingPhraseParts}
		targetTrainingPhrases = append(targetTrainingPhrases, &targetTrainingPhrase)
	}

	intentMessageTexts := dialogflowpb.Intent_Message_Text{Text: messageTexts}
	wrappedIntentMessageTexts := dialogflowpb.Intent_Message_Text_{Text: &intentMessageTexts}
	intentMessage := dialogflowpb.Intent_Message{Message: &wrappedIntentMessageTexts}

	target := dialogflowpb.Intent{DisplayName: displayName, WebhookState: dialogflowpb.Intent_WEBHOOK_STATE_UNSPECIFIED, TrainingPhrases: targetTrainingPhrases, Messages: []*dialogflowpb.Intent_Message{&intentMessage}}

	request := dialogflowpb.CreateIntentRequest{Parent: parent, Intent: &target}

	_, requestErr := intentsClient.CreateIntent(ctx, &request)
	if requestErr != nil {
		return requestErr
	}

	return nil
}

func DeleteIntent(userID, systemID, IntentName, DisplayName string) error {
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ? and deleted_at is null", systemID).First(&system)
	if system.ID == "" {
		return errors.New("system not found.")
	}
	admin := models.Admin{}
	db.Where("user_id = ? and system_id = ?", userID, system.ID).First(&admin)
	if admin.ID == 0 {
		return errors.New("you not admin.")
	}

	df := models.DialogflowProcessor{}
	db.Where("system_id = ? and deleted_at is null", system.ID).First(&df)
	if df.ID == 0 {
		return errors.New("system not connect dialogflow.")
	}

	message := models.Message{}
	db.Where("intent_name = ? and dialogflow_id = ? and deleted_at is null", DisplayName, df.ID).First(&message)
	if message.ID != 0 {
		return errors.New("You can't delete this intent")
	}

	err := DowloadFileJSON(df.AuthJSONFilePath, df.ProjectID+".json")
	if err != nil {
		return err
	}
	defer os.Remove("dialogflow/" + df.ProjectID + ".json")
	df.AuthJSONFilePath = "dialogflow/" + df.ProjectID + ".json"
	err = df.Init()

	ctx := context.Background()

	intentsClient, clientErr := dialogflow.NewIntentsClient(ctx, option.WithCredentialsFile(df.AuthJSONFilePath))
	if clientErr != nil {
		log.Print(clientErr)
		return clientErr
	}
	defer intentsClient.Close()

	request := dialogflowpb.DeleteIntentRequest{Name: IntentName}

	requestErr := intentsClient.DeleteIntent(ctx, &request)
	if requestErr != nil {
		return requestErr
	}

	return nil
}
