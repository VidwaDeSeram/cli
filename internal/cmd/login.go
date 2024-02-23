package cmd

import (
	"fmt"
	"net/http"

	"github.com/recode-sh/cli/internal/dependencies"
	"github.com/recode-sh/cli/internal/features"
	"github.com/spf13/cobra"
)

func TokenCallBack(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := features.ExchangeCodeForToken(code)
	if err != nil {
		// Handle the error, maybe return it in the response
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Now you can print the token
	fmt.Println("Token:", token)

	// Add additional logic to handle the OAuth code and token
}

// LoginHandler handles the login requests
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	login := dependencies.ProvideLoginFeature()
	loginInput := features.LoginInput{}

	// Execute the login feature which initiates the OAuth flow
	err := login.Execute(loginInput)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Normally, the login feature would redirect to the OAuth URL
	// For API, return the URL in the response instead
	// This is a placeholder, replace with actual OAuth URL generation logic
	oAuthURL := "https://github.com/login/oauth/authorize"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"url": "` + oAuthURL + `"}`))

}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to your GitHub account",
	Long:  `Log in to your GitHub account via OAuth.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Handler logic here
	},
}

func init() {

	rootCmd.AddCommand(loginCmd)
}
