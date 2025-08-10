package helper

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  string      `json:"errors,omitempty"`
}

func failResponseWriter(w http.ResponseWriter, err error, errStatusCode int) {
	w.Header().Set("Content-Type", "application/json")

	var resp Response
	w.WriteHeader(errStatusCode)
	resp.Code = errStatusCode
	resp.Message = err.Error()
	resp.Data = nil

	responseBytes, _ := json.Marshal(resp)
	if _, writeErr := w.Write(responseBytes); writeErr != nil {
		log.Error().Err(writeErr).Msg("failed to write response")
	}
}

func successResponseWriter(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	var resp Response
	w.WriteHeader(statusCode)
	resp.Code = statusCode
	resp.Message = "success"
	resp.Data = data

	responseBytes, _ := json.Marshal(resp)
	if _, writeErr := w.Write(responseBytes); writeErr != nil {
		log.Error().Err(writeErr).Msg("failed to write response")
	}
}

func WriteResponse(w http.ResponseWriter, err error, data any) {
	switch err.(type) {
	case *ErrForbidden, ErrForbidden:
		failResponseWriter(w, err, http.StatusForbidden)
	case *ErrUnauthorized, ErrUnauthorized:
		failResponseWriter(w, err, http.StatusUnauthorized)
	case *ErrNotFound, ErrNotFound:
		failResponseWriter(w, err, http.StatusNotFound)
	case *ErrBadRequest, ErrBadRequest:
		failResponseWriter(w, err, http.StatusBadRequest)
	case *ErrInternalServer, ErrInternalServer:
		failResponseWriter(w, err, http.StatusInternalServerError)
	case nil:
		successResponseWriter(w, data, http.StatusOK)
	default:
		failResponseWriter(w, err, http.StatusInternalServerError)
	}
}
