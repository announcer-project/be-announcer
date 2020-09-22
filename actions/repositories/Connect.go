package repositories

import (
	"be_nms/database"
	"be_nms/models"
)

func CheclConnectLineOA(systemid string) bool {
	db := database.Open()
	lineoa := models.LineOA{}
	db.Where("system_id = ?", systemid).First(&lineoa)
	if lineoa.ID == 0 {
		return false
	}
	return true
}
