package main

import (
	"be_nms/database"
	"be_nms/routes"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	e := routes.Init()
	db := database.Open()
	defer db.Close()
	e.Logger.Fatal(e.Start(":1323"))
}
