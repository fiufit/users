package utils

import (
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

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
	params.SetBody("Gracias por registrarte en FiuFit! 💪\nTu código de verificación es *" + pin + "*\nIngrésalo en la app para comenzar a entrenar ahora mismo!")
	_, err := w.twilioClient.Api.CreateMessage(params)
	return err
}
