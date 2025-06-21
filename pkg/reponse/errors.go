package response

import "net/http"

type Err struct {
	StatusCode int
	Message    string `json:"message"`
}

func NewInternalErr(message string) *Err {

	return &Err{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
	}
}

func NewBadREquest(message string) *Err {

	return &Err{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}
