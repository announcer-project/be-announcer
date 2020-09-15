package repositories

import (
	"be_nms/database"
	"be_nms/models"
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
