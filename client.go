package appcenter

import "net/http"

type roundTripper struct {
	token string
}

// RoundTrip ...
func (rt roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add(
		"x-api-token", rt.token,
	)
	req.Header.Add(
		"content-type", "application/json; charset=utf-8",
	)
	return http.DefaultTransport.RoundTrip(req)
}

// Client ...
type Client struct {
	httpClient *http.Client
	debug      bool
}

// Apps ...
func (c Client) Apps(owner, name string) App {
	return App{client: c, owner: owner, name: name}
}

// NewClient returns an AppCenter authenticated client
func NewClient(token string, debug bool) Client {
	return Client{
		httpClient: &http.Client{
			Transport: &roundTripper{
				token: token,
			},
		},
		debug: debug,
	}
}
