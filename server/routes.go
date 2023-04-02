package server

import "github.com/fiufit/users/middleware"

func (s *Server) InitRoutes() {
	s.router.POST("/register", s.register.Handle())
	s.router.POST("/:userID/finish-register", middleware.BindUserIDFromUri(), s.finishRegister.Handle())
}
