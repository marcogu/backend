package OAuth

import (
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/felipeweb/osin-mysql"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"net/http"
	"net/url"
)

type GinAuthServer struct {
	osinServer *osin.Server
}

func GetDefaultAuthServer(storage *mysql.Storage) *GinAuthServer {
	initDbData(storage)
	s := osin.NewServer(NewDefaultAuthServiceCfg(), storage)
	authServer := &GinAuthServer{}
	authServer.osinServer = s
	return authServer
}

func initDbData(storage *mysql.Storage) {
	testClientId := "5678"
	testClient, err := storage.GetClient(testClientId)
	if testClient != nil {
		storage.RemoveClient(testClientId)
	} else if err != nil && err != osin.ErrNotFound {
		log.Warn("Try add default client error, stop update default application client:", err.Error())
		return
	}
	createErr := storage.CreateClient(
		&osin.DefaultClient{
			Id:          testClientId,
			Secret:      "9527abcdefg0039",
			RedirectUri: "http://localhost:8091/oauth/callback",
			UserData:    ""},
	)
	if createErr == nil {
		log.Debug("default client did update.")
	} else {
		log.Warn("Insert default client error.", createErr.Error())
	}
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
	var authReq *osin.AuthorizeRequest = s.osinServer.HandleAuthorizeRequest(resp, c.Request)
	if cookieContainLoginedInfo(c.Request) == false && isSuccessMockLoginRequest(c) == false && authReq != nil {
		templateArg := gin.H{
			"clientId":      authReq.Client.GetId(),
			"state":         c.Query("state"),
			"scope":         c.Query("scope"),
			"response_type": c.Query("response_type"),
			"redirect_uri":  c.Query("redirect_uri")}
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
	succLoginRs := c.Request.Method == http.MethodPost &&
		c.Request.Form.Get("login") == "test" &&
		c.Request.Form.Get("password") == "test"
	return succLoginRs
}

func (s *GinAuthServer) AuthCodeReceiveHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.HTML(http.StatusOK, "token.tmpl", gin.H{"noneErr": false, "errInfo": "Receive Code is empty, Nothing to do"})
	} else {
		accessTokenRequestUrl := fmt.Sprintf("/token?grant_type=authorization_code&client_id=1234&client_secret=aabbccdd&state=xyz&redirect_uri=%s&code=%s",
			url.QueryEscape("http://localhost:14000/appauth/code"), url.QueryEscape(code))
		htmlParam := gin.H{
			"gotoTokenUrl": accessTokenRequestUrl,
			"code":         code,
			"clientId":     "1234",
			"secury":       "aabbccdd",
			"noneErr":      true,
		}
		c.HTML(http.StatusOK, "token.tmpl", htmlParam)
	}
}

func (s *GinAuthServer) AuthClientIdxPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "authidx.tmpl", nil)
}

func (s *GinAuthServer) AccessTokenHandler(c *gin.Context) {
	resp := s.osinServer.NewResponse()
	defer resp.Close()

	if ar := s.osinServer.HandleAccessRequest(resp, c.Request); ar != nil {
		ar.Authorized = true
		s.osinServer.FinishAccessRequest(resp, c.Request, ar)
	}
	if resp.IsError && resp.InternalError != nil {
		fmt.Printf("ERROR: %s\n", resp.InternalError)
	}
	osin.OutputJSON(resp, c.Writer, c.Request)
}

func (s *GinAuthServer) TokenInfoHandler(c *gin.Context) {
	resp := s.osinServer.NewResponse()
	defer resp.Close()
	if ir := s.osinServer.HandleInfoRequest(resp, c.Request); ir != nil {
		s.osinServer.FinishInfoRequest(resp, c.Request, ir)
	}
	osin.OutputJSON(resp, c.Writer, c.Request)
}

func (s *GinAuthServer) SetupGinRouter(router *gin.Engine) {
	router.Any("authorize", s.AuthorizeReqHandler)       // role : Authorize server
	router.GET("appauth/code", s.AuthCodeReceiveHandler) // role : Client server
	router.GET("app", s.AuthClientIdxPageHandler)        // role : Client (Resources) server, Index page.
	router.Any("/token", s.AccessTokenHandler)           // role : Authorize server
	router.Any("/api/token/info", s.TokenInfoHandler)
}
