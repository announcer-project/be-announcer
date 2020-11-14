package handlers

import (
	"be_nms/actions/repositories"
	"be_nms/models"
	"be_nms/models/modelsMember"
	"be_nms/models/modelsNews"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateMember(c echo.Context) error {
	var data struct {
		FName          string
		LName          string
		RoleID         int
		NewsInterested []modelsNews.NewsType
		SystemID       string
		Line           string
	}
	var message struct {
		Message string `json:"message"`
	}
	if err := c.Bind(&data); err != nil {
		message.Message = "server error."
		return c.JSON(500, message)
	}
	log.Print(data)
	err := repositories.RegisterGetNews(data)
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	message.Message = "register for get news success."
	return c.JSON(http.StatusOK, message)
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
	return c.JSON(http.StatusOK, members)
}

func GetMemberByLineID(c echo.Context) error {
	var message struct {
		Message string `json:"message"`
	}
	lineid := c.Param("lineid")
	if lineid == "" {
		message.Message = "not have query param."
		return c.JSON(400, message)
	}
	member, err := repositories.GetMemberByLineID(lineid)
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	role, err := repositories.GetRoleByID(member.(modelsMember.Member).RoleID)
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	var MemberDetail struct {
		Member modelsMember.Member `json:"member"`
		Role   models.Role         `json:"role"`
	}
	MemberDetail.Member = member.(modelsMember.Member)
	MemberDetail.Role = role.(models.Role)
	return c.JSON(200, MemberDetail)
}

func UpdateMemberName(c echo.Context) error {
	var message struct {
		Message string `json:"message"`
	}
	memberid := c.Param("memberid")
	if memberid == "" {
		message.Message = "not have param."
		return c.JSON(400, message)
	}
	var data struct {
		FName string
		LName string
	}
	if err := c.Bind(&data); err != nil {
		message.Message = "server error."
		return c.JSON(500, message)
	}
	err := repositories.UpdateMemberName(data.FName, data.LName, memberid)
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	message.Message = "update success"
	return c.JSON(200, message)
}

func UpdateMemberRole(c echo.Context) error {
	var message struct {
		Message string `json:"message"`
	}
	memberid := c.Param("memberid")
	if memberid == "" {
		message.Message = "not have param."
		return c.JSON(400, message)
	}
	var data struct {
		RoleID string
	}
	if err := c.Bind(&data); err != nil {
		message.Message = "server error."
		return c.JSON(500, message)
	}
	log.Print(data)
	err := repositories.UpdateMemberRole(data.RoleID, memberid)
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	message.Message = "update success"
	return c.JSON(200, message)
}

func UpdateMemberNewstype(c echo.Context) error {
	var message struct {
		Message string `json:"message"`
	}
	memberid := c.Param("memberid")
	if memberid == "" {
		message.Message = "not have param."
		return c.JSON(400, message)
	}
	var data struct {
		Newstypes []struct {
			Newstype   modelsNews.NewsType
			Interested bool
		}
	}
	if err := c.Bind(&data); err != nil {
		message.Message = "server error."
		return c.JSON(500, message)
	}
	log.Print(data)
	err := repositories.UpdateMemberNewstype(memberid, data.Newstypes)
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	message.Message = "update success"
	return c.JSON(200, message)
}
