package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/recode-sh/cli/internal/cmd"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/login", cmd.LoginHandler).Methods("POST")
	router.HandleFunc("/github/oauth/callback", cmd.TokenCallBack).Methods("get")

	// Add other routes here

	http.ListenAndServe(":8080", router)
}
