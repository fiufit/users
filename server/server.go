package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/fiufit/users/database"
	"github.com/fiufit/users/handlers"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories"
	"github.com/fiufit/users/usecases/accounts"
	"github.com/fiufit/users/usecases/users"
	"github.com/fiufit/users/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type Server struct {
	router *gin.Engine

	register          handlers.Register
	finishRegister    handlers.FinishRegister
	adminRegister     handlers.AdminRegister
	adminLogin        handlers.AdminLogin
	getUserByID       handlers.GetUserByID
	getUserByNickname handlers.GetUsers
	updateUser        handlers.UpdateUser
	deleteUser        handlers.DeleteUser
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

	opt := option.WithCredentialsJSON(sdkJson)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(err)
	}

	auth, err := app.Auth(context.Background())
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
	userRepo := repositories.NewUserRepository(db, logger, auth)
	adminRepo := repositories.NewAdminRepository(db, logger)

	// USECASES
	registerUc := accounts.NewRegisterImpl(userRepo, logger, auth)
	adminRegisterUc := accounts.NewAdminRegistererImpl(adminRepo, logger, toker)
	getUserUc := users.NewUserGetterImpl(userRepo, logger)
	updateUserUc := users.NewUserUpdaterImpl(userRepo)
	deleteUserUc := users.NewUserDeleterImpl(userRepo)

	// HANDLERS
	register := handlers.NewRegister(&registerUc, logger)
	finishRegister := handlers.NewFinishRegister(&registerUc, logger)
	adminRegister := handlers.NewAdminRegister(&adminRegisterUc, logger)
	adminLogin := handlers.NewAdminLogin(&adminRegisterUc, logger)

	getUserByID := handlers.NewGetUserByID(&getUserUc, logger)
	getUserByNickname := handlers.NewGetUsers(&getUserUc, logger)
	updateUser := handlers.NewUpdateUser(&updateUserUc, logger)
	deleteUser := handlers.NewDeleteUser(&deleteUserUc, logger)

	return &Server{
		router:            gin.Default(),
		register:          register,
		finishRegister:    finishRegister,
		adminRegister:     adminRegister,
		adminLogin:        adminLogin,
		getUserByID:       getUserByID,
		getUserByNickname: getUserByNickname,
		updateUser:        updateUser,
		deleteUser:        deleteUser,
	}
}
