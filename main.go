package main

import (
	_ "web_server/infrastructure"
	"web_server/router"
)

func main() {
	router.RunServer()
}
