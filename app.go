package appcenter

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"
)

// App ...
type App struct {
	client      Client
	owner, name string
}

// NewRelease ...
func (a App) NewRelease(filePath string, opts ...ReleaseOptions) (Release, error) {
	// send file upload request
	var (
		postURL      = fmt.Sprintf("%s/v0.1/apps/%s/%s/release_uploads", baseURL, a.owner, a.name)
		postBody     interface{}
		postResponse struct {
			UploadID  string `json:"upload_id"`
			UploadURL string `json:"upload_url"`
		}
	)

	if len(opts) > 0 {
		postBody = &opts[0]
	}

	statusCode, err := a.client.jsonRequest(http.MethodPost, postURL, postBody, &postResponse)
	if err != nil {
		return Release{}, err
	}

	if statusCode != http.StatusCreated {
		return Release{}, fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, postURL, postBody)
	}

	// upload file to {upload_url}
	statusCode, err = a.client.uploadRequest(postResponse.UploadURL, map[string]string{"ipa": filePath})
	if err != nil {
		return Release{}, err
	}

	if statusCode != http.StatusNoContent {
		return Release{}, fmt.Errorf("invalid status code: %d, url: %s", statusCode, postResponse.UploadURL)
	}

	var (
		patchURL  = fmt.Sprintf("%s/v0.1/apps/%s/%s/release_uploads/%s", baseURL, a.owner, a.name, postResponse.UploadID)
		patchBody = map[string]string{
			"status": "committed",
		}
		patchResponse struct {
			ReleaseID  string `json:"release_id"`
			ReleaseURL string `json:"release_url"`
		}
	)

	statusCode, err = a.client.jsonRequest(http.MethodPatch, patchURL, patchBody, &patchResponse)
	if err != nil {
		return Release{}, err
	}

	if statusCode != http.StatusOK {
		return Release{}, fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, patchURL, patchResponse)
	}

	// fetch release details
	var (
		getURL      = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases/%s", baseURL, a.owner, a.name, patchResponse.ReleaseID)
		getResponse Release
	)

	statusCode, err = a.client.jsonRequest(http.MethodGet, getURL, nil, &getResponse)
	if err != nil {
		return Release{}, err
	}

	if statusCode != http.StatusOK {
		return Release{}, fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, getURL, getResponse)
	}

	getResponse.app = a

	return getResponse, nil
}

// Groups ...
func (a App) Groups(name string) (Group, error) {
	var (
		getURL      = fmt.Sprintf("%s/v0.1/apps/%s/%s/distribution_groups/%s", baseURL, a.owner, a.name, name)
		getResponse Group
	)

	statusCode, err := a.client.jsonRequest(http.MethodGet, getURL, nil, &getResponse)
	if err != nil {
		return Group{}, err
	}

	if statusCode != http.StatusOK {
		return Group{}, fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, getURL, getResponse)
	}

	return getResponse, nil
}

// Stores ...
func (a App) Stores(name string) (Store, error) {
	var (
		getURL      = fmt.Sprintf("%s/v0.1/apps/%s/%s/distribution_stores/%s", baseURL, a.owner, a.name, name)
		getResponse Store
	)

	statusCode, err := a.client.jsonRequest(http.MethodGet, getURL, nil, &getResponse)
	if err != nil {
		return Store{}, err
	}

	if statusCode != http.StatusOK {
		return Store{}, fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, getURL, getResponse)
	}

	return getResponse, nil
}

// UploadSymbol - build and version is required for Android and optional for iOS
func (a App) UploadSymbol(filePath, build, version string, symbolType SymbolType) error {
	// send file upload request
	var (
		postURL  = fmt.Sprintf("%s/v0.1/apps/%s/%s/symbol_uploads", baseURL, a.owner, a.name)
		postBody = struct {
			SymbolType SymbolType `json:"symbol_type"`
			FileName   string     `json:"file_name,omitempty"`
			Build      string     `json:"build,omitempty"`
			Version    string     `json:"version,omitempty"`
		}{
			FileName: filepath.Base(filePath),
			Build:    build,
			Version:  version,
		}
		postResponse struct {
			SymbolUploadID string    `json:"symbol_upload_id"`
			UploadURL      string    `json:"upload_url"`
			ExpirationDate time.Time `json:"expiration_date"`
		}
	)

	statusCode, err := a.client.jsonRequest(http.MethodPost, postURL, postBody, &postResponse)
	if err != nil {
		return err
	}

	if statusCode != http.StatusCreated {
		return fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, postURL, postBody)
	}

	// upload file to {upload_url}
	statusCode, err = a.client.uploadRequest(postResponse.UploadURL, map[string]string{"dsym": filePath})
	if err != nil {
		return err
	}

	if statusCode != http.StatusNoContent {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, postResponse.UploadURL)
	}

	var (
		patchURL  = fmt.Sprintf("%s/v0.1/apps/%s/%s/symbol_uploads/%s", baseURL, a.owner, a.name, postResponse.SymbolUploadID)
		patchBody = map[string]string{
			"status": "committed",
		}
	)

	statusCode, err = a.client.jsonRequest(http.MethodPatch, patchURL, patchBody, nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, patchURL)
	}

	return nil
}
