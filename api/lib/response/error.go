package response

import (
	"net/http"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"
)

type DefaultError struct {
	Message string `json:"message"`
}

func Error(w http.ResponseWriter, code int, msg string) {
	resp := DefaultError{Message: msg}
	Write(w, code, resp)
}

/*
RespondErrorIfNeeded writes error response if err != nil and map grpc status code to http status code
If err == nil do nothing and respond false
*/
func RespondErrorIfNeeded(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	code := status.Code(err)
	httpCode := http.StatusInternalServerError
	switch code {
	case codes.InvalidArgument:
		httpCode = http.StatusBadRequest
	case codes.NotFound:
		httpCode = http.StatusNotFound
	}
	Error(w, httpCode, err.Error())
	return true
}
