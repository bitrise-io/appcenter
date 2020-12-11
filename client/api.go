package client

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/bitrise-io/appcenter/model"
)

// CreateAPIWithClientParams ...
func CreateAPIWithClientParams(token string, debug bool) API {
	return API{
		Client: NewClient(token, debug),
	}
}

// API ...
type API struct {
	Client Client
}

// GetLatestReleases ...
func (api API) GetLatestReleases(app model.App) (model.Release, error) {
	//fetch releases and find the latest
	var (
		getReleasesURL   = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases", baseURL, app.Owner, app.AppName)
		releasesResponse []model.Release
	)

	statusCode, err := api.Client.jsonRequest(http.MethodGet, getReleasesURL, nil, &releasesResponse)
	if err != nil {
		return model.Release{}, err
	}

	if statusCode != http.StatusOK {
		return model.Release{}, fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, getReleasesURL, releasesResponse)
	}

	latestReleaseID := releasesResponse[0].ID

	var (
		releaseShowURL = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases/%s", baseURL, app.Owner, app.AppName, strconv.Itoa(latestReleaseID))
		release        model.Release
	)

	statusCode, err = api.Client.jsonRequest(http.MethodGet, releaseShowURL, nil, &release)
	if err != nil {
		return model.Release{}, err
	}

	if statusCode != http.StatusOK {
		return model.Release{}, fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, releaseShowURL, release)
	}

	return release, err
}

// GetGroupByName ...
func (api API) GetGroupByName(groupName string, app model.App) (model.Group, error) {
	var (
		getURL      = fmt.Sprintf("%s/v0.1/apps/%s/%s/distribution_groups/%s", baseURL, app.Owner, app.AppName, groupName)
		getResponse model.Group
	)

	statusCode, err := api.Client.jsonRequest(http.MethodGet, getURL, nil, &getResponse)
	if err != nil {
		return model.Group{}, err
	}

	if statusCode != http.StatusOK {
		return model.Group{}, fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, getURL, getResponse)
	}

	return getResponse, err
}

// GetStore ...
func (api API) GetStore(storeName string, app model.App) (model.Store, error) {
	var (
		getURL      = fmt.Sprintf("%s/v0.1/apps/%s/%s/distribution_stores/%s", baseURL, app.Owner, app.AppName, storeName)
		getResponse model.Store
	)

	statusCode, err := api.Client.jsonRequest(http.MethodGet, getURL, nil, &getResponse)
	if err != nil {
		return model.Store{}, err
	}

	if statusCode != http.StatusOK {
		return model.Store{}, fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, getURL, getResponse)
	}

	return getResponse, nil
}

// AddReleaseToGroup ...
func (api API) AddReleaseToGroup(g model.Group, releaseID int, opts model.ReleaseOptions) error {
	var (
		postURL     = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases/%d/groups", baseURL, opts.App.Owner, opts.App.AppName, releaseID)
		postRequest = struct {
			ID              string `json:"id"`
			MandatoryUpdate bool   `json:"mandatory_update"`
			NotifyTesters   bool   `json:"notify_testers"`
		}{
			ID:              g.ID,
			MandatoryUpdate: opts.Mandatory,
			NotifyTesters:   opts.NotifyTesters,
		}
	)

	statusCode, err := api.Client.jsonRequest(http.MethodPost, postURL, postRequest, nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusCreated {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, postURL)
	}

	return nil
}

// AddReleaseToStore ...
func (api API) AddReleaseToStore(s model.Store, releaseID int, opts model.ReleaseOptions) error {
	var (
		postURL     = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases/%d/stores", baseURL, opts.App.Owner, opts.App.AppName, releaseID)
		postRequest = struct {
			ID string `json:"id"`
		}{
			ID: s.ID,
		}
	)

	statusCode, err := api.Client.jsonRequest(http.MethodPost, postURL, postRequest, nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusCreated {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, postURL)
	}

	return nil
}

// AddTesterToRelease ...
func (api API) AddTesterToRelease(email string, releaseID int, opts model.ReleaseOptions) error {
	var (
		postURL     = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases/%d/testers", baseURL, opts.App.Owner, opts.App.AppName, releaseID)
		postRequest = struct {
			Email           string `json:"email"`
			MandatoryUpdate bool   `json:"mandatory_update"`
			NotifyTesters   bool   `json:"notify_testers"`
		}{
			Email:           email,
			MandatoryUpdate: opts.Mandatory,
			NotifyTesters:   opts.NotifyTesters,
		}
	)

	statusCode, err := api.Client.jsonRequest(http.MethodPost, postURL, postRequest, nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusCreated {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, postURL)
	}

	return nil
}

// SetReleaseNoteOnRelease ...
func (api API) SetReleaseNoteOnRelease(releaseNote string, releaseID int, opts model.ReleaseOptions) error {
	var (
		putURL     = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases/%d", baseURL, opts.App.Owner, opts.App.AppName, releaseID)
		putRequest = struct {
			ReleaseNotes string `json:"release_notes,omitempty"`
		}{
			ReleaseNotes: releaseNote,
		}
	)

	statusCode, err := api.Client.jsonRequest(http.MethodPut, putURL, putRequest, nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, putURL)
	}

	return nil
}

// UploadSymbolToRelease - build and version is required for Android and optional for iOS
func (api API) UploadSymbolToRelease(filePath string, release model.Release, opts model.ReleaseOptions) error {
	var symbolType = model.SymbolTypeDSYM
	if release.AppOs == "Android" {
		symbolType = model.SymbolTypeMapping
	}

	// send file upload request
	var (
		postURL  = fmt.Sprintf("%s/v0.1/apps/%s/%s/symbol_uploads", baseURL, opts.App.Owner, opts.App.AppName)
		postBody = struct {
			SymbolType model.SymbolType `json:"symbol_type"`
			FileName   string           `json:"file_name,omitempty"`
			Build      string           `json:"build,omitempty"`
			Version    string           `json:"version,omitempty"`
		}{
			FileName:   filepath.Base(filePath),
			Build:      release.Version,
			Version:    release.ShortVersion,
			SymbolType: symbolType,
		}
		postResponse struct {
			SymbolUploadID string    `json:"symbol_upload_id"`
			UploadURL      string    `json:"upload_url"`
			ExpirationDate time.Time `json:"expiration_date"`
		}
	)

	statusCode, err := api.Client.jsonRequest(http.MethodPost, postURL, postBody, &postResponse)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, postURL, postBody)
	}

	// upload file to {upload_url}
	statusCode, err = api.Client.uploadFile(postResponse.UploadURL, filePath)
	if err != nil {
		return err
	}

	if statusCode != http.StatusCreated {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, postResponse.UploadURL)
	}

	var (
		patchURL  = fmt.Sprintf("%s/v0.1/apps/%s/%s/symbol_uploads/%s", baseURL, opts.App.Owner, opts.App.AppName, postResponse.SymbolUploadID)
		patchBody = map[string]string{
			"status": "committed",
		}
	)

	statusCode, err = api.Client.jsonRequest(http.MethodPatch, patchURL, patchBody, nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, patchURL)
	}

	return nil
}
