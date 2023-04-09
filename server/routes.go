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
	s.InitAdminroutes(adminRouter)

}

func (s *Server) InitBaseRoutes(router *gin.RouterGroup) {
	router.POST("/register", middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.register.Handle(),
	}))
}

func (s *Server) InitUserRoutes(router *gin.RouterGroup) {
	router.POST("/:userID/finish-register", middleware.BindUserIDFromUri(), middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.finishRegister.Handle(),
	}))

	router.GET("/:userID", middleware.BindUserIDFromUri(), middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.getUserByID.Handle(),
	}))

	router.GET("", middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.getUserByNickname.Handle(),
	}))
}

func (s *Server) InitAdminroutes(router *gin.RouterGroup) {
	router.POST("/register", middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.adminRegister.Handle(),
	}))

	router.POST("/login", middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.adminLogin.Handle(),
	}))
}
