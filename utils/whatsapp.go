package utils

import (
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

//go:generate mockery --name WhatsApper
type WhatsApper interface {
	SendWhatsAppMessage(to string, pin string) error
}

type WhatsApperImpl struct {
	fromPhoneNumber string
	twilioClient    *twilio.RestClient
}

func NewWhatsApperImpl(fromPhoneNumber string, twilioClient *twilio.RestClient) WhatsApperImpl {
	return WhatsApperImpl{fromPhoneNumber, twilioClient}
}

func (w WhatsApperImpl) SendWhatsAppMessage(to string, pin string) error {
	params := &openapi.CreateMessageParams{}
	params.SetTo("whatsapp:" + to)
	params.SetFrom("whatsapp:" + w.fromPhoneNumber)
	params.SetBody("Thanks for signing up for FiuFit! 💪\nYour verification pin is *" + pin + "*\nEnter this code in our app to activate your account and start training!")
	_, err := w.twilioClient.Api.CreateMessage(params)
	return err
}
