package main

import (
	"github.com/fiufit/users/server"
	_ "github.com/lib/pq"
)

func main() {
	srv := server.NewServer()
	srv.InitRoutes()
	srv.Run()
}
