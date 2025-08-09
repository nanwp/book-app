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

func failResponseWriter(w http.ResponseWriter, r *http.Request, err error, errStatusCode int) {
	w.Header().Set("Content-Type", "application/json")

	log.Error().Err(err).Str("request", r.URL.Path).Msg("error processing request")

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

func successResponseWriter(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	log.Info().Str("request", r.URL.Path).Msg("success processing request")

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

func WriteResponse(w http.ResponseWriter, r *http.Request, err error, data any) {
	switch err.(type) {
	case *ErrForbidden, ErrForbidden:
		failResponseWriter(w, r, err, http.StatusForbidden)
	case *ErrUnauthorized, ErrUnauthorized:
		failResponseWriter(w, r, err, http.StatusUnauthorized)
	case *ErrNotFound, ErrNotFound:
		failResponseWriter(w, r, err, http.StatusNotFound)
	case *ErrBadRequest, ErrBadRequest:
		failResponseWriter(w, r, err, http.StatusBadRequest)
	case *ErrInternalServer, ErrInternalServer:
		failResponseWriter(w, r, err, http.StatusInternalServerError)
	case nil:
		successResponseWriter(w, r, data, http.StatusOK)
	default:
		failResponseWriter(w, r, err, http.StatusInternalServerError)
	}
}
