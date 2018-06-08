package main

import (
	"backend/OAuth"
	"backend/github"
	"github.com/gin-gonic/gin"
)

func main() {
	bootWithAuthService()
	//bootWithGithubClient()
}

func bootWithGithubClient() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	github.SetupGinRoute(router)
	router.Run(":14002")
}

func bootWithAuthService() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	aserver := OAuth.GetDefaultAuthServer()
	aserver.SetupGinRouter(router)
	router.Run(":14001")
}
