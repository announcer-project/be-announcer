package handlers

import (
	"be_nms/actions/repositories"
	"be_nms/models"
	"be_nms/models/modelsMember"
	"be_nms/models/modelsNews"
	"net/http"

	"github.com/labstack/echo/v4"
)

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
	targetgroups, err := repositories.GetAllTargetGroup(c.QueryParam("systemid"))
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
	newstypes, err := repositories.GetAllNewsType(c.QueryParam("systemid"), false)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	roles, err := repositories.GetAllRole(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	aboutSystem := AboutForLineRegister{NewsTypes: newstypes.([]modelsNews.NewsType), RolesUser: roles.([]models.Role)}
	return c.JSON(http.StatusOK, aboutSystem)
}

func GetAllSystems(c echo.Context) error {
	systems, err := repositories.GetAllsystems(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, systems)
}

func CreateSystem(c echo.Context) error {
	_, err := repositories.CreateSystem(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, "Create Success.")
}
