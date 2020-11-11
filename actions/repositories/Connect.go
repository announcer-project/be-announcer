package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsLineAPI"
	"be_nms/models/modelsMember"
	"errors"
	"fmt"
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

func CheckConnectLineOA(systemid, userid string) (bool, error) {
	db := database.Open()
	system := models.System{}
	db.Where("id = ?", systemid).First(&system)
	if system.ID == "" {
		return false, errors.New("not found system.")
	}
	admin := models.Admin{}
	db.Where("system_id = ? and user_id = ?", system.ID, userid).First(&admin)
	if admin.ID == 0 {
		return false, errors.New("you noy admin.")
	}
	lineoa := models.LineOA{}
	db.Where("system_id = ? and deleted_at is null", system.ID).First(&lineoa)
	log.Print(lineoa)
	if lineoa.ID == 0 {
		return false, nil
	}
	return true, nil
}

func DisconnectLineOA(systemid, userid string) error {
	db := database.Open()
	system := models.System{}
	db.Where("id = ?", systemid).First(&system)
	if system.ID == "" {
		return errors.New("not found system.")
	}
	admin := models.Admin{}
	db.Where("system_id = ? and user_id = ?", system.ID, userid).First(&admin)
	if admin.ID == 0 {
		return errors.New("you not admin.")
	}
	lineoa := models.LineOA{}
	db.Where("system_id = ?", system.ID).First(&lineoa)
	if lineoa.ID == 0 {
		return errors.New("not found line oa.")
	}
	richmenus := []modelsLineAPI.RichMenu{}
	db.Where("line_oa_id = ?", lineoa.ID).Find(&richmenus)
	//delete richmenu
	for _, richmenu := range richmenus {
		err := DeleteRichmenu(lineoa.ChannelID, lineoa.ChannelSecret, richmenu.RichID)
		if err != nil {
			return errors.New("delete richmenu fail.")
		}
		err = DeleteFile("/richmenu", richmenu.RichID+".png")
		if err != nil {
			return errors.New("delete imagerichmenu fail.")
		}
	}
	members := []modelsMember.Member{}
	db.Where("system_id = ? and deleted_at is null", system.ID).Find(&members)
	tx := db.Begin()
	for _, member := range members {
		tx.Where("member_id = ? and deleted_at is null", member.ID).Delete(&modelsMember.MemberInterested{})
		tx.Where("member_id = ? and deleted_at is null", member.ID).Delete(&modelsMember.MemberGroup{})
	}
	tx.Where("system_id = ? and deleted_at is null", system.ID).Delete(&models.Role{})
	tx.Where("system_id = ? and deleted_at is null", system.ID).Delete(&modelsMember.Member{})
	tx.Where("line_oa_id = ? and deleted_at is null", lineoa.ID).Delete(&modelsLineAPI.RichMenu{})
	tx.Where("system_id = ? and deleted_at is null", system.ID).Delete(&models.LineOA{})
	tx.Commit()
	return nil
}

