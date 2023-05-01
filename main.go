package main

import (
	"github.com/fiufit/users/server"
	_ "github.com/lib/pq"
)

// @title           Fiufit Users API
// @version        	dev
// @description     Fiufit's Users service documentation. This service manages accounts, profiles, admin authentication, etc.

// @host      fiufit-users.fly.dev
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	srv := server.NewServer()
	srv.InitRoutes()
	srv.Run()
}
