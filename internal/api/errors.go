package api

import (
	"encoding/json"
	"github.com/broswen/vex/internal/account"
	"github.com/broswen/vex/internal/flag"
	"github.com/broswen/vex/internal/project"
	"net/http"
)

var (
	ErrUnknown        = NewAPIError(http.StatusInternalServerError, 9999, "unknown error")
	ErrInternalServer = NewAPIError(http.StatusInternalServerError, 9500, "internal error")
	ErrBadRequest     = NewAPIError(http.StatusBadRequest, 9400, "bad request")
	ErrNotFound       = NewAPIError(http.StatusNotFound, 9404, "not found")
	ErrUnauthorized   = NewAPIError(http.StatusUnauthorized, 9401, "unauthorized")
)

type APIError struct {
	Status  int    `json:"-"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e APIError) Error() string {
	return e.Message
}

func NewAPIError(status, code int, message string) *APIError {
	return &APIError{
		Status:  status,
		Code:    code,
		Message: message,
		Err:     nil,
	}
}

func translateError(err error) *APIError {
	switch err.(type) {
	case account.ErrAccountNotFound,
		project.ErrProjectNotFound,
		flag.ErrFlagNotFound:
		return ErrNotFound
	case account.ErrInvalidData,
		project.ErrInvalidData,
		flag.ErrInvalidData:
		return ErrBadRequest.WithError(err)
	case flag.ErrKeyNotUnique:
		return ErrBadRequest.WithError(err)
	default:
		return ErrUnknown
	}
}

func (e *APIError) WithError(err error) *APIError {
	temp := &APIError{
		Status: e.Status,
		Code:   e.Code,
		Err:    err,
	}
	if e.Message == "" {
		temp.Message = err.Error()
	} else {
		temp.Message = e.Message + ": " + err.Error()
	}
	return temp
}

func renderError(rw http.ResponseWriter, apiError *APIError) error {
	j, err := json.Marshal(apiError)
	if err != nil {
		return err
	}
	rw.WriteHeader(apiError.Status)
	_, err = rw.Write(j)
	if err != nil {
		return err
	}
	return nil
}
