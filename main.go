package main

import (
	"backend/OAuth"
	"backend/github"
	"database/sql"
	"errors"
	"fmt"
	"github.com/felipeweb/osin-mysql"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/prometheus/common/log"
	"os"
)

type Server struct {
	WebServer  *gin.Engine
	AuthServer *OAuth.GinAuthServer
	Db         *gorm.DB
}

func NewServer() (*Server, error) {
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
		log.Warn("No OSIN_TABLE_PREFIX env is provided and default value \"osin_\" is used.")
		osinTablePrefix = "osin_"
	}

	store := mysql.New(db, osinTablePrefix)
	if err := store.CreateSchemas(); err != nil {
		log.Warnf("Error creating osin schemas: %v", err.Error())
	}

	authServer := OAuth.GetDefaultAuthServer(store)

	ginMode, ok := os.LookupEnv(gin.ENV_GIN_MODE)
	if !ok || len(ginMode) == 0 {
		gin.SetMode(gin.ReleaseMode)
	}

	server := &Server{
		WebServer:  gin.Default(),
		AuthServer: authServer,
	}

	server.WebServer.LoadHTMLGlob("templates/*")
	server.AuthServer.SetupGinRouter(server.WebServer)

	// For github client tests.
	server.WebServer.LoadHTMLGlob("templates/*")
	github.SetupGinRoute(server.WebServer)

	return server, nil
}

func (s *Server) Run() error {
	return s.WebServer.Run(":14000")
}

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.Run())
}
