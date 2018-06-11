package OAuth

import (
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/RangelReale/osin/example"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

type GinAuthServer struct {
	osinServer *osin.Server
}

func GetDefaultAuthServer() *GinAuthServer {
	fackeStorage := example.NewTestStorage()
	s := osin.NewServer(NewDefaultAuthServiceCfg(), fackeStorage)
	authServer := &GinAuthServer{}
	authServer.osinServer = s
	return authServer
}

func NewDefaultAuthServiceCfg() *osin.ServerConfig {
	sconfig := osin.NewServerConfig()
	sconfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN}
	sconfig.AllowedAccessTypes = osin.AllowedAccessType{osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN, osin.PASSWORD, osin.CLIENT_CREDENTIALS, osin.ASSERTION}
	sconfig.AllowGetAccessRequest = true
	sconfig.AllowClientSecretInParams = true
	return sconfig
}

func (s *GinAuthServer) AuthorizeReqHandler(c *gin.Context) {
	resp := s.osinServer.NewResponse()
	defer resp.Close()
	// Create AuthorizeRequest struct and return it's pointer for grant Authorization process by storage component,
	// it will write error info to Response stream when
	// client information invalid or other error(like authorize http request param format error) occur
	var authReq *osin.AuthorizeRequest = s.osinServer.HandleAuthorizeRequest(resp, c.Request)
	if cookieContainLoginedInfo(c.Request) == false && isSuccessMockLoginRequest(c) == false && authReq != nil {
		templateArg := gin.H{
			"clientId":     authReq.Client.GetId(),
			"state":        c.Query("state"),
			"scope":        c.Query("scope"),
			"redirect_uri": c.Query("redirect_uri")}
		c.HTML(http.StatusOK, "login.tmpl", templateArg)
	} else {
		if authReq != nil {
			authReq.UserData = struct{ Login string }{Login: "test"}
			authReq.Authorized = true
			s.osinServer.FinishAuthorizeRequest(resp, c.Request, authReq)
		}
		if resp.IsError && resp.InternalError != nil { // TODO: log error to database
			fmt.Printf("ERROR: %s\n", resp.InternalError)
		}
		if !resp.IsError {
			resp.Output["custom_parameter"] = 187723
		}
		osin.OutputJSON(resp, c.Writer, c.Request)
	}
}

func cookieContainLoginedInfo(r *http.Request) bool {
	_, err := r.Cookie("user_session")
	if err != nil {
		return false
	} else {
		// TODO: validate login info, and get user
		return false
	}
}

func isSuccessMockLoginRequest(c *gin.Context) bool {
	// TODO: This is test stub, please complete the logic.
	return c.Request.Method == http.MethodPost && c.Query("login") == "test" && c.Query("password") == "test"
}

func (s *GinAuthServer) AuthCodeReqHandler(c *gin.Context) {
	code := c.Query("code")
	c.Writer.Write([]byte("<html><body>"))
	c.Writer.Write([]byte("APP AUTH - CODE<br/>"))
	defer c.Writer.Write([]byte("</body></html>"))

	if code == "" {
		c.Writer.Write([]byte("Nothing to do"))
		return
	}
	jr := make(map[string]interface{})
	// build access code url
	aurl := fmt.Sprintf("/token?grant_type=authorization_code&client_id=1234&state=xyz&redirect_uri=%s&code=%s",
		url.QueryEscape("http://localhost:14000/appauth/code"), url.QueryEscape(code))

	if c.Query("doparse") == "1" { // if parse, download and parse json
		err := example.DownloadAccessToken(fmt.Sprintf("http://localhost:14000%s", aurl),
			&osin.BasicAuth{"1234", "aabbccdd"}, jr)
		if err != nil {
			c.Writer.Write([]byte(err.Error()))
			c.Writer.Write([]byte("<br/>"))
		}
	}

	if erd, ok := jr["error"]; ok { // show json error
		c.Writer.Write([]byte(fmt.Sprintf("ERROR: %s<br/>\n", erd)))
	}

	if at, ok := jr["access_token"]; ok { // show json access token
		c.Writer.Write([]byte(fmt.Sprintf("ACCESS TOKEN: %s<br/>\n", at)))
	}
	c.Writer.Write([]byte(fmt.Sprintf("FULL RESULT: %+v<br/>\n", jr)))
	// output links
	c.Writer.Write([]byte(fmt.Sprintf("<a href=\"%s\">Goto Token URL</a><br/>", aurl)))

	cururl := *c.Request.URL
	curq := cururl.Query()
	curq.Add("doparse", "1")
	cururl.RawQuery = curq.Encode()
	c.Writer.Write([]byte(fmt.Sprintf("<a href=\"%s\">Download Token</a><br/>", cururl.String())))
}

func (s *GinAuthServer) AuthClientIdxPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "authidx.tmpl", nil)
}

func (s *GinAuthServer) SetupGinRouter(router *gin.Engine) {
	router.Any("authorize", s.AuthorizeReqHandler)
	// Client Http API for receive the Auth Code from Authrize Server Redirect request.
	router.GET("appauth/code", s.AuthCodeReqHandler)
	router.GET("app", s.AuthClientIdxPageHandler)
}
