package main

import (
	"backend/github"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/idx", github.GithubRequestAuth(true))
	router.GET("/oauth/callback", github.GithubEuthCallbackHandler)

	router.Run(":14002")
}
