package smartthings

import (
	"encoding/json"
	"fmt"
	"io"
)

type ErrorResponse struct {
	RequestID string `json:"requestId"`
	Error     *Error `json:"error"`
}

type Error struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Target  string   `json:"target"`
	Details []*Error `json:"details"`
}

func (e *Error) Error() string {
	return fmt.Sprint(e.Code, e.Message, e.Target)
}

func checkErrorResponse(r io.ReadCloser) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	var errResponse *ErrorResponse
	if err := json.Unmarshal(data, &errResponse); err == nil {
		if errResponse.Error != nil {
			return errResponse.Error
		}
	}

	return nil
}