func ConnectLineOA(
	systemid,
	userid,
	channelid,
	liffid,
	channelaccesstoken string,
	role []struct {
		Rolename string
		Require  bool
	}) error {

	db := database.Open()
	defer db.Close()
	user := models.User{}
	db.Where("id = ?", userid).First(&user)
	if user.ID == "" {
		return errors.New("you not user.")
	}
	system := models.System{}
	db.Where("id = ?", systemid).First(&system)
	if system.ID == "" {
		return errors.New("not found system.")
	}
	admin := models.Admin{}
	db.Where("system_id = ? and user_id = ?", system.ID, user.ID).First(&admin)
	if admin.ID == 0 {
		return errors.New("ypu not admin.")
	}
	tx := db.Begin()
	richMenuPreRegister := linebot.RichMenu{
		Size:        linebot.RichMenuSize{Width: 2500, Height: 1686},
		Selected:    true,
		Name:        "Register",
		ChatBarText: "Register",
		Areas: []linebot.AreaDetail{
			{
				Bounds: linebot.RichMenuBounds{X: 0, Y: 0, Width: 2500, Height: 1686},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypeURI,
					URI:  getEnv("LINE_LIFF", "") + "/" + system.SystemName + "/" + fmt.Sprint(system.ID) + "/register",
					Text: "click me",
				},
			},
		},
	}
	richmenuidPreRegister, err := CreateRichmenu(channelid, channelaccesstoken, "Register", richMenuPreRegister)
	if err != nil {
		tx.Rollback()
		return errors.New("Channel ID or Channel Access Token invalid")
	}
	richmenuPreRegister := modelsLineAPI.RichMenu{RichID: richmenuidPreRegister.(string), Status: "preregister"}
	if err = SetImageToRichMenu(richmenuPreRegister.RichID, channelid, channelaccesstoken, "richmenu-register.png"); err != nil {
		tx.Rollback()
		return errors.New("set image richmenu 1 error.")
	}
	if err = SetDefaultRichMenu(richmenuPreRegister.RichID, channelid, channelaccesstoken); err != nil {
		tx.Rollback()
		return errors.New("set richmenu 1 error.")
	}
	richMenuWaitApprove := linebot.RichMenu{
		Size:        linebot.RichMenuSize{Width: 2500, Height: 1686},
		Selected:    true,
		Name:        "Wait",
		ChatBarText: "Wait",
		Areas: []linebot.AreaDetail{
			{
				Bounds: linebot.RichMenuBounds{X: 0, Y: 0, Width: 2500, Height: 1686},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypeURI,
					URI:  getEnv("LINE_LIFF", "") + "/" + system.SystemName + "/" + fmt.Sprint(system.ID) + "/register",
					Text: "click me",
				},
			},
		},
	}
	richmenuidWaitApprove, err := CreateRichmenu(channelid, channelaccesstoken, "WaitApprove", richMenuWaitApprove)
	if err != nil {
		tx.Rollback()
		return errors.New("Channel ID or Channel Access Token invalid")
	}
	richmenuWaitApprove := modelsLineAPI.RichMenu{RichID: richmenuidWaitApprove.(string), Status: "waitapprove"}
	if err = SetImageToRichMenu(richmenuWaitApprove.RichID, channelid, channelaccesstoken, "richmenu-waitapprove.png"); err != nil {
		tx.Rollback()
		return errors.New("set image richmenu wait error.")
	}
	richMenuAfterRegister := linebot.RichMenu{
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
	lineoa := models.LineOA{
		ChannelID:     channelid,
		ChannelSecret: channelaccesstoken,
		LiffID:        liffid,
	}
	lineoa.AddRichMenu(richmenuPreRegister)
	lineoa.AddRichMenu(richmenuWaitApprove)
	for _, r := range role {
		richmenuidAfterRegister, err := CreateRichmenu(channelid, channelaccesstoken, "Menu", richMenuAfterRegister)
		if err != nil {
			tx.Rollback()
			return errors.New("richmenu 2 error (create rich menu after register fail.)")
		}
		richmenuAfterRegister := modelsLineAPI.RichMenu{RichID: richmenuidAfterRegister.(string), Status: "afterregister" + r.Rolename}
		if err = SetImageToRichMenu(richmenuAfterRegister.RichID, channelid, channelaccesstoken, "richmenu-afterregister.png"); err != nil {
			tx.Rollback()
			return errors.New("set image richmenu 2 error.")
		}
		lineoa.AddRichMenu(richmenuAfterRegister)
	}

	for _, role := range role {
		system.AddRole(models.Role{RoleName: role.Rolename, Require: role.Require})
	}
	system.AddLineOA(lineoa)
	if err = tx.Save(&system).Error; err != nil {
		tx.Rollback()
		return errors.New("server error.")
	}
	tx.Commit()
	return nil
}

func GetLiffID(systemid string) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ? and deleted is null", systemid).First(&system)
	if system.ID == "" {
		return nil, errors.New("system not found.")
	}
	lineoa := models.LineOA{}
	db.Where("system_id = ? and deleted is null", system.ID).First(&lineoa)
	if lineoa.ID == 0 {
		return nil, errors.New("system not connect lineoa.")
	}
	return lineoa.LiffID, nil
}
