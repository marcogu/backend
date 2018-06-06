package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	didAccept := true

	router.GET("/idx", githubRequestAuth(didAccept))
	router.GET("/oauth/callback", githubEuthCallbackHandler)

	router.Run(":14002")
}

func githubRequestAuth(accepted bool) func(c *gin.Context) {
	if accepted {
		return func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently,
				"https://github.com/login/oauth/authorize?scope=user:email&client_id=a4c5ffb112e47562979e")
		}
	} else {
		return func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"title": "Main website",
			})
		}
	}
}

// http://localhost:14002/idx
func githubEuthCallbackHandler(c *gin.Context) {
	code := c.DefaultQuery("code", "none")
	fmt.Println(code)

	mBody := map[string]interface{}{
		"client_id":     "a4c5ffb112e47562979e",
		"client_secret": "d0c47635f20d5447ec67090b32f6e73ca0cce50d",
		"code":          code,
		"accept":        "application/json",
	}
	jobdy, _ := json.Marshal(mBody)
	request, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(jobdy))
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil || response.StatusCode != http.StatusOK {
		fmt.Println(err)
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	fmt.Println("response is2:", buf.String())
	c.Writer.WriteString("---------------")
}
