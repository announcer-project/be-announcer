package main

import (
	"be_nms/database"
	"be_nms/routes"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := routes.Init()
	db := database.Open()
	defer db.Close()
	database.Migration(db)
	// database.SetData(db)
	e.Use(middleware.CORS())
	e.Logger.Fatal(e.Start(":8080"))
}
