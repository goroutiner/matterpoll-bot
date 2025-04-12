package handlers_test

import (
	"matterpoll-bot/config"
	"matterpoll-bot/internal/handlers"
	"matterpoll-bot/internal/storage/store_mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestTokenValidatorMiddleware проверяет работу middleware для проверки токена.
func TestTokenValidatorMiddleware(t *testing.T) {
	mockStore := store_mocks.NewStoreInterface(t)

	// Создание обработчика для тестирования
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	config.Mode = "database"
	cmdPath := "test_path"
	handler := handlers.TokenValidatorMiddleware(mockStore, nextHandler)

	t.Run("valid token", func(t *testing.T) {
		token := "valid_token"

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		respRec := httptest.NewRecorder()

		req.ParseForm()
		req.Form.Add("command", cmdPath)
		req.Form.Add("token", token)

		mockStore.ExpectedCalls = nil
		mockStore.On("ValidateCmdToken", mock.Anything, mock.Anything).Return(true)
		handler.ServeHTTP(respRec, req)

		expectedResponse := "OK"
		actualResponse := respRec.Body.String()

		require.Equal(t, http.StatusOK, respRec.Code)
		require.Contains(t, actualResponse, expectedResponse)
		mockStore.AssertCalled(t, "ValidateCmdToken", cmdPath, token)
	})

	t.Run("invalid token", func(t *testing.T) {
		token := "invalids_token"

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		respRec := httptest.NewRecorder()

		req.ParseForm()
		req.Form.Add("command", cmdPath)
		req.Form.Add("token", token)

		mockStore.ExpectedCalls = nil
		mockStore.On("ValidateCmdToken", mock.Anything, mock.Anything).Return(false)
		handler.ServeHTTP(respRec, req)

		expectedResponse := "Invalid token"
		actualResponse := respRec.Body.String()

		require.Equal(t, http.StatusUnauthorized, respRec.Code)
		require.Contains(t, actualResponse, expectedResponse)
		mockStore.AssertCalled(t, "ValidateCmdToken", cmdPath, token)

        // empty token
		token = ""

		req = httptest.NewRequest(http.MethodGet, "/", nil)
		respRec = httptest.NewRecorder()

		req.ParseForm()
		req.Form.Add("command", cmdPath)
		req.Form.Add("token", token)

		mockStore.ExpectedCalls = nil
		handler.ServeHTTP(respRec, req)

		expectedResponse = "'command' or 'token' are empty in the form data"
		actualResponse = respRec.Body.String()

		require.Equal(t, http.StatusBadRequest, respRec.Code)
		require.Contains(t, actualResponse, expectedResponse)
	})
}
