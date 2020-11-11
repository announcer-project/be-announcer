package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsLineAPI"
	"be_nms/models/modelsMember"
	"be_nms/models/modelsNews"
	"errors"
	"log"
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

func RegisterGetNews(data struct {
	IsUser         bool
	FName          string
	LName          string
	Email          string
	ImageUrl       string
	RoleID         int
	NewsInterested []modelsNews.NewsType
	SystemID       string
	Line           string
}) error {
	db := database.Open()
	defer db.Close()
	log.Print("2")
	system := models.System{}
	db.Where("id = ? and deleted_at is null", data.SystemID).First(&system)
	if system.ID == "" {
		return errors.New("not found system.")
	}
	lineoa := models.LineOA{}
	db.Where("system_id = ? and deleted_at is null", system.ID).First(&lineoa)
	if lineoa.ID == 0 {
		return errors.New("not found line oa.")
	}
	role := models.Role{}
	db.Where("system_id = ? and id = ? and deleted_at is null", system.ID, data.RoleID).First(&role)
	if role.ID == 0 {
		return errors.New("not found role.")
	}
	richmenu := modelsLineAPI.RichMenu{}
	if role.Require {
		db.Where("line_oa_id = ? and status = ? and deleted_at is null", lineoa.ID, "waitapprove").First(&richmenu)
		if richmenu.ID == 0 {
			return errors.New("not found richmenu afterregister")
		}
	} else {
		db.Where("line_oa_id = ? and status = ? and deleted_at is null", lineoa.ID, "afterregister"+role.RoleName).First(&richmenu)
		if richmenu.ID == 0 {
			return errors.New("not found richmenu afterregister")
		}
	}
	tx := db.Begin()
	member := modelsMember.Member{}
	member.LineID = data.Line
	member.FName = data.FName
	member.LName = data.LName
	member.SystemID = system.ID
	member.RoleID = role.ID
	if role.Require {
		member.Approve = false
	} else {
		member.Approve = true
	}
	for _, newstype := range data.NewsInterested {
		memberInterrested := modelsMember.MemberInterested{NewsTypeID: newstype.ID}
		member.AddNewsTypeInterested(memberInterrested)
	}
	tx.Create(&member)
	if member.ID == "" {
		return errors.New("Create member fail")
	}
	targetgroup := modelsMember.TargetGroup{}
	db.Where("target_group_name = ? and deleted is null", role.RoleName).First(&targetgroup)
	membergroup := modelsMember.MemberGroup{MemberID: member.ID, TargetGroupID: targetgroup.ID}
	tx.Create(&membergroup)
	targetgroup.NumberOfMembers = targetgroup.NumberOfMembers + 1
	tx.Update(&targetgroup)
	if err := SetLinkRichMenu(richmenu.RichID, lineoa.ChannelID, lineoa.ChannelSecret, member.LineID); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func GetAllMember(userid, systemid string) (interface{}, error) {
	members := []modelsMember.Member{}
	db := database.Open()
	defer db.Close()
	admin := models.Admin{}
	db.Where("user_id = ? and system_id = ? and deleted_at is null", userid, systemid).First(&admin)
	if admin.ID == 0 {
		return nil, errors.New("you not admin.")
	}
	db.Where("system_id = ? and approve = ?", systemid, true).Find(&members)
	return members, nil
}
