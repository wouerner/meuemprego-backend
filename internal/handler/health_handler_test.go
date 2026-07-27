package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wouerner/runter-backend/internal/handler"
)

func TestHealthHandler_Check(t *testing.T) {
	// Arrange
	healthHandler := handler.NewHealthHandler()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/health", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	// Act
	healthHandler.Check(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var body map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &body)
	assert.NoError(t, err)

	assert.Equal(t, "ok", body["status"])
	assert.Equal(t, "v1", body["version"])
	assert.Equal(t, "runter-backend-api", body["service"])
	assert.NotEmpty(t, body["timestamp"])
}
