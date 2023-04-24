package server

import (
	"github.com/fiufit/users/middleware"
	"github.com/gin-gonic/gin"
)

func (s *Server) InitRoutes() {
	baseRouter := s.router.Group("/:version")
	userRouter := baseRouter.Group("/users")
	adminRouter := baseRouter.Group("/admin")

	s.InitUserRoutes(userRouter)
	s.InitAdminRoutes(adminRouter)

}

func (s *Server) InitUserRoutes(router *gin.RouterGroup) {
	router.POST("/register", middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.register.Handle(),
	}))

	router.POST("/:userID/finish-register", middleware.BindUserIDFromUri(), middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.finishRegister.Handle(),
	}))

	router.GET("/:userID", middleware.BindUserIDFromUri(), middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.getUserByID.Handle(),
	}))

	router.PATCH("/:userID", middleware.BindUserIDFromUri(), middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.updateUser.Handle(),
	}))

	router.DELETE("/:userID", middleware.BindUserIDFromUri(), middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.deleteUser.Handle(),
	}))

	router.GET("", middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.getUsers.Handle(),
	}))

}

func (s *Server) InitAdminRoutes(router *gin.RouterGroup) {
	router.POST("/register", middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.adminRegister.Handle(),
	}))

	router.POST("/login", middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.adminLogin.Handle(),
	}))
}
