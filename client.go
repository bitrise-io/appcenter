package appcenter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
)

const (
	baseURL = `https://api.appcenter.ms`
)

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

func (c Client) jsonRequest(method, url string, body interface{}, response interface{}) (int, error) {
	var reader io.Reader

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return -1, err
		}
		reader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return -1, err
	}

	if c.debug {
		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			return -1, err
		}
		fmt.Println(string(b))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return -1, err
	}

	if c.debug {
		b, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return -1, err
		}
		fmt.Println(string(b))
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
		}
	}()

	if response != nil {
		rb, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return -1, err
		}

		if err := json.Unmarshal(rb, response); err != nil {
			return resp.StatusCode, fmt.Errorf("error: %s, response: %s", err, string(rb))
		}
	}

	return resp.StatusCode, nil
}

func (c Client) uploadFile(url string, filePath string) (int, error) {
	fb, err := ioutil.ReadFile(filePath)
	if err != nil {
		return -1, err
	}

	uploadReq, err := http.NewRequest("PUT", url, bytes.NewReader(fb))
	if err != nil {
		return -1, err
	}

	uploadReq.Header.Set("x-ms-blob-type", "BlockBlob")
	uploadReq.Header.Set("content-length", strconv.Itoa(len(fb)))

	resp, err := c.httpClient.Do(uploadReq)
	if err != nil {
		return -1, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
		}
	}()

	return resp.StatusCode, nil
}

// API ...
type API struct {
	Client Client
}

// GetLatestReleases ...
func (api API) GetLatestReleases(app App) (Release, error) {
	//fetch releases and find the latest
	var (
		getReleasesURL   = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases", baseURL, app.owner, app.name)
		releasesResponse []Release
	)

	statusCode, err := api.Client.jsonRequest(http.MethodGet, getReleasesURL, nil, &releasesResponse)
	if err != nil {
		return Release{}, err
	}

	if statusCode != http.StatusOK {
		return Release{}, fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, getReleasesURL, releasesResponse)
	}

	latestReleaseID := releasesResponse[0].ID

	var (
		releaseShowURL = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases/%s", baseURL, app.owner, app.name, strconv.Itoa(latestReleaseID))
		release        Release
	)

	statusCode, err = api.Client.jsonRequest(http.MethodGet, releaseShowURL, nil, &release)
	if err != nil {
		return Release{}, err
	}

	if statusCode != http.StatusOK {
		return Release{}, fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, releaseShowURL, release)
	}

	return release, err
}

// GetGroupByName ...
func (api API) GetGroupByName(groupName string, app App) (Group, error) {
	var (
		getURL      = fmt.Sprintf("%s/v0.1/apps/%s/%s/distribution_groups/%s", baseURL, app.owner, app.name, groupName)
		getResponse Group
	)

	statusCode, err := api.Client.jsonRequest(http.MethodGet, getURL, nil, &getResponse)
	if err != nil {
		return Group{}, err
	}

	if statusCode != http.StatusOK {
		return Group{}, fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, getURL, getResponse)
	}

	return getResponse, err
}

// GetStore ...
func (api API) GetStore(storeName string, app App) (Store, error) {
	var (
		getURL      = fmt.Sprintf("%s/v0.1/apps/%s/%s/distribution_stores/%s", baseURL, app.owner, app.name, storeName)
		getResponse Store
	)

	statusCode, err := api.Client.jsonRequest(http.MethodGet, getURL, nil, &getResponse)
	if err != nil {
		return Store{}, err
	}

	if statusCode != http.StatusOK {
		return Store{}, fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, getURL, getResponse)
	}

	return getResponse, nil
}
