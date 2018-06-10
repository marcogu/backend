package main

import (
	"backend/OAuth"
	"backend/github"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	prepareWithAuthService(server)
	prepareWithGithubClient(server)
	server.Run(":14000")
}

func prepareWithGithubClient(route *gin.Engine) {
	route.LoadHTMLGlob("templates/*")
	github.SetupGinRoute(route)
}

func prepareWithAuthService(route *gin.Engine) {
	route.LoadHTMLGlob("templates/*")
	aserver := OAuth.GetDefaultAuthServer()
	aserver.SetupGinRouter(route)
}
