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
	"github.com/fiufit/users/usecases/certifications"
	"github.com/fiufit/users/usecases/users"
	"github.com/fiufit/users/utils"
	"github.com/gin-gonic/gin"
	twilio "github.com/twilio/twilio-go"
	"go.uber.org/zap"
)

type Server struct {
	router                *gin.Engine
	register              handlers.Register
	finishRegister        handlers.FinishRegister
	adminRegister         handlers.AdminRegister
	adminLogin            handlers.AdminLogin
	getUserByID           handlers.GetUserByID
	getUsers              handlers.GetUsers
	updateUser            handlers.UpdateUser
	deleteUser            handlers.DeleteUser
	followUser            handlers.FollowUser
	unfollowUser          handlers.UnfollowUser
	getUserFollowers      handlers.GetUserFollowers
	getFollowedUsers      handlers.GetFollowedUsers
	enableUser            handlers.EnableUser
	disableUser           handlers.DisableUser
	notifyUserLogin       handlers.NotifyUserLogin
	notifyPasswordRecover handlers.NotifyPasswordRecover
	getClosestUsers       handlers.GetClosestUsers
	sendVerificationPin   handlers.SendVerificationPin
	verifyUser            handlers.VerifyUser
	createCert            handlers.CreateCertification
	updateCert            handlers.UpdateCertification
	getCert               handlers.GetCertifications
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

	err = db.AutoMigrate(
		&models.User{},
		&models.Administrator{},
		&models.Interest{},
		&models.VerificationPin{},
		&models.Certification{},
	)
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

	reverseLocator, _ := utils.NewReverseLocator()

	whatsAppSender := utils.NewWhatsApperImpl(os.Getenv("TWILIO_PHONE_NUMBER"), twilio.NewRestClient())

	metricsUrl := os.Getenv("METRICS_SERVICE_URL")
	notificationUrl := os.Getenv("NOTIFICATION_SERVICE_URL")

	// REPOSITORIES
	firebaseRepo, err := repositories.NewFirebaseRepository(logger, sdkJson, os.Getenv("FIREBASE_BUCKET_NAME"))
	if err != nil {
		panic(err)
	}
	userRepo := repositories.NewUserRepository(db, logger, firebaseRepo, reverseLocator)
	adminRepo := repositories.NewAdminRepository(db, logger)
	metricsRepo := repositories.NewMetricsRepository(metricsUrl, "v1", logger)
	notificationRepo := repositories.NewNotificationRepository(notificationUrl, logger, "v1")
	verificationRepo := repositories.NewVerificationPinRepository(db, logger)
	certificationRepo := repositories.NewCertificationRepository(db, logger, firebaseRepo)

	// USECASES
	registerUc := accounts.NewRegisterImpl(userRepo, logger, firebaseRepo, metricsRepo)
	adminRegisterUc := accounts.NewAdminRegistererImpl(adminRepo, logger, toker)
	getUserUc := users.NewUserGetterImpl(userRepo, logger)
	updateUserUc := users.NewUserUpdaterImpl(userRepo, metricsRepo)
	deleteUserUc := users.NewUserDeleterImpl(userRepo)
	followUserUc := users.NewUserFollowerImpl(userRepo, notificationRepo, metricsRepo, logger)
	enableUserUc := users.NewUserEnablerImpl(userRepo, firebaseRepo, metricsRepo, logger)
	verificationUc := accounts.NewVerifierImpl(verificationRepo, firebaseRepo, whatsAppSender, logger)
	createCertUc := certifications.NewCertificationCreator(certificationRepo, userRepo)
	updateCertUc := certifications.NewCertificationUpdaterImpl(certificationRepo, userRepo, notificationRepo, firebaseRepo, logger)
	getCertUc := certifications.NewCertificationGetterImpl(certificationRepo, userRepo)

	// HANDLERS
	register := handlers.NewRegister(&registerUc, logger)
	finishRegister := handlers.NewFinishRegister(&registerUc, logger)
	adminRegister := handlers.NewAdminRegister(&adminRegisterUc, logger)
	adminLogin := handlers.NewAdminLogin(&adminRegisterUc, logger)
	sendVerificationPin := handlers.NewSendVerificationPin(&verificationUc, logger)
	verifyUser := handlers.NewVerifyUser(&verificationUc, logger)

	getUserByID := handlers.NewGetUserByID(&getUserUc, logger)
	getUsers := handlers.NewGetUsers(&getUserUc, logger)
	getClosestUsers := handlers.NewGetClosestUsers(&getUserUc, logger)
	updateUser := handlers.NewUpdateUser(&updateUserUc, logger)
	deleteUser := handlers.NewDeleteUser(&deleteUserUc, logger)

	createCertification := handlers.NewCreateCertification(createCertUc)
	updateCertification := handlers.NewUpdateCertification(updateCertUc)
	getCertifications := handlers.NewGetCertifications(getCertUc)

	followUser := handlers.NewFollowUser(&followUserUc, logger)
	unfollowUser := handlers.NewUnfollowUser(&followUserUc, logger)
	getUserFollowers := handlers.NewGetUserFollowers(&getUserUc, logger)
	getFollowedUsers := handlers.NewGetFollowedUsers(&getUserUc, logger)
	enableUser := handlers.NewEnableUser(&enableUserUc, logger)
	disableUser := handlers.NewDisableUser(&enableUserUc, logger)
	notifyPasswordRecover := handlers.NewNotifyPasswordRecover(metricsRepo)
	notifyUserLogin := handlers.NewNotifyUserLogin(metricsRepo)

	return &Server{
		router:                gin.Default(),
		register:              register,
		finishRegister:        finishRegister,
		adminRegister:         adminRegister,
		adminLogin:            adminLogin,
		getUserByID:           getUserByID,
		getUsers:              getUsers,
		updateUser:            updateUser,
		deleteUser:            deleteUser,
		followUser:            followUser,
		unfollowUser:          unfollowUser,
		getUserFollowers:      getUserFollowers,
		getFollowedUsers:      getFollowedUsers,
		enableUser:            enableUser,
		disableUser:           disableUser,
		getClosestUsers:       getClosestUsers,
		notifyUserLogin:       notifyUserLogin,
		notifyPasswordRecover: notifyPasswordRecover,
		sendVerificationPin:   sendVerificationPin,
		verifyUser:            verifyUser,
		createCert:            createCertification,
		updateCert:            updateCertification,
		getCert:               getCertifications,
	}
}
