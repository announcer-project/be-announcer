package main

import (
	"be_nms/routes"
)

func main() {
	e := routes.Init()
	e.Logger.Fatal(e.Start(":1323"))
}
