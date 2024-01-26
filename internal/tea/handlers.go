package tea

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/antunesleo/tea/internal/map_utils"
)

type ErrorResponse struct {
	Error   string `json: "error"`
	Message string `json: "message"`
}

type ExpectedRequestPayload struct {
	Body    json.RawMessage   `json:"body"`
	Headers map[string]string `json:"headers"`
	Method  string            `json:"method"`
	URL     string            `json:"url"`
}

type WantedResponsePayload struct {
	Body       json.RawMessage   `json:"body"`
	StatusCode int               `json:"statusCode`
	Headers    map[string]string `json:"headers"`
}

type RegisterHandlerPayload struct {
	ExpectedRequest ExpectedRequestPayload `json:"expectedRequest"`
	WantedResponse  WantedResponsePayload  `json:"wantedResponse"`
}

type RegisterHandler struct {
	RequestsStore *RequestsStore
}

func (rh *RegisterHandler) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var registerPayload RegisterHandlerPayload
	err := json.NewDecoder(r.Body).Decode(&registerPayload)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	missingFields := rh.validateMissingFields(registerPayload)

	if len(missingFields) > 0 {
		errorMessage := fmt.Sprintf("Missing mandatory fields: %v", missingFields)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	rh.RequestsStore.Register(&StoredRequest{
		ExpectedRequest: ExpectedRequest{
			Body:    registerPayload.ExpectedRequest.Body,
			Headers: registerPayload.ExpectedRequest.Headers,
			Method:  registerPayload.ExpectedRequest.Method,
			URL:     registerPayload.ExpectedRequest.URL,
		},
		WantedResponse: WantedResponse{
			Body:       registerPayload.WantedResponse.Body,
			StatusCode: registerPayload.WantedResponse.StatusCode,
			Headers:    registerPayload.WantedResponse.Headers,
		},
	})
	w.WriteHeader(http.StatusCreated)
}

func NewRegisterHandler(rs *RequestsStore) *RegisterHandler {
	return &RegisterHandler{RequestsStore: rs}
}

func (*RegisterHandler) validateMissingFields(registerPayload RegisterHandlerPayload) []string {
	missingFields := []string{}
	if registerPayload.ExpectedRequest.URL == "" {
		missingFields = append(missingFields, "expectedRequest.url")
	}
	if registerPayload.ExpectedRequest.Method == "" {
		missingFields = append(missingFields, "expectedRequest.method")
	}
	if registerPayload.WantedResponse.Body == nil {
		missingFields = append(missingFields, "wantedResponse.body")
	}
	if registerPayload.WantedResponse.StatusCode == 0 {
		missingFields = append(missingFields, "wantedResponse.statusCode")
	}
	return missingFields
}

type UnderTestRequest struct {
	RequestBody json.RawMessage   `json:"requestBody"`
	Headers     map[string]string `json:"headers"`
	Method      string            `json:"method"`
	URL         string            `json:url`
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
	storeRequest, err := h.RequestsStore.MatchRequest(
		&UnderTestRequest{
			RequestBody: body,
			Method:      r.Method,
			Headers:     map_utils.HeaderToMap(r.Header),
			URL:         r.URL.Path,
		},
	)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(storeRequest.WantedResponse.StatusCode)
		w.Write(storeRequest.WantedResponse.Body)
		for key, value := range storeRequest.WantedResponse.Headers {
			w.Header().Set(key, value)
		}
		return
	}

	http.Error(w, "Unconfigured call", http.StatusNotFound)
}

func NewApiUnderTestHandler(rs *RequestsStore) *ApiUnderTestHandler {
	return &ApiUnderTestHandler{RequestsStore: rs}
}
