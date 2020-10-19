package handlers

import (
	"be_nms/actions/repositories"
	"be_nms/models"
	"be_nms/models/modelsMember"
	"be_nms/models/modelsNews"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetSystemByID(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	if c.Param("systemid") == "" {
		message.Message = "not have param."
		return c.JSON(400, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	system, err := repositories.GetSystemByID(tokens["user_id"].(string), c.Param("systemid"))
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	return c.JSON(200, system)
}

type About struct {
	News         []modelsNews.News          `json:"news"`
	NewsTypes    []modelsNews.NewsType      `json:"newstypes"`
	TatgetGroups []modelsMember.TargetGroup `json:"targetgroups"`
}

func GetAllAboutSystem(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	news, err := repositories.GetAllNews(tokens["user_id"].(string), c.QueryParam("systemid"), "publish")
	if err != nil {
		return err
	}
	newstypes, err := repositories.GetAllNewsType(c.QueryParam("systemid"), true)
	if err != nil {
		return err
	}
	targetgroups, err := repositories.GetAllTargetGroup(tokens["user_id"].(string), c.QueryParam("systemid"))
	if err != nil {
		return err
	}
	aboutSystem := About{News: news.([]modelsNews.News), NewsTypes: newstypes.([]modelsNews.NewsType), TatgetGroups: targetgroups.([]modelsMember.TargetGroup)}
	return c.JSON(http.StatusOK, aboutSystem)
}

type AboutForLineRegister struct {
	NewsTypes []modelsNews.NewsType `json:"newstypes"`
	RolesUser []models.Role         `json:"roles"`
}

func GetAboutSystemForLineRegister(c echo.Context) error {
	var message struct {
		Message string `json:"message"`
	}
	if c.QueryParam("systemid") == "" {
		message.Message = "not have query param."
		return c.JSON(400, message)
	}
	newstypes, err := repositories.GetAllNewsType(c.QueryParam("systemid"), false)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	roles, err := repositories.GetAllRole(c.QueryParam("systemid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	aboutSystem := AboutForLineRegister{NewsTypes: newstypes.([]modelsNews.NewsType), RolesUser: roles.([]models.Role)}
	return c.JSON(http.StatusOK, aboutSystem)
}

func GetAllSystems(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	systems, err := repositories.GetAllsystems(tokens["user_id"].(string))
	if err != nil {
		message.Message = "Server error."
		return c.JSON(500, message)
	}
	return c.JSON(http.StatusOK, systems)
}

func DeleteSystem(c echo.Context) error {
	var message struct {
		Message string `json:"message"`
	}
	authorization := c.Request().Header.Get("Authorization")
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	if c.Param("systemid") == "" {
		message.Message = "not have query param."
		return c.JSON(400, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	err := repositories.DeleteSystem(c.Param("systemid"), tokens["user_id"].(string))
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	message.Message = "delete success."
	return c.JSON(200, message)
}

func CreateSystem(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	var data struct {
		SystemProfile string
		Systemname    string
		NewsTypes     []string
		LineOA        struct {
			ChannelID          string
			ChannelAccessToken string
			RoleUsers          []struct {
				RoleName string
				Require  bool
			}
		}
	}
	if err := c.Bind(&data); err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	_, err := repositories.CreateSystem(tokens["user_id"].(string), data)
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	message.Message = "create system success."
	return c.JSON(http.StatusOK, message)
}
