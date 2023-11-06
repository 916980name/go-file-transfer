package errno

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	OK = &Errno{HTTP: 200, Code: "", Message: ""}

	InternalServerError = &Errno{HTTP: 500, Code: "InternalError", Message: "Internal server error."}

	ErrPageNotFound = &Errno{HTTP: 404, Code: "ResourceNotFound.PageNotFound", Message: "Page not found."}

	ErrBind = &Errno{HTTP: 400, Code: "InvalidParameter.BindError", Message: "Error occurred while binding the request body to the struct."}

	ErrInvalidParameter = &Errno{HTTP: 400, Code: "InvalidParameter", Message: "Parameter verification failed."}
)

type ErrResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteResponse(c context.Context, w http.ResponseWriter, err error, data interface{}) {
	if err != nil {
		hcode, code, message := Decode(err)
		marshalResponse(w, hcode, ErrResponse{
			Code:    code,
			Message: message,
		})

		return
	}
	marshalResponse(w, 200, data)
}

func marshalResponse(w http.ResponseWriter, hcode int, data interface{}) {
	result, eR := json.Marshal(data)
	if eR != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(hcode)
	w.Write(result)
}

type Errno struct {
	HTTP    int
	Code    string
	Message string
}

func (err *Errno) Error() string {
	return err.Message
}

func (err *Errno) SetMessage(format string, args ...interface{}) *Errno {
	err.Message = fmt.Sprintf(format, args...)
	return err
}

func Decode(err error) (int, string, string) {
	if err == nil {
		return OK.HTTP, OK.Code, OK.Message
	}

	switch typed := err.(type) {
	case *Errno:
		return typed.HTTP, typed.Code, typed.Message
	default:
	}

	return InternalServerError.HTTP, InternalServerError.Code, err.Error()
}
