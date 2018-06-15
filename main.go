package main

import (
	"backend/handlers"
	"backend/models"
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

func NewServer() *Server {
	authStorage := newAuthStorage()
	authServer := osin.NewServer(newAuthServerConfig(), authStorage)
	webDb := newWebDb()
	webServer := newWebServer()

	server := &Server{
		WebServer:   webServer,
		WebDb:       webDb,
		AuthServer:  authServer,
		AuthStorage: authStorage,
	}

	initTestClient(authStorage)
	registerAuthRoutes(webServer, authServer, authStorage)
	registerWebRoutes(webServer, webDb)

	return server
}

func newAuthStorage() *mysql.Storage {
	osinDbUsername, ok := os.LookupEnv("OSIN_DB_USERNAME")
	if !ok || len(osinDbUsername) == 0 {
		panic(errors.New("no OSIN_DB_USERNAME env is provided"))
	}

	osinDbPassword, ok := os.LookupEnv("OSIN_DB_PASSWORD")
	if !ok || len(osinDbPassword) == 0 {
		log.Warn("No OSIN_DB_PASSWORD env is provided. If it is, your database server is NOT SECURE!")
	}

	osinDbHost, ok := os.LookupEnv("OSIN_DB_HOST")
	if !ok || len(osinDbHost) == 0 {
		panic(errors.New("no OSIN_DB_HOST env is provided"))
	}

	osinDbPort, ok := os.LookupEnv("OSIN_DB_PORT")
	if !ok || len(osinDbPort) == 0 {
		log.Warn("No OSIN_DB_PORT env is provided and default value 3306 is used.")
		osinDbPort = "3306"
	}

	osinDbDatabse, ok := os.LookupEnv("OSIN_DB_DATABASE")
	if !ok || len(osinDbDatabse) == 0 {
		panic(errors.New("no OSIN_DB_DATABASE env is provided"))
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		osinDbUsername,
		osinDbPassword,
		osinDbHost,
		osinDbPort,
		osinDbDatabse,
	))
	if err != nil {
		panic(err)
	}

	osinTablePrefix, ok := os.LookupEnv("OSIN_TABLE_PREFIX")
	if !ok || len(osinTablePrefix) == 0 {
		log.Warn(`No OSIN_TABLE_PREFIX env is provided and default value "osin_" is used.`)
		osinTablePrefix = "osin_"
	}

	store := mysql.New(db, osinTablePrefix)
	if err := store.CreateSchemas(); err != nil {
		panic(err)
	}

	return store
}

func newWebServer() *gin.Engine {
	ginMode, ok := os.LookupEnv(gin.ENV_GIN_MODE)
	if !ok || len(ginMode) == 0 {
		gin.SetMode(gin.ReleaseMode)
	}

	staticPathFromEnv := func() string {
		osinStaticRs, ok := os.LookupEnv("OSIN_STATIC_RESOURCES")
		if !ok || len(osinStaticRs) == 0 {
			osinStaticRs = "./statics"
		}

		if fi, err := os.Stat(osinStaticRs); err != nil {
			if os.IsNotExist(err) {
				os.Mkdir(osinStaticRs, os.ModeDir)
				return osinStaticRs
			} else {
				panic(err)
			}
		} else if fi.Mode().IsDir() == false {
			panic(errors.New("OSIN_STATIC_RESOURCES file is exist, but it is not a folder!"))
		} else {
			return osinStaticRs
		}
	}()

	server := gin.Default()
	server.Static("/assets", staticPathFromEnv)
	server.LoadHTMLGlob("templates/*")
	return server
}

func initTestClient(storage *mysql.Storage) {
	testClientId := "5678"

	_, err := storage.GetClient(testClientId)
	if err != nil {
		if err != osin.ErrNotFound {
			panic(err)
		}
	} else {
		if err = storage.RemoveClient(testClientId); err != nil {
			panic(err)
		}
	}

	if err := storage.CreateClient(
		&osin.DefaultClient{
			Id:          testClientId,
			Secret:      "9527abcdefg0039",
			RedirectUri: "http://localhost:8091/oauth/callback",
			UserData:    "",
		},
	); err != nil {
		panic(err)
	}
}

func registerAuthRoutes(webServer *gin.Engine, authServer *osin.Server, authOperationStorage *mysql.Storage) {
	authorizeHandler := handlers.AuthorizeReqHandler(authServer)
	accesstokenHandler := handlers.AccessTokenHandler(authServer)
	tokeninfoHandler := handlers.TokenInfoHandler(authServer)

	webServer.Group("/oauth").GET("/authorize", authorizeHandler).POST("/authorize", authorizeHandler).
		GET("/token", accesstokenHandler).POST("/token", accesstokenHandler).
		GET("/api/token/info", tokeninfoHandler).POST("/api/token/info", tokeninfoHandler).
		POST("/register/client", handlers.RegistClientApp(authOperationStorage))
}

func newAuthServerConfig() *osin.ServerConfig {
	config := osin.NewServerConfig()
	config.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN}
	config.AllowedAccessTypes = osin.AllowedAccessType{osin.AUTHORIZATION_CODE, osin.REFRESH_TOKEN, osin.PASSWORD, osin.CLIENT_CREDENTIALS, osin.ASSERTION}
	config.AllowGetAccessRequest = true
	config.AllowClientSecretInParams = true
	config.AuthorizationExpiration = 250
	config.AccessExpiration = 3600
	return config
}

func newWebDb() *gorm.DB {
	webDbUsername, ok := os.LookupEnv("WEB_DB_USERNAME")
	if !ok || len(webDbUsername) == 0 {
		panic(errors.New("no WEB_DB_USERNAME env is provided"))
	}

	webDbPassword, ok := os.LookupEnv("WEB_DB_PASSWORD")
	if !ok || len(webDbPassword) == 0 {
		log.Warn("No WEB_DB_PASSWORD env is provided. If it is, your database server is NOT SECURE!")
	}

	webDbHost, ok := os.LookupEnv("WEB_DB_HOST")
	if !ok || len(webDbHost) == 0 {
		panic(errors.New("no WEB_DB_HOST env is provided"))
	}

	webDbPort, ok := os.LookupEnv("WEB_DB_PORT")
	if !ok || len(webDbPort) == 0 {
		log.Warn("No WEB_DB_PORT env is provided and default value 3306 is used.")
		webDbPort = "3306"
	}

	webDbDatabse, ok := os.LookupEnv("WEB_DB_DATABASE")
	if !ok || len(webDbDatabse) == 0 {
		panic(errors.New("no WEB_DB_DATABASE env is provided"))
	}

	db, err := gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		webDbUsername,
		webDbPassword,
		webDbHost,
		webDbPort,
		webDbDatabse,
	))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{})

	return db
}

func registerWebRoutes(webServer *gin.Engine, webDb *gorm.DB) {
	g := webServer.Group("/member")

	vcodeHandler := handlers.VCodeHandler(webDb)
	loginHandler := handlers.VCodeLoginHandler()

	g.POST("/vcode", vcodeHandler)
	g.POST("/login", loginHandler)
}

func main() {
	log.Fatal(NewServer().WebServer.Run())
}
