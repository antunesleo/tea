package main

import (
	"fmt"
	"net/http"

	"github.com/antunesleo/tea/internal/tea"
)

func main() {
	rs := &tea.RequestsStore{}
	registerHandler := tea.RegisterHandler{rs}
	apiUnderTestHandler := tea.ApiUnderTestHandler{rs}
	http.HandleFunc("/register-request", registerHandler.Handler)
	http.HandleFunc("/", apiUnderTestHandler.Handler)
	fmt.Println("Server is running on port 7111")
	err := http.ListenAndServe(":7111", nil)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
