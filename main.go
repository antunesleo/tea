package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

func mapsEqual(a, b map[string]string) bool {
	for keyA, valueA := range a {
		// Check if the key exists in map b
		if valueB, ok := b[keyA]; ok {
			// If the values are not equal, maps are not equal
			if valueA != valueB {
				return false
			}
		} else {
			// If the key doesn't exist in map b, maps are not equal
			return false
		}
	}
	// All keys in map a exist in map b with equal values
	return true
}

func HeaderToMap(header http.Header) map[string]string {
	result := make(map[string]string)
	for key, values := range header {
		// Use the first value, assuming you only want a single string value for each key
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}

type RequestsStore struct {
	storedRequests []*StoredRequest
}

func (rs *RequestsStore) Register(storedRequest *StoredRequest) {
	rs.storedRequests = append(rs.storedRequests, storedRequest)
}

func (rs *RequestsStore) MatchRequest(underTestReq *UnderTestRequest) (bool, StoredRequest) {
	for _, storedReq := range rs.storedRequests {
		if storedReq.Method != underTestReq.Method {
			continue
		}

		var data1, data2 interface{}
		if err := json.Unmarshal(storedReq.RequestBody, &data1); err != nil {
			fmt.Println("Error unmarshaling rawMessage1:", err)
			return false, StoredRequest{} // TODO: refactor to return error
		}

		if err := json.Unmarshal(underTestReq.RequestBody, &data2); err != nil {
			fmt.Println("Error unmarshaling rawMessage2:", err)
			return false, StoredRequest{} // TODO: refactor to return error
		}

		if !reflect.DeepEqual(data1, data2) {
			continue
		}

		if !mapsEqual(storedReq.Headers, underTestReq.Headers) {
			continue
		}

		return true, *storedReq
	}
	return false, StoredRequest{}
}

type StoredRequest struct {
	RequestBody  json.RawMessage
	Headers      map[string]string
	ResponseBody json.RawMessage
	Method       string
}

type UnderTestRequest struct {
	RequestBody json.RawMessage   `json:"requestBody"`
	Headers     map[string]string `json:"headers"`
	Method      string
}

type RegisterHandler struct {
	RequestsStore *RequestsStore
}

type RegisterHandlerPayload struct {
	RequestBody  json.RawMessage   `json:"requestBody"`
	Headers      map[string]string `json:"headers"`
	ResponseBody json.RawMessage   `json:"ResponseBody"`
}

func (rh *RegisterHandler) Handler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Parse the JSON body into RegisterRequest struct
	var registerPayload RegisterHandlerPayload
	err = json.Unmarshal(body, &registerPayload)
	if err != nil {
		errorMsg := fmt.Sprint("Error parsing JSON body: %v", err)
		http.Error(w, errorMsg, http.StatusBadRequest)
		return
	}

	// Validate mandatory fields
	missingFields := []string{}
	if registerPayload.RequestBody == nil {
		missingFields = append(missingFields, "requestBody")
	}
	if registerPayload.Headers == nil {
		missingFields = append(missingFields, "headers")
	}
	if registerPayload.ResponseBody == nil {
		missingFields = append(missingFields, "ResponseBody")
	}

	// If any mandatory field is missing, return a 400 Bad Request response
	if len(missingFields) > 0 {
		errorMessage := fmt.Sprintf("Missing mandatory fields: %v", missingFields)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	rh.RequestsStore.Register(&StoredRequest{
		RequestBody:  registerPayload.RequestBody,
		Headers:      registerPayload.Headers,
		ResponseBody: registerPayload.ResponseBody,
		Method:       r.Method,
	})
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Registed the request successfully!")
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
		&UnderTestRequest{RequestBody: body, Method: r.Method, Headers: HeaderToMap(r.Header)},
	)
	if matched {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(storeRequest.ResponseBody)
		return
	}

	http.Error(w, "Unconfigured call", http.StatusNotFound)
}

func main() {
	rs := &RequestsStore{}
	registerHandler := RegisterHandler{rs}
	apiUnderTestHandler := ApiUnderTestHandler{rs}
	http.HandleFunc("/register-request", registerHandler.Handler)
	http.HandleFunc("/", apiUnderTestHandler.Handler)
	fmt.Println("Server is running on port 7111")
	err := http.ListenAndServe(":7111", nil)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
