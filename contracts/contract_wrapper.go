package contracts

import "encoding/json"

func UnwrapOkResponse(data []byte, contract interface{}) (interface{}, error) {
	okRes := OkResponse{contract}
	err := json.Unmarshal(data, &okRes)
	if err != nil {
		return nil, err
	}
	contract = okRes.Data
	return contract, nil
}

func UnwrapError(data []byte) error {
	var errRes ErrResponse
	err := json.Unmarshal(data, &errRes)
	if err != nil {
		return err
	}
	code := errRes.Err.Code
	err, ok := externalCodes[code]
	if !ok {
		return ErrInternal
	}
	return err
}
