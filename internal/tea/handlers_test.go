package tea

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RegisterPayload struct {
	RequestBody  map[string]string `json:"requestBody"`
	Headers      map[string]string `json:"headers"`
	ResponseBody map[string]string `json:"responseBody"`
	URL          string            `json:"URL"`
	Method       string            `json:"Method"`
}

func TestRegisterHandler(t *testing.T) {
	payload := RegisterPayload{
		RequestBody:  map[string]string{"a": "b"},
		Headers:      map[string]string{"h1": "v1"},
		ResponseBody: map[string]string{"b": "c"},
		URL:          "/test",
		Method:       "POST",
	}
	var reqBuffer bytes.Buffer
	err := json.NewEncoder(&reqBuffer).Encode(payload)
	if err != nil {
		t.Error("failed to encode request")
	}
	req := httptest.NewRequest("POST", "/register-request", &reqBuffer)

	rs := &RequestsStore{}
	registerHandler := &RegisterHandler{RequestsStore: rs}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(registerHandler.Handler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

}
