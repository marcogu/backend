package handlers

import (
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginForm struct {
	Login            string `form:"login"`
	Pwd              string `form:"password"`
	ClientId         string `form:"client_id"`
	AuthorizeGateway string `form:"authorizeGateway"`
	State            string `form:"state"`
	Scope            string `form:"scope"`
	Response_type    string `form:"response_type"`
	Redirect_uri     string `form:"redirect_uri"`
}

func AuthorizeReqHandler(s *osin.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := s.NewResponse()
		defer resp.Close()

		var loginForm LoginForm
		c.ShouldBind(&loginForm)

		var authReq = s.HandleAuthorizeRequest(resp, c.Request)
		if cookieContainLoginedInfo(c.Request) == false && isSuccessMockLoginRequest(&loginForm) == false && authReq != nil {
			loginForm.ClientId = authReq.Client.GetId()
			loginForm.AuthorizeGateway = c.Request.URL.Path
			c.HTML(http.StatusOK, "login.tmpl", loginForm)
		} else {
			if authReq != nil {
				authReq.UserData = struct{ Login string }{Login: "test"}
				authReq.Authorized = true
				s.FinishAuthorizeRequest(resp, c.Request, authReq)
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

func isSuccessMockLoginRequest(form *LoginForm) bool {
	// TODO: This is test stub, please complete the logic.
	succLoginRs := form.Login == "test" && form.Pwd == "test"
	return succLoginRs
}

func AccessTokenHandler(s *osin.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := s.NewResponse()
		defer resp.Close()

		if ar := s.HandleAccessRequest(resp, c.Request); ar != nil {
			ar.Authorized = true
			s.FinishAccessRequest(resp, c.Request, ar)
		}
		if resp.IsError && resp.InternalError != nil {
			fmt.Printf("ERROR: %s\n", resp.InternalError)
		}
		osin.OutputJSON(resp, c.Writer, c.Request)
	}

}

func TokenInfoHandler(s *osin.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := s.NewResponse()
		defer resp.Close()
		if ir := s.HandleInfoRequest(resp, c.Request); ir != nil {
			s.FinishInfoRequest(resp, c.Request, ir)
		}
		osin.OutputJSON(resp, c.Writer, c.Request)
	}
}
