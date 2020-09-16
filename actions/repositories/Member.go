package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsLineAPI"
	"be_nms/models/modelsMember"
	"be_nms/models/modelsNews"
	"errors"
	"log"

	"github.com/labstack/echo/v4"
)

type User struct {
	ImageUrl       string
	Fname          string
	Lname          string
	RoleID         uint
	NewsInterested []modelsNews.NewsType
	Email          string
	Line           string
	SystemID       string
}

func RegisterGetNews(c echo.Context) (interface{}, error) {
	data := User{}
	if err := c.Bind(&data); err != nil {
		return nil, err
	}
	log.Print("Data: ", data)
	db := database.Open()
	defer db.Close()
	user := models.User{}
	system := models.System{}
	db.Where("id = ?", data.SystemID).First(&system)
	if system.ID == "" {
		return nil, errors.New("not found system")
	}
	lineoa := models.LineOA{}
	db.Where("system_id = ?", system.ID).First(&lineoa)
	if lineoa.ID == 0 {
		return nil, errors.New("not found line oa")
	}
	richmenu := modelsLineAPI.RichMenu{}
	db.Where("line_oa_id = ? and status = ?", lineoa.ID, "afterregister").First(&richmenu)
	if richmenu.ID == 0 {
		return nil, errors.New("not found richmenu afterregister")
	}
	if data.Fname == "" && data.Lname == "" {
		log.Print("user have account")
		db.Where("line_id = ?", data.Line).First(&user)
		log.Print("user", user)
		tx := db.Begin()
		member := modelsMember.Member{UserID: user.ID, SystemID: data.SystemID, RoleID: data.RoleID}
		for _, newstype := range data.NewsInterested {
			memberInterrested := modelsMember.MemberInterested{NewsTypeID: newstype.ID}
			member.AddNewsTypeInterested(memberInterrested)
		}
		tx.Create(&member)
		if member.ID == 0 {
			return nil, errors.New("Create member fail")
		}
		if err := SetAfterRegisterRichMenu(richmenu.RichID, lineoa.ChannelID, lineoa.ChannelSecret, user.LineID); err != nil {
			return nil, err
		}
		tx.Commit()
		return member, nil
	} else {
		log.Print("user not have account")
		user, err := Register(data.Email, data.Fname, data.Lname, data.Line, "", "", true, data.ImageUrl, "")
		if err != nil {
			return nil, errors.New("Register fail")
		}
		tx := db.Begin()
		member := modelsMember.Member{UserID: user.(models.User).ID, SystemID: data.SystemID, RoleID: data.RoleID}
		for _, newstype := range data.NewsInterested {
			memberInterrested := modelsMember.MemberInterested{NewsTypeID: newstype.ID}
			member.AddNewsTypeInterested(memberInterrested)
		}
		tx.Create(&member)
		if member.ID == 0 {
			return nil, errors.New("Create member fail")
		}
		if err := SetAfterRegisterRichMenu(richmenu.RichID, lineoa.ChannelID, lineoa.ChannelSecret, user.(models.User).LineID); err != nil {
			return nil, err
		}
		tx.Commit()
		return member, nil
	}
}

func GetAllMember(userid, systemid string) (interface{}, error) {
	members := []modelsMember.Member{}
	db := database.Open()
	defer db.Close()
	admin := models.Admin{}
	db.Where("user_id = ? and system_id = ?", userid, systemid).First(&admin)
	if admin.ID == 0 {
		return nil, errors.New("you not admin.")
	}
	db.Where("system_id = ?", systemid).Find(&members)
	return members, nil
}
