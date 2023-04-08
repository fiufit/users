package server

import (
	"github.com/fiufit/users/middleware"
)

func (s *Server) InitRoutes() {
	baseRouter := s.router.Group("/:version")
	userRouter := baseRouter.Group("/users")

	baseRouter.POST("/register", middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.register.Handle(),
	}))

	userRouter.POST("/:userID/finish-register", middleware.BindUserIDFromUri(), middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.finishRegister.Handle(),
	}))

	userRouter.GET("/:userID", middleware.BindUserIDFromUri(), middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.getUserByID.Handle(),
	}))

	userRouter.GET("", middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.getUserByNickname.Handle(),
	}))
}
