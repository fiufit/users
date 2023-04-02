package contracts

type OkResponse struct {
	Data interface{} `json:"data"`
}

type ErrResponse struct {
	Error string `json:"error"`
}

func FormatOkResponse(data interface{}) OkResponse {
	return OkResponse{data}
}

func FormatErrResponse(err error) ErrResponse {
	return ErrResponse{err.Error()}
}
