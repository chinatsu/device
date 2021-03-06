package auth

import (
	"context"
	"fmt"
	"github.com/nais/device/device-agent/open"
	"net/http"
	"strings"
	"time"

	"github.com/nais/device/pkg/random"
	codeverifier "github.com/nirasan/go-oauth-pkce-code-verifier"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func AzureAuthenticatedClient(ctx context.Context, conf oauth2.Config) (*http.Client, error) {
	token, err := runAuthFlow(ctx, conf)

	if err != nil {
		return nil, fmt.Errorf("running authorization code flow: %w", err)
	}

	return conf.Client(ctx, token), nil
}

func runAuthFlow(ctx context.Context, conf oauth2.Config) (*oauth2.Token, error) {
	// Ignoring impossible error
	codeVerifier, _ := codeverifier.CreateCodeVerifier()

	// TODO check this in response from Azure
	tokenChan := make(chan *oauth2.Token)
	handler := http.NewServeMux()

	// define a handler that will get the authorization code, call the token endpoint, and close the HTTP server
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Catch if user has not approved terms
		if strings.HasPrefix(r.URL.Query().Get("error_description"), "AADSTS50105") {
			http.Redirect(w, r, "https://naisdevice-approval.nais.io/", http.StatusSeeOther)
			tokenChan <- nil
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			log.Errorf("Error: could not find 'code' URL query parameter")
			failureResponse(w, "Error: could not find 'code' URL query parameter")
			tokenChan <- nil
			return
		}

		ctx, cancel := context.WithDeadline(ctx, time.Now().Add(30*time.Second))
		defer cancel()

		codeVerifierParam := oauth2.SetAuthURLParam("code_verifier", codeVerifier.String())
		t, err := conf.Exchange(ctx, code, codeVerifierParam)
		if err != nil {
			failureResponse(w, "Error: during code exchange")
			log.Errorf("exchanging code for tokens: %v", err)
			tokenChan <- nil
			return
		}

		successfulResponse(w, "Successfully authenticated 👌 Close me pls")
		tokenChan <- t
	})

	server := &http.Server{Addr: "127.0.0.1:51800", Handler: handler}
	go server.ListenAndServe()
	defer server.Close()

	url := conf.AuthCodeURL(
		random.RandomString(16, random.LettersAndNumbers),
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeVerifier.CodeChallengeS256()))

	err := open.Open(url)
	if err != nil {
		log.Errorf("opening browser, err: %v", err)
		// Don't return, as this is not fatal (user can open browser manually)
	}
	fmt.Printf("If the browser didn't open, visit this url to sign in: %v\n", url)

	token := <-tokenChan

	if token == nil {
		return nil, fmt.Errorf("no token received")
	}

	return token, nil
}
