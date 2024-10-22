package server

import (
	_ "github.com/fiufit/users/docs"
	"github.com/fiufit/users/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) InitRoutes() {
	baseRouter := s.router.Group("/:version")
	userRouter := baseRouter.Group("/users")
	adminRouter := baseRouter.Group("/admin")

	s.InitDocRoutes(baseRouter)
	s.InitUserRoutes(userRouter)
	s.InitAdminRoutes(adminRouter)

}

func (s *Server) InitDocRoutes(router *gin.RouterGroup) {
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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

	router.POST("/:userID/followers", middleware.BindUserIDFromUri(), middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.followUser.Handle(),
	}))

	router.DELETE("/:userID/followers/:followerID", middleware.BindUserIDFromUri(), middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.unfollowUser.Handle(),
	}))

	router.GET("/:userID/followers", middleware.BindUserIDFromUri(), middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.getUserFollowers.Handle(),
	}))

	router.GET("/:userID/followed", middleware.BindUserIDFromUri(), middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.getFollowedUsers.Handle(),
	}))

	router.GET("/:userID/closest", middleware.BindUserIDFromUri(), middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.getClosestUsers.Handle(),
	}))

	router.POST("/:userID/enable", middleware.BindUserIDFromUri(), middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.enableUser.Handle(),
	}))

	router.DELETE("/:userID/disable", middleware.BindUserIDFromUri(), middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.disableUser.Handle(),
	}))

	router.POST("/password-recover", middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.notifyPasswordRecover.Handle(),
	}))

	router.POST("/login", middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.notifyUserLogin.Handle(),
	}))

	router.POST("/:userID/verification/send", middleware.BindUserIDFromUri(), middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.sendVerificationPin.Handle(),
	}))

	router.POST("/:userID/verification/verify", middleware.BindUserIDFromUri(), middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.verifyUser.Handle(),
	}))

	router.POST("/certifications", middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.createCert.Handle(),
	}))

	router.GET("/certifications", middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.getCert.Handle(),
	}))

	router.PUT("/certifications/:certificationID", middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.updateCert.Handle(),
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
