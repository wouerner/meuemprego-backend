package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/wouerner/runter-backend/internal/middleware"
)

const testJWTSecret = "test_secret_key_123"

func TestJWTMiddleware_MissingHeader(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handlerToTest := middleware.JWTMiddleware(testJWTSecret)(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()

	handlerToTest.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestJWTMiddleware_InvalidFormat(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handlerToTest := middleware.JWTMiddleware(testJWTSecret)(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "InvalidTokenWithoutBearer")
	rec := httptest.NewRecorder()

	handlerToTest.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestJWTMiddleware_ValidToken(t *testing.T) {
	// Gerar token válido
	claims := jwt.MapClaims{
		"user_id": float64(10),
		"email":   "teste@exemplo.com",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(testJWTSecret))
	assert.NoError(t, err)

	nextCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
		assert.True(t, ok)
		assert.Equal(t, uint(10), userID)

		userEmail, ok := r.Context().Value(middleware.UserEmailKey).(string)
		assert.True(t, ok)
		assert.Equal(t, "teste@exemplo.com", userEmail)

		w.WriteHeader(http.StatusOK)
	})

	handlerToTest := middleware.JWTMiddleware(testJWTSecret)(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rec := httptest.NewRecorder()

	handlerToTest.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, nextCalled)
}
