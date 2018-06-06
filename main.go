package main

import (
	"backend/github"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	github.SetupGinRoute(router)
	router.Run(":14002")
}
