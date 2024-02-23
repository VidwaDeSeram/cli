package config

var (
	GitHubOAuthClientID     = "aa380e0c7a818c10acc9"
	GitHubOAuthClientSecret = "5b2b21756b41c0fc90e010bdbed4ec4b5b0b3d50"
	GitHubOAuthCLIToAPIURL  = "http://127.0.0.1:8080/github/oauth/callback"

	GitHubOAuthAPIToCLIURLPath = "/github/oauth/callback"

	GitHubOAuthScopes = []string{
		"read:user",
		"user:email",
		"repo",
		"admin:public_key",
		"admin:gpg_key",
	}
)
