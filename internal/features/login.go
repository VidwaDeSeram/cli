package features

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/recode-sh/cli/internal/config"
	"github.com/recode-sh/cli/internal/constants"
	"github.com/recode-sh/cli/internal/exceptions"
	"github.com/recode-sh/cli/internal/interfaces"
	"golang.org/x/oauth2"
)

type LoginInput struct{}

type LoginResponseContent struct{}

type LoginResponse struct {
	Error   error
	Content LoginResponseContent
}

type LoginPresenter interface {
	PresentToView(LoginResponse)
}

type LoginFeature struct {
	presenter  LoginPresenter
	logger     interfaces.Logger
	browser    interfaces.BrowserManager
	userConfig interfaces.UserConfigManager
	sleeper    interfaces.Sleeper
	github     interfaces.GitHubManager
}

func NewLoginFeature(
	presenter LoginPresenter,
	logger interfaces.Logger,
	browser interfaces.BrowserManager,
	config interfaces.UserConfigManager,
	sleeper interfaces.Sleeper,
	github interfaces.GitHubManager,
) LoginFeature {
	return LoginFeature{
		presenter:  presenter,
		logger:     logger,
		browser:    browser,
		userConfig: config,
		sleeper:    sleeper,
		github:     github,
	}
}

func (l LoginFeature) Execute(input LoginInput) error {
	handleError := func(err error) error {
		l.presenter.PresentToView(LoginResponse{
			Error: exceptions.ErrLoginError{
				Reason: err.Error(),
			},
		})

		return err
	}

	gitHubOAuthCbHandlerResp := struct {
		Error       error
		AccessToken string
	}{}

	gitHubOauthCbHandlerDoneChan := make(chan struct{})

	gitHubOAuthCbHandler := func(w http.ResponseWriter, r *http.Request) {
		defer close(gitHubOauthCbHandlerDoneChan)

		queryComponents, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			gitHubOAuthCbHandlerResp.Error = err
			return
		}

		code, hasCodeInQuery := queryComponents["code"]
		if !hasCodeInQuery {
			gitHubOAuthCbHandlerResp.Error = errors.New("no code returned after authorization")
			return
		}

		accessToken, err := ExchangeCodeForToken(code[0])
		if err != nil {
			gitHubOAuthCbHandlerResp.Error = err
			return
		}

		gitHubOAuthCbHandlerResp.AccessToken = accessToken
		w.WriteHeader(200)
		w.Write([]byte("<h1>Success!</h1><p>Your GitHub account is now connected.</p>"))
	}

	http.HandleFunc(
		config.GitHubOAuthAPIToCLIURLPath,
		gitHubOAuthCbHandler,
	)

	httpListener, err := net.Listen("tcp", ":0")
	if err != nil {
		return handleError(err)
	}

	httpServerServeErrorChan := make(chan error, 1)
	go func() {
		httpServerServeErrorChan <- http.Serve(httpListener, nil)
	}()

	httpListenPort := httpListener.Addr().(*net.TCPAddr).Port
	gitHubOAuthClient := &oauth2.Config{
		ClientID:    config.GitHubOAuthClientID,
		Scopes:      config.GitHubOAuthScopes,
		Endpoint:    oauth2.Endpoint{AuthURL: "https://github.com/login/oauth/authorize"},
		RedirectURL: config.GitHubOAuthCLIToAPIURL,
	}
	gitHubOAuthAuthorizeURL := gitHubOAuthClient.AuthCodeURL(fmt.Sprintf("%d", httpListenPort))

	bold := constants.Bold
	l.logger.Log(bold("\nYou will be taken to your browser to connect your GitHub account...\n"))
	l.logger.Info("If your browser doesn't open automatically, go to the following link:\n")
	l.logger.Log("%s", gitHubOAuthAuthorizeURL)
	l.sleeper.Sleep(4 * time.Second)

	if err := l.browser.OpenURL(gitHubOAuthAuthorizeURL); err != nil {
		l.logger.Error("\nCannot open browser! Please visit above URL ↑")
	}

	l.logger.Warning("\nWaiting for GitHub authorization... (Press Ctrl-C to quit)\n")
	select {
	case httpServerServeError := <-httpServerServeErrorChan:
		return handleError(httpServerServeError)
	case <-gitHubOauthCbHandlerDoneChan:
		_ = httpListener.Close()
	}

	if gitHubOAuthCbHandlerResp.Error != nil {
		return handleError(gitHubOAuthCbHandlerResp.Error)
	}

	githubUser, err := l.github.GetAuthenticatedUser(gitHubOAuthCbHandlerResp.AccessToken)
	if err != nil {
		return handleError(err)
	}

	l.userConfig.Set(config.UserConfigKeyUserIsLoggedIn, true)
	l.userConfig.Set(config.UserConfigKeyGitHubAccessToken, gitHubOAuthCbHandlerResp.AccessToken)
	l.userConfig.PopulateFromGitHubUser(githubUser)

	if err := l.userConfig.WriteConfig(); err != nil {
		return handleError(err)
	}

	l.presenter.PresentToView(LoginResponse{})
	return nil
}

func ExchangeCodeForToken(code string) (string, error) {

	postData := url.Values{
		"client_id":     {config.GitHubOAuthClientID},
		"client_secret": {config.GitHubOAuthClientSecret},
		"code":          {code},
	}

	resp, err := http.PostForm("https://github.com/login/oauth/access_token", postData)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK response from GitHub: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("Raw response body:", string(body))

	// Parse the URL-encoded response body
	vals, err := url.ParseQuery(string(body))
	if err != nil {
		return "", fmt.Errorf("error parsing response: %v, response: %s", err, string(body))
	}

	// Extract the access token from the parsed values
	accessToken := vals.Get("access_token")
	if accessToken == "" {
		return "", fmt.Errorf("access token not found in response: %s", string(body))
	}

	fmt.Println("Extracted Token:", accessToken)

	return accessToken, nil
}
