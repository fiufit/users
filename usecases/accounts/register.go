package accounts

import (
	"context"

	"firebase.google.com/go/v4/auth"
	"github.com/fiufit/users/utils"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Register interface {
	Execute(ctx context.Context, email, password, nickname string) error
}

type RegisterImpl struct {
	db     *sqlx.DB
	logger *zap.Logger
	auth   *auth.Client
	mailer utils.Mailer
}

func NewRegisterImpl(db *sqlx.DB, logger *zap.Logger, auth *auth.Client, mailer utils.Mailer) *RegisterImpl {
	return &RegisterImpl{db: db, logger: logger, auth: auth, mailer: mailer}
}

func (uc RegisterImpl) Execute(ctx context.Context, email, password, nickname string) error {
	params := (&auth.UserToCreate{}).Email(email).Password(password).DisplayName(nickname).EmailVerified(false)
	user, err := uc.auth.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	verificationLink, err := uc.auth.EmailVerificationLink(ctx, user.Email)
	if err != nil {
		uc.logger.Error("Unable to generate verification link for email", zap.String("email", email), zap.Error(err))
		return err
	}

	/*TODO Find out how to user firebase's email verification instead of our own mail account + SMTP server. Apparently
	the Go firebase SDK doesn't have auth.SendEmailVerification()
	*/
	return uc.mailer.SendAccountVerificationEmail(user.Email, verificationLink)
}
