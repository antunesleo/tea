package tea

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/antunesleo/tea/internal/map_utils"
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
			"body":       map[string]string{"b": "c"},
			"statusCode": http.StatusOK,
			"headers": map[string]string{
				"header1": "value1",
				"header2": "value2",
			},
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
			"wantedResponse": map[string]any{
				"statusCode": http.StatusOK,
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
				"body": map[string]any{
					"b":          "c",
					"statusCode": http.StatusOK,
				},
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
				"body": map[string]any{
					"b":          "c",
					"statusCode": http.StatusOK,
				},
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

	t.Run("missing wanted response status code", func(t *testing.T) {
		payload := map[string]any{
			"expectedRequest": map[string]any{
				"body": map[string]string{
					"a": "b",
				},
				"url":     "/randomjoke/1",
				"headers": map[string]string{"h1": "v1"},
				"method":  http.MethodPost,
			},
			"wantedResponse": map[string]any{
				"body": map[string]any{
					"b": "c",
				},
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
	t.Helper()
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

func HeaderToMap(header http.Header) map[string]string {
	m := map[string]string{}
	for key, values := range header {
		m[key] = values[0]
	}
	return m
}

func TestApiUnderTestHandler(t *testing.T) {
	t.Run("test succeed", func(t *testing.T) {
		wantedResponseHeader := map[string]string{
			"header1": "value1",
			"header2": "value2",
		}
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
				"statusCode": http.StatusAccepted,
				"body": map[string]string{
					"c": "d",
				},
				"headers": wantedResponseHeader,
			},
		}

		requestStore := NewRequestsStore()
		reqBuffer := reqBufferFromPayload(payload, t)
		req := httptest.NewRequest(http.MethodPost, "/register-request", &reqBuffer)
		registerHandler := NewRegisterHandler(requestStore)
		responseRecorder := httptest.NewRecorder()
		handler := http.HandlerFunc(registerHandler.Handler)
		handler.ServeHTTP(responseRecorder, req)
		assertStatusCode(responseRecorder.Code, http.StatusCreated, t)

		payload2 := map[string]any{
			"a": "b",
		}
		reqBuffer = reqBufferFromPayload(payload2, t)
		req = httptest.NewRequest(http.MethodPost, "/randomjoke/1", &reqBuffer)
		req.Header.Add("h1", "v1")
		apiUnderTestHandler := NewApiUnderTestHandler(requestStore)
		responseRecorder = httptest.NewRecorder()
		handler = http.HandlerFunc(apiUnderTestHandler.Handler)
		handler.ServeHTTP(responseRecorder, req)
		assertStatusCode(responseRecorder.Code, http.StatusAccepted, t)
		gotResponseHeaders := map_utils.LowercaseMap(HeaderToMap(responseRecorder.Header()))

		for wantedKey, wantedValue := range wantedResponseHeader {
			gotValue, ok := gotResponseHeaders[wantedKey]
			if !ok {
				t.Errorf("want response to have header %v", wantedKey)
			}

			if gotValue != wantedValue {
				t.Errorf("want value %v but got %v for header %v", wantedValue, gotValue, wantedKey)
			}
		}
	})
}
