package server

import (
	"context"
	"os"
	"strconv"

	firebase "firebase.google.com/go/v4"
	"github.com/fiufit/users/database"
	"github.com/fiufit/users/handlers"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/repositories"
	"github.com/fiufit/users/usecases/accounts"
	"github.com/fiufit/users/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type Server struct {
	router *gin.Engine

	register       handlers.Register
	finishRegister handlers.FinishRegister
}

func (s *Server) Run() {
	err := s.router.Run(":8888")
	if err != nil {
		panic(err)
	}
}

func NewServer() *Server {
	db, err := database.NewPostgresDBClient()
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.User{})

	logger, _ := zap.NewDevelopment()

	opt := option.WithCredentialsFile("./firebase-adminsdk.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(err)
	}

	auth, err := app.Auth(context.Background())

	fromMail := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	mail := utils.NewMailerImpl(fromMail, password, host, port)

	// REPOSITORIES
	userRepo := repositories.NewUserRepository(db, logger)

	// USECASES
	registerUc := accounts.NewRegisterImpl(userRepo, logger, auth, mail)

	// HANDLERS
	register := handlers.NewRegister(&registerUc, logger)
	finishRegister := handlers.NewFinishRegister(&registerUc, logger)

	return &Server{
		router:         gin.Default(),
		register:       register,
		finishRegister: finishRegister,
	}
}
