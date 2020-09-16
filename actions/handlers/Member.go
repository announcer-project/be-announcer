package handlers

import (
	"be_nms/actions/repositories"
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsMember"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateMember(c echo.Context) error {
	_, err := repositories.RegisterGetNews(c)
	if err != nil {
		log.Print("Error : ", err)
		return c.JSON(401, err)
	}
	return c.String(http.StatusOK, "OK!")
}

func GetAllMember(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	if c.QueryParam("systemid") == "" {
		message.Message = "not have query param."
		return c.JSON(400, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	members, err := repositories.GetAllMember(tokens["user_id"].(string), c.QueryParam("systemid"))
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	membersArr := members.([]modelsMember.Member)
	db := database.Open()
	defer db.Close()
	type Members struct {
		Member modelsMember.Member `json:"member"`
		User   models.User         `json:"user"`
	}
	membersRes := []Members{}
	for _, member := range membersArr {
		user := models.User{}
		db.Where("id = ?", member.UserID).First(&user)
		m := Members{Member: member, User: user}
		membersRes = append(membersRes, m)
	}
	return c.JSON(http.StatusOK, membersRes)
}
