package main

import (
	"backend/OAuth"
	"backend/github"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/prometheus/common/log"
	"github.com/RangelReale/osin"
	"backend/models"
	"database/sql"
	"github.com/felipeweb/osin-mysql"
)

type Server struct {
	WebServer *gin.Engine
	AuthServer *OAuth.GinAuthServer
	Db *gorm.DB
}

func NewServer() (*Server, error) {
	url := "user:password@tcp(host:3306)/dbname?parseTime=true"
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	store := mysql.New(db,"osin_")
	if err := store.CreateSchemas(); err != nil {
		log.Warnf("Error creating osin schemas: %v", err.Error())
	}

	authServer := OAuth.GetDefaultAuthServer(store)

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

func NewAuthStorage() osin.Storage {
	return models.NewTAuthMySQLStorage()
}

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.Run())
}
