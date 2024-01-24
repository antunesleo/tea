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
		"expectedRequest": map[string]any{
			"body": map[string]string{
				"a": "b",
			},
			"url":     "/randomjoke/1",
			"method":  http.MethodPost,
			"headers": map[string]string{"h1": "v1"},
		},
		"wantedResponse": map[string]any{
			"body": map[string]string{"b": "c"},
		},
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
	t.Run("missing expected request headers", func(t *testing.T) {
		payload := map[string]any{
			"expectedRequest": map[string]any{
				"body": map[string]string{
					"a": "b",
				},
				"url":    "/randomjoke/1",
				"method": http.MethodPost,
			},
			"wantedResponse": map[string]any{
				"body": map[string]string{"b": "c"},
			},
		}
		reqBuffer := reqBufferFromPayload(payload, t)
		req := httptest.NewRequest(http.MethodPost, "/register-request", &reqBuffer)

		registerHandler := NewRegisterHandler(NewRequestsStore())
		responseRecorder := httptest.NewRecorder()
		handler := http.HandlerFunc(registerHandler.Handler)
		handler.ServeHTTP(responseRecorder, req)

		assertStatusCode(responseRecorder.Code, http.StatusBadRequest, t)
	})

	t.Run("missing wanted response body", func(t *testing.T) {
		payload := map[string]any{
			"expectedRequest": map[string]any{
				"body": map[string]string{
					"a": "b",
				},
				"url":     "/randomjoke/1",
				"method":  http.MethodPost,
				"headers": map[string]string{"h1": "v1"},
			},
			"wantedResponse": map[string]any{},
		}

		reqBuffer := reqBufferFromPayload(payload, t)
		req := httptest.NewRequest(http.MethodPost, "/register-request", &reqBuffer)

		registerHandler := NewRegisterHandler(NewRequestsStore())
		responseRecorder := httptest.NewRecorder()
		handler := http.HandlerFunc(registerHandler.Handler)
		handler.ServeHTTP(responseRecorder, req)

		assertStatusCode(responseRecorder.Code, http.StatusBadRequest, t)
	})

	t.Run("missing expected request url", func(t *testing.T) {
		payload := map[string]any{
			"expectedRequest": map[string]any{
				"body": map[string]string{
					"a": "b",
				},
				"method":  http.MethodPost,
				"headers": map[string]string{"h1": "v1"},
			},
			"wantedResponse": map[string]any{
				"body": map[string]string{"b": "c"},
			},
		}

		reqBuffer := reqBufferFromPayload(payload, t)
		req := httptest.NewRequest(http.MethodPost, "/register-request", &reqBuffer)

		registerHandler := NewRegisterHandler(NewRequestsStore())
		responseRecorder := httptest.NewRecorder()
		handler := http.HandlerFunc(registerHandler.Handler)
		handler.ServeHTTP(responseRecorder, req)

		assertStatusCode(responseRecorder.Code, http.StatusBadRequest, t)
	})

	t.Run("missing expected request method", func(t *testing.T) {
		payload := map[string]any{
			"expectedRequest": map[string]any{
				"body": map[string]string{
					"a": "b",
				},
				"url":     "/randomjoke/1",
				"headers": map[string]string{"h1": "v1"},
			},
			"wantedResponse": map[string]any{
				"body": map[string]string{"b": "c"},
			},
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
