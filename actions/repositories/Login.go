package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"errors"
)

func GetUserBySocialId(SocialID, Social string) (interface{}, error) {
	db := database.Open()
	user := models.User{}
	if Social == "line" {
		db.First(&user, "line_id = ?", SocialID)
		if user.ID == "" {
			return nil, errors.New("You don't register.")
		}
	} else if Social == "facebook" {
		db.First(&user, "facebook_id = ?", SocialID)
		if user.ID == "" {
			return nil, errors.New("You don't register.")
		}
	}
	return user, nil
}
