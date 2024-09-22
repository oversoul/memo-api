package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type response struct {
	Error     string            `json:"error,omitempty"`
	Status    string            `json:"status"`
	HTTPCode  int               `json:"http_code"`
	Datetime  string            `json:"datetime"`
	Timestamp int64             `json:"timestamp"`
	Errors    map[string]string `json:"errors,omitempty"`
}

func ErrMessage(w http.ResponseWriter, msg string, status int) {
	RespondErr(w, ErrorResponse{Message: msg, Status: status})
}

func ValidationErr(w http.ResponseWriter, errors map[string]string) {
	resp := response{
		Status:    "fail",
		Datetime:  time.Now().Format(time.DateTime),
		Timestamp: time.Now().Unix(),
		HTTPCode:  http.StatusUnprocessableEntity,
		Error:     "request not valid",
		Errors:    errors,
	}

	Respond(w, resp, resp.HTTPCode)
}

func RespondErr(w http.ResponseWriter, err error) {
	var errresp ErrorResponse

	if !errors.As(err, &errresp) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := response{
		Status:    "fail",
		Datetime:  time.Now().Format(time.DateTime),
		Timestamp: time.Now().Unix(),
		HTTPCode:  errresp.Status,
		Error:     errresp.Message,
	}

	Respond(w, resp, resp.HTTPCode)
}

func RespondSuccess(w http.ResponseWriter) {
	resp := response{
		Status:    "success",
		Datetime:  time.Now().Format(time.DateTime),
		Timestamp: time.Now().Unix(),
		HTTPCode:  http.StatusOK,
	}

	Respond(w, resp, resp.HTTPCode)
}

func Respond(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
