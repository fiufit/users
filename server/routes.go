package server

import "github.com/fiufit/users/middleware"

func (s *Server) InitRoutes() {
	baseRouter := s.router.Group("/:version")
	userRouter := baseRouter.Group("/users")

	baseRouter.POST("/register", s.register.Handle())
	userRouter.POST("/:userID/finish-register", middleware.BindUserIDFromUri(), s.finishRegister.Handle())
}
