package main

import (
	"context"
	"os"
	"strconv"

	firebase "firebase.google.com/go/v4"
	"github.com/fiufit/users/database"
	"github.com/fiufit/users/usecases/accounts"
	"github.com/fiufit/users/utils"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

func main() {
	db, err := database.NewPostgresDBClient()
	if err != nil {
		panic(err)
	}

	defer db.Close()

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

	registerUc := accounts.NewRegisterImpl(db, logger, auth, mail)
	err = registerUc.Execute(context.Background(), "ejemplo@fi.uba.ar", "123456", "ejemplo")
	if err != nil {
		panic(err)
	}
}
