package tea

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	payload := map[string]any{
		"requestBody":  map[string]string{"a": "b"},
		"headers":      map[string]string{"h1": "v1"},
		"responseBody": map[string]string{"b": "c"},
		"url":          "/test",
		"method":       http.MethodPost,
	}
	reqBuffer := reqBufferFromPayload(payload, t)
	req := httptest.NewRequest(http.MethodPost, "/register-request", &reqBuffer)
	registerHandler := NewRegisterHandler(NewRequestsStore())
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(registerHandler.Handler)
	handler.ServeHTTP(responseRecorder, req)
	assertStatusCode(responseRecorder.Code, http.StatusCreated, t)
}

func TestRegisterHandlerBadRequest(t *testing.T) {
	t.Run("missing request body", func(t *testing.T) {
		payload := map[string]any{
			"headers":      map[string]string{"h1": "v1"},
			"responseBody": map[string]string{"b": "c"},
			"url":          "/test",
			"method":       http.MethodPost,
		}

		reqBuffer := reqBufferFromPayload(payload, t)
		req := httptest.NewRequest(http.MethodPost, "/register-request", &reqBuffer)

		registerHandler := NewRegisterHandler(NewRequestsStore())
		responseRecorder := httptest.NewRecorder()
		handler := http.HandlerFunc(registerHandler.Handler)
		handler.ServeHTTP(responseRecorder, req)

		assertStatusCode(responseRecorder.Code, http.StatusBadRequest, t)
	})

	t.Run("missing headers", func(t *testing.T) {
		payload := map[string]any{
			"requestBody":  map[string]string{"a": "b"},
			"responseBody": map[string]string{"b": "c"},
			"url":          "/test",
			"method":       http.MethodPost,
		}

		reqBuffer := reqBufferFromPayload(payload, t)
		req := httptest.NewRequest(http.MethodPost, "/register-request", &reqBuffer)

		registerHandler := NewRegisterHandler(NewRequestsStore())
		responseRecorder := httptest.NewRecorder()
		handler := http.HandlerFunc(registerHandler.Handler)
		handler.ServeHTTP(responseRecorder, req)

		assertStatusCode(responseRecorder.Code, http.StatusBadRequest, t)
	})

	t.Run("missing responseBody", func(t *testing.T) {
		payload := map[string]any{
			"requestBody": map[string]string{"a": "b"},
			"headers":     map[string]string{"h1": "v1"},
			"url":         "/test",
			"method":      http.MethodPost,
		}

		reqBuffer := reqBufferFromPayload(payload, t)
		req := httptest.NewRequest(http.MethodPost, "/register-request", &reqBuffer)

		registerHandler := NewRegisterHandler(NewRequestsStore())
		responseRecorder := httptest.NewRecorder()
		handler := http.HandlerFunc(registerHandler.Handler)
		handler.ServeHTTP(responseRecorder, req)

		assertStatusCode(responseRecorder.Code, http.StatusBadRequest, t)
	})

	t.Run("missing url", func(t *testing.T) {
		payload := map[string]any{
			"requestBody":  map[string]string{"a": "b"},
			"headers":      map[string]string{"h1": "v1"},
			"responseBody": map[string]string{"b": "c"},
			"method":       http.MethodPost,
		}

		reqBuffer := reqBufferFromPayload(payload, t)
		req := httptest.NewRequest(http.MethodPost, "/register-request", &reqBuffer)

		registerHandler := NewRegisterHandler(NewRequestsStore())
		responseRecorder := httptest.NewRecorder()
		handler := http.HandlerFunc(registerHandler.Handler)
		handler.ServeHTTP(responseRecorder, req)

		assertStatusCode(responseRecorder.Code, http.StatusBadRequest, t)
	})

	t.Run("missing method", func(t *testing.T) {
		payload := map[string]any{
			"requestBody":  map[string]string{"a": "b"},
			"headers":      map[string]string{"h1": "v1"},
			"responseBody": map[string]string{"b": "c"},
			"url":          "/test",
		}

		reqBuffer := reqBufferFromPayload(payload, t)
		req := httptest.NewRequest(http.MethodPost, "/register-request", &reqBuffer)

		registerHandler := NewRegisterHandler(NewRequestsStore())
		responseRecorder := httptest.NewRecorder()
		handler := http.HandlerFunc(registerHandler.Handler)
		handler.ServeHTTP(responseRecorder, req)

		assertStatusCode(responseRecorder.Code, http.StatusBadRequest, t)
	})

}

func assertStatusCode(got, want int, t *testing.T) {
	if got != want {
		t.Errorf("handler returned wrong status code: got %v want %v", got, want)
	}
}

func reqBufferFromPayload(payload map[string]any, t *testing.T) bytes.Buffer {
	var reqBuffer bytes.Buffer
	err := json.NewEncoder(&reqBuffer).Encode(payload)
	if err != nil {
		t.Error("failed to encode request")
	}
	return reqBuffer
}
