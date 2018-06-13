package main

import (
	"backend/handlers"
	"database/sql"
	"errors"
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/felipeweb/osin-mysql"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/prometheus/common/log"
	"os"
)

type Server struct {
	WebServer   *gin.Engine
	WebDb       *gorm.DB
	AuthServer  *osin.Server
	AuthStorage *mysql.Storage
}

func NewServer() (*Server, error) {
	storage, err := newAuthStorage()
	if err != nil {
		return nil, err
	}

	webServer := newWebServer()
	authServer := osin.NewServer(newDefaultAuthServiceCfg(), storage)

	server := &Server{
		WebServer:   webServer,
		AuthServer:  authServer,
		AuthStorage: storage,
	}

	initTestClient(storage)

	authorizeHandler := handlers.AuthorizeReqHandler(authServer)
	accesstokenHandler := handlers.AccessTokenHandler(authServer)
	tokeninfoHandler := handlers.TokenInfoHandler(authServer)

	oauthGroup := webServer.Group("/oauth")
	oauthGroup.GET("/authorize", authorizeHandler).POST("/authorize", authorizeHandler).
		GET("/token", accesstokenHandler).POST("/token", accesstokenHandler).
		GET("/api/token/info", tokeninfoHandler).POST("/api/token/info", tokeninfoHandler)

	return server, nil
}

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.WebServer.Run())
}

func newAuthStorage() (*mysql.Storage, error) {
	osinDbUsername, ok := os.LookupEnv("OSIN_DB_USERNAME")
	if !ok || len(osinDbUsername) == 0 {
		return nil, errors.New("No OSIN_DB_USERNAME env is provided!")
	}

	osinDbPassword, ok := os.LookupEnv("OSIN_DB_PASSWORD")
	if !ok || len(osinDbPassword) == 0 {
		log.Warn("No OSIN_DB_PASSWORD env is provided. If it is, your database server is NOT SECURE!")
	}

	osinDbHost, ok := os.LookupEnv("OSIN_DB_HOST")
	if !ok || len(osinDbHost) == 0 {
		return nil, errors.New("No OSIN_DB_HOST env is provided!")
	}

	osinDbPort, ok := os.LookupEnv("OSIN_DB_PORT")
	if !ok || len(osinDbPort) == 0 {
		log.Warn("No OSIN_DB_PORT env is provided and default value 3306 is used.")
		osinDbPort = "3306"
	}

	osinDbDatabse, ok := os.LookupEnv("OSIN_DB_DATABASE")
	if !ok || len(osinDbDatabse) == 0 {
		return nil, errors.New("No OSIN_DB_DATABASE env is provided!")
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		osinDbUsername,
		osinDbPassword,
		osinDbHost,
		osinDbPort,
		osinDbDatabse,
	))
	if err != nil {
		return nil, err
	}

	osinTablePrefix, ok := os.LookupEnv("OSIN_TABLE_PREFIX")
	if !ok || len(osinTablePrefix) == 0 {
		log.Warn(`No OSIN_TABLE_PREFIX env is provided and default value "osin_" is used.`)
		osinTablePrefix = "osin_"
	}

	store := mysql.New(db, osinTablePrefix)
	if err := store.CreateSchemas(); err != nil {
		log.Warnf("Error creating osin schemas: %v", err.Error())
	}

	return store, nil
}

func newWebServer() *gin.Engine {
	ginMode, ok := os.LookupEnv(gin.ENV_GIN_MODE)
	if !ok || len(ginMode) == 0 {
		gin.SetMode(gin.ReleaseMode)
	}
	server := gin.Default()
	server.LoadHTMLGlob("templates/*")
	return server
}

func initTestClient(storage *mysql.Storage) {
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
			UserData:    "",
		},
	)
	if createErr == nil {
		log.Debug("default client did update.")
	} else {
		log.Warn("Insert default client error.", createErr.Error())
	}
}

func newDefaultAuthServiceCfg() *osin.ServerConfig {
	sconfig := osin.NewServerConfig()
	sconfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN}
	sconfig.AllowedAccessTypes = osin.AllowedAccessType{osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN, osin.PASSWORD, osin.CLIENT_CREDENTIALS, osin.ASSERTION}
	sconfig.AllowGetAccessRequest = true
	sconfig.AllowClientSecretInParams = true
	return sconfig
}
