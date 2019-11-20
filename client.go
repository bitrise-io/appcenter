package appcenter

import "net/http"

type roundTripper struct {
	customHeaders map[string]string
	token         string
}

// RoundTrip ...
func (rt roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add(
		"x-api-token", rt.token,
	)
	req.Header.Add(
		"content-type", "application/json; charset=utf-8",
	)
	for k, v := range rt.customHeaders {
		req.Header.Add(k, v)
	}
	return http.DefaultTransport.RoundTrip(req)
}

// Client ...
type Client struct {
	roundTripper *roundTripper
	httpClient   *http.Client
	debug        bool
}

// Apps ...
func (c Client) Apps(owner, name string) App {
	return App{client: c, owner: owner, name: name}
}

// NewClient returns an AppCenter authenticated client
func NewClient(token string, debug bool) Client {
	rt := &roundTripper{
		token: token,
	}
	return Client{
		roundTripper: rt,
		httpClient: &http.Client{
			Transport: rt,
		},
		debug: debug,
	}
}
