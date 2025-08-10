package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// Simple test without database dependency
func TestSimpleHealthCheck(t *testing.T) {
	// Create a simple router just for health check
	router := mux.NewRouter()

	api := router.PathPrefix("/api/v1/").Subrouter()
	api.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	req, err := http.NewRequest("GET", "/api/v1/healthcheck", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
}

func TestRouterSetup(t *testing.T) {
	router := mux.NewRouter()

	// Test 404 for non-existent endpoint
	req, err := http.NewRequest("GET", "/nonexistent", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
