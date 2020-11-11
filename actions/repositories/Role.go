package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsLineAPI"
	"be_nms/models/modelsMember"
	"errors"
)

func CreateRole(userid, systemid, rolename string, require bool) (interface{}, error) {
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
		return nil, errors.New("not have system.")
	}
	role := models.Role{RoleName: rolename, Require: require, SystemID: system.ID}
	tx := db.Begin()
	tx.Create(&role)
	if role.ID == 0 {
		tx.Rollback()
		return nil, errors.New("create role fail.")
	}
	targetgroup := modelsMember.TargetGroup{TargetGroupName: rolename, NumberOfMembers: 0, SystemID: system.ID}
	tx.Create(&targetgroup)
	if role.ID == 0 {
		tx.Rollback()
		return nil, errors.New("create targetgroup of role fail.")
	}
	tx.Commit()
	return role, nil
}

func GetAllRole(systemid string) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	roleuser := []models.Role{}
	db.Where("system_id = ?", systemid).Find(&roleuser)
	return roleuser, nil
}

func GetRoleRequest(systemid string) (interface{}, error) {
	members := []modelsMember.Member{}
	db := database.Open()
	defer db.Close()
	db.Where("system_id = ? and approve = ?", systemid, false).Find(&members)
	var memberrequests []struct {
		Key  string `json:"key"`
		Name string `json:"name"`
		Role string `json:"role"`
	}
	for _, member := range members {
		role := models.Role{}
		var memberrequest struct {
			Key  string `json:"key"`
			Name string `json:"name"`
			Role string `json:"role"`
		}
		db.Where("id = ?", member.RoleID).First(&role)
		memberrequest.Key = member.ID
		memberrequest.Name = member.FName + " " + member.LName
		memberrequest.Role = role.RoleName
		memberrequests = append(memberrequests, memberrequest)
	}
	return memberrequests, nil
}

func ApproveRoleRequest(memberid uint, userid, systemid string) error {
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ? and deleted_at is null", systemid).First(&system)
	if system.ID == "" {
		return errors.New("not found system.")
	}
	admin := models.Admin{}
	db.Where("user_id = ? and system_id = ? and deleted_at is null", userid, systemid).First(&admin)
	if admin.ID == 0 {
		return errors.New("you not admin.")
	}
	lineoa := models.LineOA{}
	db.Where("system_id = ? and deleted_at is null", system.ID).First(&lineoa)
	if lineoa.ID == 0 {
		return errors.New("not found line oa.")
	}
	member := modelsMember.Member{}
	db.Where("id = ? and system_id = ? and deleted_at is null", memberid, system.ID).First(&member)
	if member.ID == "" {
		return errors.New("not found member request.")
	}
	role := models.Role{}
	db.Where("id = ? and deleted_at is null", member.RoleID).First(&role)
	if role.ID == 0 {
		return errors.New("not found role.")
	}
	richmenu := modelsLineAPI.RichMenu{}
	db.Where("line_oa_id = ? and status = ? and deleted_at is null", lineoa.ID, "afterregister"+role.RoleName).First(&richmenu)
	if richmenu.ID == 0 {
		return errors.New("not found richmenu.")
	}
	member.Approve = true
	tx := db.Begin()
	tx.Save(&member)
	err := SetLinkRichMenu(richmenu.RichID, lineoa.ChannelID, lineoa.ChannelSecret, member.LineID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func RejectRoleRequest(memberid uint, userid, systemid string) error {
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ? and deleted_at is null", systemid).First(&system)
	if system.ID == "" {
		return errors.New("not found system.")
	}
	admin := models.Admin{}
	db.Where("user_id = ? and system_id = ? and deleted_at is null", userid, system.ID).First(&admin)
	if admin.ID == 0 {
		return errors.New("you not admin.")
	}
	member := modelsMember.Member{}
	db.Where("id = ? and system_id = ? and deleted_at is null", memberid, system.ID).First(&member)
	if member.ID == "" {
		return errors.New("not found member request.")
	}
	lineoa := models.LineOA{}
	db.Where("system_id = ? and deleted_at is null", system.ID).First(&lineoa)
	if lineoa.ID == 0 {
		return errors.New("not found line oa.")
	}
	richmenu := modelsLineAPI.RichMenu{}
	db.Where("status = ? and line_oa_id = ? and deleted_at is null", "preregister", lineoa.ID).First(&richmenu)
	err := SetLinkRichMenu(richmenu.RichID, lineoa.ChannelID, lineoa.ChannelSecret, member.LineID)
	if err != nil {
		return err
	}
	db.Delete(&member)
	return nil
}

func DeleteRole(systemid, userid, roleid string) error {
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ? and deleted_at is null", systemid).First(&system)
	if system.ID == "" {
		return errors.New("system not found.")
	}
	admin := models.Admin{}
	db.Where("system_id = ? and user_id = ? and deleted_at is null", system.ID, userid).First(&admin)
	if admin.ID == 0 {
		return errors.New("you not admin.")
	}
	role := models.Role{}
	db.Where("id = ? and deleted_at is null", roleid).First(&role)
	if role.ID == 0 {
		return errors.New("role not found.")
	}
	targetgroup := modelsMember.TargetGroup{}
	db.Where("target_group_name = ? and system_id = ? and deleted_at is null", role.RoleName, system.ID).First(&targetgroup)
	tx := db.Begin()
	if targetgroup.ID != 0 {
		tx.Where("target_group_id = ? and deleted_at is null", targetgroup.ID).Delete(&modelsMember.MemberGroup{})
		tx.Where("target_group_name = ? and system_id = ? and deleted_at is null", role.RoleName, system.ID).Delete(&modelsMember.TargetGroup{})
	}
	members := []modelsMember.Member{}
	db.Where("role_id = ? and deleted_at is null", role.ID).Find(&members)
	lineoa := models.LineOA{}
	db.Where("system_id =? and deleted_at is null", system.ID).First(&lineoa)
	richmenu := modelsLineAPI.RichMenu{}
	db.Where("line_oa_id = ? and status = ? and deleted_at is null", lineoa.ID, "preregister").First(&richmenu)
	for _, member := range members {
		tx.Where("member_id = ? and deleted_at is null", member.ID).Delete(&modelsMember.MemberInterested{})
		SetLinkRichMenu(richmenu.RichID, lineoa.ChannelID, lineoa.ChannelSecret, member.LineID)
	}
	tx.Where("role_id = ? and deleted_at is null", role.ID).Delete(&modelsMember.Member{})
	tx.Where("id = ? and deleted_at is null", role.ID).Delete(&models.Role{})
	tx.Commit()
	return nil
}
