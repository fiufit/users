package contracts

var errCodes = map[error]string{
	ErrInternal:          "U0",
	ErrBadRequest:        "U1",
	ErrUserNotFound:      "U2",
	ErrUserAlreadyExists: "U3",
}

type OkResponse struct {
	Data interface{} `json:"data"`
}

type ErrResponse struct {
	Err ErrPayload `json:"error"`
}

type ErrPayload struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func FormatOkResponse(data interface{}) OkResponse {
	return OkResponse{data}
}

func FormatErrResponse(err error) ErrResponse {
	errCode, ok := errCodes[err]
	if !ok {
		errCode = "U0"
	}

	payload := ErrPayload{
		Description: err.Error(),
		Code:        errCode,
	}

	return ErrResponse{payload}
}
