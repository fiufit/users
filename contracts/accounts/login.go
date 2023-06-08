package accounts

import "errors"

type NotifyLoginMethodRequest struct {
	Method string `form:"method"`
}

func ValidateMethod(method string) error {
	if method != "mail" && method != "federated_entity" {
		return errors.New("invalid method")
	}
	return nil
}
