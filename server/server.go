package server

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/fiufit/users/database"
	"github.com/fiufit/users/handlers"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories"
	"github.com/fiufit/users/usecases/accounts"
	"github.com/fiufit/users/usecases/users"
	"github.com/fiufit/users/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	router *gin.Engine

	register         handlers.Register
	finishRegister   handlers.FinishRegister
	adminRegister    handlers.AdminRegister
	adminLogin       handlers.AdminLogin
	getUserByID      handlers.GetUserByID
	getUsers         handlers.GetUsers
	updateUser       handlers.UpdateUser
	deleteUser       handlers.DeleteUser
	followUser       handlers.FollowUser
	unfollowUser     handlers.UnfollowUser
	getUserFollowers handlers.GetUserFollowers
	getFollowedUsers handlers.GetFollowedUsers
}

func (s *Server) Run() {
	err := s.router.Run(fmt.Sprintf("0.0.0.0:%v", os.Getenv("SERVICE_PORT")))
	if err != nil {
		panic(err)
	}
}

func NewServer() *Server {
	db, err := database.NewPostgresDBClient()
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Administrator{})
	if err != nil {
		panic(err)
	}

	logger, _ := zap.NewDevelopment()

	sdkJson, err := base64.StdEncoding.DecodeString(os.Getenv("FIREBASE_B64_SDK_JSON"))
	if err != nil {
		panic(err)
	}

	pubJwtKey, err := base64.StdEncoding.DecodeString(os.Getenv("PUB_RSA_B64"))
	if err != nil {
		panic(err)
	}

	privJwtKey, err := base64.StdEncoding.DecodeString(os.Getenv("PRIV_RSA_B64"))
	if err != nil {
		panic(err)
	}
	toker, err := utils.NewJwtToker(privJwtKey, pubJwtKey)
	if err != nil {
		panic(err)
	}

	// REPOSITORIES
	firebaseRepo, err := repositories.NewFirebaseRepository(logger, sdkJson, os.Getenv("FIREBASE_BUCKET_NAME"))
	if err != nil {
		panic(err)
	}
	userRepo := repositories.NewUserRepository(db, logger, firebaseRepo)
	adminRepo := repositories.NewAdminRepository(db, logger)

	// USECASES
	registerUc := accounts.NewRegisterImpl(userRepo, logger, firebaseRepo)
	adminRegisterUc := accounts.NewAdminRegistererImpl(adminRepo, logger, toker)
	getUserUc := users.NewUserGetterImpl(userRepo, firebaseRepo, logger)
	updateUserUc := users.NewUserUpdaterImpl(userRepo, firebaseRepo)
	deleteUserUc := users.NewUserDeleterImpl(userRepo)
	followUserUc := users.NewUserFollowerImpl(userRepo)

	// HANDLERS
	register := handlers.NewRegister(&registerUc, logger)
	finishRegister := handlers.NewFinishRegister(&registerUc, logger)
	adminRegister := handlers.NewAdminRegister(&adminRegisterUc, logger)
	adminLogin := handlers.NewAdminLogin(&adminRegisterUc, logger)

	getUserByID := handlers.NewGetUserByID(&getUserUc, logger)
	getUsers := handlers.NewGetUsers(&getUserUc, logger)
	updateUser := handlers.NewUpdateUser(&updateUserUc, logger)
	deleteUser := handlers.NewDeleteUser(&deleteUserUc, logger)

	followUser := handlers.NewFollowUser(&followUserUc, logger)
	unfollowUser := handlers.NewUnfollowUser(&followUserUc, logger)
	getUserFollowers := handlers.NewGetUserFollowers(&getUserUc, logger)
	getFollowedUsers := handlers.NewGetFollowedUsers(&getUserUc, logger)

	return &Server{
		router:           gin.Default(),
		register:         register,
		finishRegister:   finishRegister,
		adminRegister:    adminRegister,
		adminLogin:       adminLogin,
		getUserByID:      getUserByID,
		getUsers:         getUsers,
		updateUser:       updateUser,
		deleteUser:       deleteUser,
		followUser:       followUser,
		unfollowUser:     unfollowUser,
		getUserFollowers: getUserFollowers,
		getFollowedUsers: getFollowedUsers,
	}
}
