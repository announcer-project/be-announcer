package routes

import (
	"be_nms/actions/handlers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Init() *echo.Echo {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to Announcer")
	})
	//Account
	e.POST("/login", handlers.Login)
	e.POST("/login/line", handlers.LineLogin)
	e.POST("/register", handlers.Register)
	e.POST("/register/sendotp", handlers.SendOTP)
	e.POST("/register/checkuser", handlers.CheckUserByEmail)
	e.POST("/register/checkuserbylineid", handlers.CheckUserByLineID)
	e.POST("/register/connectsocial", handlers.ConnectSocialWithAccount)
	//System
	e.GET("/system/all", handlers.GetAllSystems)
	e.POST("/system/create", handlers.CreateSystem)
	//NewsManagement
	e.GET("/aboutsystem", handlers.GetAllAboutSystem)
	e.GET("/news/all", handlers.GetAllNewsByClassify)
	e.GET("/news/all/draft", handlers.GetAllNewsDraft)
	e.GET("/news/all/publish", handlers.GetAllNewsPublish)
	e.POST("/news/create", handlers.CreateNews)
	e.GET("/news/publish/:id", handlers.GetNewsByID)
	e.GET("/news/draft/:id", handlers.GetNewsDraftByID)
	e.POST("/news/newstype/create", handlers.CreateNewsType)
	e.POST("/news/newstype/delete", handlers.DeleteNewsType)
	e.GET("/news/newstype/all", handlers.GetAlNewsType)
	//TargetGroup
	e.POST("/targetgroup/create", handlers.CreateTargetGroup)
	e.GET("/targetgroup/all", handlers.GetAllTargetGroup)
	//Member
	e.GET("/member/all", handlers.GetAllMember)
	//Role
	e.POST("/role/create", handlers.CreateRole)
	e.GET("/role/all", handlers.GetAllRole)
	//Social
	e.POST("/webhooklineoa", handlers.WebhookLineOA)
	e.GET("/connect/line/check", handlers.CheckConnectLineOA)
	e.DELETE("/connect/line/:systemid", handlers.DisconnectLinaOA)
	e.POST("/connect/line", handlers.ConenctLineOA)
	//Broadcast
	e.GET("/broadcast/line/aboutsystem", handlers.GetAboutLineBroadcast)
	e.POST("/broadcast/line", handlers.BroadcastNewsToLine)
	//Liff
	e.POST("/line/register", handlers.CreateMember)
	e.GET("/line/register/aboutsystem", handlers.GetAboutSystemForLineRegister)

	//Line API Richmenu
	e.POST("/webhook/:hookid", handlers.WebhookLineOA)
	// e.GET("/richmenu/setdefaultregister", handlers.SetDefaultRichMenuRegister)
	return e
}
