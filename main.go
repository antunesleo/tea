package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RegisterRequest struct {
	RequestBody  json.RawMessage   `json:"requestBody"`
	Headers      map[string]string `json:"headers"`
	ResponseBody json.RawMessage   `json:"ResponseBody"`
}

type RequestsStore struct {
	requests []*RegisterRequest
}

func (rs *RequestsStore) Register(r *RegisterRequest) {
	rs.requests = append(rs.requests, r)
}

func (rs *RequestsStore) MatchRequest(rr *RegisterRequest) {
	// TODO
	// for _, request := range rs.requests {
	// }
}

type RegisterHandler struct {
	RequestsStore *RequestsStore
}

func (rh *RegisterHandler) Handler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Parse the JSON body into RegisterRequest struct
	var registerReq RegisterRequest
	err = json.Unmarshal(body, &registerReq)
	if err != nil {
		http.Error(w, "Error parsing JSON body", http.StatusBadRequest)
		return
	}

	// Validate mandatory fields
	missingFields := []string{}
	if registerReq.RequestBody == nil {
		missingFields = append(missingFields, "requestBody")
	}
	if registerReq.Headers == nil {
		missingFields = append(missingFields, "headers")
	}
	if registerReq.ResponseBody == nil {
		missingFields = append(missingFields, "ResponseBody")
	}

	// If any mandatory field is missing, return a 400 Bad Request response
	if len(missingFields) > 0 {
		errorMessage := fmt.Sprintf("Missing mandatory fields: %v", missingFields)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	rh.RequestsStore.Register(&registerReq)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Registed the request successfully!")
}

type ApiUnderTestHandler struct {
	RequestsStore *RequestsStore
}

func (h *ApiUnderTestHandler) Handler(w http.ResponseWriter, r *http.Request) {
	return
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
