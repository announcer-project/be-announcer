package openapi

import (
	"be_nms/database"
	"be_nms/models/modelsNews"
	"errors"
)

func GetNewsByID(id string, systemid string) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	news := modelsNews.News{}
	db.Where("id = ? AND system_id = ?", id, systemid).First(&news)
	if news.ID == 0 {
		return nil, errors.New("Not have news")
	}
	return news, nil
}

func GetAllNews(status string, systemid string) interface{} {
	db := database.Open()
	defer db.Close()
	news := []modelsNews.News{}
	db.Where("status = ? AND system_id = ?", status, systemid).Find(&news)
	return news
}
