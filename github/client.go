package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

//
// The 3part service handlers, In OAuth2, the role is client
//
type TokenResponse struct {
	Token string `json:"access_token"`
	Type  string `json:"token_type"`
	Scope string `json:"scope"`
}

func GithubRequestAuth(accepted bool) func(c *gin.Context) {
	if accepted {
		return directReqAuth
	} else {
		return responseNavUI
	}
}

func directReqAuth(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently,
		"https://github.com/login/oauth/authorize?scope=user&client_id=a4c5ffb112e47562979e")
}

func responseNavUI(c *gin.Context) {
	c.HTML(http.StatusOK, "githubauthnav.tmpl", gin.H{
		"title": "Main website",
	})
}

func GithubEuthCallbackHandler(c *gin.Context) {
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

	var mRespBody TokenResponse
	json.NewDecoder(response.Body).Decode(&mRespBody)
	c.Writer.WriteString(getUserInfo(mRespBody.Token))
	c.Writer.WriteString("\n\n")
	c.Writer.WriteString("token is ")
	c.Writer.WriteString(mRespBody.Token)
}

func printTokenAccessResponseBodyStr(body io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	return buf.String()
}

func getUserInfo(token string) string {
	fmt.Printf("[%v]\n", token)
	request, _ := http.NewRequest("GET", "https://api.github.com/user",
		bytes.NewBuffer([]byte{}))
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("token %v", token))
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil || response.StatusCode != http.StatusOK {
		fmt.Println(err)
		return ""
	}
	return printTokenAccessResponseBodyStr(response.Body)
}
