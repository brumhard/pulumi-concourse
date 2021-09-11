package provider

import (
	"context"
	"github.com/concourse/concourse/go-concourse/concourse"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"net/http"
)

var ErrUsernamePasswordRequired = errors.New("username and password are required")

func newPasswordGrantHTTPClient(url, username, password string) (*http.Client, error) {
	if username == "" || password == "" {
		return nil, ErrUsernamePasswordRequired
	}

	// copied from concourse code
	oauth2Config := oauth2.Config{
		ClientID:     "fly",
		ClientSecret: "Zmx5",
		Endpoint:     oauth2.Endpoint{TokenURL: url + "/sky/issuer/token"},
		Scopes:       []string{"openid", "profile", "email", "federated:id", "groups"},
	}

	// TODO: set proper default http.Client
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, http.Client{})

	token, err := oauth2Config.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		return nil, err
	}

	return oauth2Config.Client(ctx, token), nil
}

func NewClient(url, username, password string) (concourse.Client, error) {
	httpClient, err := newPasswordGrantHTTPClient(url, username, password)
	if err != nil {
		return nil, err
	}

	return concourse.NewClient(url, httpClient, false), nil
}
