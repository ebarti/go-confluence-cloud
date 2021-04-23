package confluentcloud

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// NewAPI implements api constructor
func NewAPI(location string, username string, token string) (API, error) {
	return newAPI(location, username, token)
}

func newAPI(location string, username string, token string) (*api, error) {
	if len(location) == 0 {
		return nil, errors.New("url empty")
	}

	u, err := url.ParseRequestURI(location)

	if err != nil {
		return nil, err
	}

	a := new(api)
	a.endPoint = u
	a.token = token
	a.username = username

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}

	a.client = &http.Client{Transport: tr}

	return a, nil
}

// NewAPIWithClient creates a new api instance using an existing HTTP client.
// Useful when using oauth or other authentication methods.
func NewAPIWithClient(location string, client *http.Client) (API, error) {
	u, err := url.ParseRequestURI(location)

	if err != nil {
		return nil, err
	}

	a := new(api)
	a.endPoint = u
	a.client = client

	return a, nil
}

// VerifyTLS to enable disable certificate checks
func (a *api) VerifyTLS(set bool) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !set},
	}
	a.client = &http.Client{Transport: tr}
}

// DebugFlag is the global debugging variable
var DebugFlag = false

// SetDebug enables debug output
func SetDebug(state bool) {
	DebugFlag = state
}

// Debug outputs debug messages
func Debug(msg interface{}) {
	if DebugFlag {
		fmt.Printf("%+v\n", msg)
	}
}
