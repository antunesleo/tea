package tea

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/antunesleo/tea/internal/map_utils"
)

type RegisterHandlerPayload struct {
	RequestBody  json.RawMessage   `json:"requestBody"`
	Headers      map[string]string `json:"headers"`
	ResponseBody json.RawMessage   `json:"responseBody"`
	URL          string            `json:url`
	Method       string            `json:Method`
}

type RegisterHandler struct {
	RequestsStore *RequestsStore
}

func (rh *RegisterHandler) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var registerPayload RegisterHandlerPayload
	err = json.Unmarshal(body, &registerPayload)
	if err != nil {
		errorMsg := fmt.Sprintf("Error parsing JSON body: %v", err)
		http.Error(w, errorMsg, http.StatusBadRequest)
		return
	}

	missingFields := []string{}
	if registerPayload.RequestBody == nil {
		missingFields = append(missingFields, "requestBody")
	}
	if registerPayload.Headers == nil {
		missingFields = append(missingFields, "headers")
	}
	if registerPayload.ResponseBody == nil {
		missingFields = append(missingFields, "responseBody")
	}
	if registerPayload.URL == "" {
		missingFields = append(missingFields, "url")
	}
	if registerPayload.Method == "" {
		missingFields = append(missingFields, "method")
	}

	if len(missingFields) > 0 {
		errorMessage := fmt.Sprintf("Missing mandatory fields: %v", missingFields)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	rh.RequestsStore.Register(&StoredRequest{
		RequestBody:  registerPayload.RequestBody,
		Headers:      registerPayload.Headers,
		ResponseBody: registerPayload.ResponseBody,
		Method:       registerPayload.Method,
		URL:          registerPayload.URL,
	})
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Registed the request successfully!")
}

type UnderTestRequest struct {
	RequestBody json.RawMessage   `json:"requestBody"`
	Headers     map[string]string `json:"headers"`
	Method      string
	URL         string
}

type ApiUnderTestHandler struct {
	RequestsStore *RequestsStore
}

func (h *ApiUnderTestHandler) Handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	matched, storeRequest := h.RequestsStore.MatchRequest(
		&UnderTestRequest{
			RequestBody: body,
			Method:      r.Method,
			Headers:     map_utils.HeaderToMap(r.Header),
			URL:         r.URL.Path,
		},
	)
	if matched {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(storeRequest.ResponseBody)
		return
	}

	http.Error(w, "Unconfigured call", http.StatusNotFound)
}
