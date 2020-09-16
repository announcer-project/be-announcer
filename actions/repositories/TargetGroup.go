package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsMember"
	"errors"
	"fmt"
)

func CreateTargetGroup(
	userid,
	systemid,
	groupname string,
	members []struct {
		MemberID uint
	}) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	admin := models.Admin{}
	db.Where("user_id = ? AND system_id = ?", userid, systemid).Find(&admin)
	if admin.ID == 0 {
		return nil, errors.New("you not admin in this system.")
	}
	system := models.System{}
	db.Where("id = ?", systemid).Find(&system)
	if system.ID == "" {
		return nil, errors.New("Not have system.")
	}
	targetGroup := modelsMember.TargetGroup{
		TargetGroupName: groupname,
		NumberOfMembers: len(members),
		SystemID:        system.ID,
	}
	for _, member := range members {
		memberDB := modelsMember.Member{}
		db.Where("id = ? and system_id = ?", member.MemberID, systemid).First(&memberDB)
		if memberDB.ID == 0 {
			return nil, errors.New("not have member id " + fmt.Sprint(member.MemberID))
		}
		memberGroup := modelsMember.MemberGroup{MemberID: member.MemberID}
		targetGroup.AddMemberGroup(memberGroup)
	}
	db.Create(&targetGroup)
	if targetGroup.ID == 0 {
		return nil, errors.New("create fail.")
	}
	return targetGroup, nil
}
func GetAllTargetGroup(userid, systemid string) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	admin := models.Admin{}
	db.Where("user_id = ? AND system_id = ?", userid, systemid).Find(&admin)
	if admin.ID == 0 {
		return nil, errors.New("you not admin in this system.")
	}
	targetGroups := []modelsMember.TargetGroup{}
	db.Where("system_id = ?", systemid).Find(&targetGroups)
	return targetGroups, nil
}
