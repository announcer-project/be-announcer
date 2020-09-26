package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"errors"
)

func GetAllAdmin(systemid, userid string) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ? and deleted_at is null", systemid).First(&system)
	if system.ID == "" {
		return nil, errors.New("system not found.")
	}
	checkadmin := models.Admin{}
	db.Where("system_id = ? and user_id = ? and deleted_at is null", systemid, userid).First(&checkadmin)
	if checkadmin.ID == 0 {
		return nil, errors.New("you not admin.")
	}
	var adminDetail struct {
		Key      uint   `json:"key"`
		UserID   string `json:"userId"`
		Name     string `json:"name"`
		Position string `json:"position"`
	}
	var adminsDetail []struct {
		Key      uint   `json:"key"`
		UserID   string `json:"userId"`
		Name     string `json:"name"`
		Position string `json:"position"`
	}
	admin := models.Admin{}
	db.Where("system_id = ? and position = ? and deleted_at is null", system.ID, "admin").First(&admin)
	userAdmin := models.User{}
	db.Where("id = ? and deleted_at is null", admin.UserID).First(&userAdmin)
	adminDetail.Key = admin.ID
	adminDetail.Name = userAdmin.FName + " " + userAdmin.LName
	adminDetail.Position = admin.Position
	adminDetail.UserID = userAdmin.ID
	adminsDetail = append(adminsDetail, adminDetail)
	coAdmins := []models.Admin{}
	db.Where("system_id = ? and position = ? and deleted_at is null", system.ID, "co-admin").Find(&coAdmins)
	for _, coAdmin := range coAdmins {
		userCoAdmin := models.User{}
		db.Where("id = ? and deleted_at is null", coAdmin.UserID).First(&userCoAdmin)
		adminDetail.Key = coAdmin.ID
		adminDetail.Name = userCoAdmin.FName + " " + userCoAdmin.LName
		adminDetail.Position = coAdmin.Position
		adminDetail.UserID = userCoAdmin.ID
		adminsDetail = append(adminsDetail, adminDetail)
	}
	return adminsDetail, nil
}
