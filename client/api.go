package client

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/bitrise-io/appcenter/util"

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

// GetReleaseOnAppByID ...
func (api API) GetReleaseOnAppByID(app model.App, releaseID int) (model.Release, error) {
	//fetch releases and find the latest
	var (
		releaseShowURL = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases/%s", baseURL, app.Owner, app.AppName, strconv.Itoa(releaseID))
		release        model.Release
	)

	statusCode, err := api.Client.jsonRequest(http.MethodGet, releaseShowURL, nil, &release)
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

	body, err := api.Client.MarshallContent(postRequest)
	if err != nil {
		return err
	}

	statusCode, err := api.Client.jsonRequest(http.MethodPost, postURL, body, nil)
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

	body, err := api.Client.MarshallContent(postRequest)
	if err != nil {
		return err
	}

	statusCode, err := api.Client.jsonRequest(http.MethodPost, postURL, body, nil)
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

	body, err := api.Client.MarshallContent(postRequest)
	if err != nil {
		return err
	}

	statusCode, err := api.Client.jsonRequest(http.MethodPost, postURL, body, nil)
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

	body, err := api.Client.MarshallContent(putRequest)
	if err != nil {
		return err
	}

	statusCode, err := api.Client.jsonRequest(http.MethodPut, putURL, body, nil)
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

	body, err := api.Client.MarshallContent(postBody)
	if err != nil {
		return err
	}

	statusCode, err := api.Client.jsonRequest(http.MethodPost, postURL, body, &postResponse)
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

	body, err = api.Client.MarshallContent(patchBody)
	if err != nil {
		return err
	}

	statusCode, err = api.Client.jsonRequest(http.MethodPatch, patchURL, body, nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, patchURL)
	}

	return nil
}

// CreateRelease ...
func (api API) CreateRelease(opts model.ReleaseOptions) (model.Release, error) {
	var (
		assetsURL = fmt.Sprintf("%s/v0.1/apps/%s/%s/uploads/releases",
			baseURL,
			opts.App.Owner,
			opts.App.AppName)
		fileAssetsResponse struct {
			ReleaseID       string `json:"id"`
			PackageAssetID  string `json:"package_asset_id"`
			Token           string `json:"token"`
			UploadDomain    string `json:"upload_domain"`
			URLEncodedToken string `json:"url_encoded_token"`
		}
	)

	statusCode, err := api.Client.jsonRequest(http.MethodPost, assetsURL, nil, &fileAssetsResponse)
	if err != nil {
		return model.Release{}, err
	}

	if statusCode != http.StatusCreated {
		return model.Release{}, fmt.Errorf("invalid status code: %d, url: %s", statusCode, assetsURL)
	}

	fmt.Println("File assets response:")
	fmt.Println(fileAssetsResponse)

	file := util.LocalFile{FilePath: opts.FilePath}
	err = file.OpenFile()
	if err != nil {
		return model.Release{}, err
	}

	fileName := file.FileName()
	fileSize := file.FileSize()

	fmt.Println("Uploading file with metadata:")
	fmt.Println(fmt.Sprintf("- File name: %s", fileName))
	fmt.Println(fmt.Sprintf("- File size: %s", strconv.Itoa(fileSize)))

	var (
		metadataURL = fmt.Sprintf("%s/upload/set_metadata/%s?file_name=%s&file_size=%s&token=%s",
			fileAssetsResponse.UploadDomain,
			fileAssetsResponse.PackageAssetID,
			fileName,
			strconv.Itoa(fileSize),
			fileAssetsResponse.URLEncodedToken)
		metadataResponse struct {
			ID             string `json:"id"`
			ChunkSize      int    `json:"chunk_size"`
			ChunkList      []int  `json:"chunk_list"`
			BlobPartitions int    `json:"blob_partitions"`
		}
	)

	statusCode, err = api.Client.jsonRequest(http.MethodPost, metadataURL, nil, &metadataResponse)
	if err != nil {
		return model.Release{}, err
	}

	if statusCode != http.StatusOK {
		return model.Release{}, fmt.Errorf("invalid status code: %d, url: %s", statusCode, assetsURL)
	}

	fmt.Println("File assets response:")
	fmt.Println(metadataResponse)

	fmt.Println("Uploading chunks ...")

	fileChunks := file.MakeChunks(metadataResponse.ChunkSize)

	for idx, chunkID := range metadataResponse.ChunkList {
		chunk := fileChunks[idx]
		fmt.Println(fmt.Sprintf("Chunk ID: %d, chunk size: %d", chunkID, len(chunk)))

		var (
			chunkUploadURL = fmt.Sprintf("%s/upload/upload_chunk/%s?block_number=%s&token=%s",
				fileAssetsResponse.UploadDomain,
				fileAssetsResponse.PackageAssetID,
				strconv.Itoa(chunkID),
				fileAssetsResponse.URLEncodedToken)
			chunkUploadResponse interface{}
		)

		statusCode, err = api.Client.jsonRequest(http.MethodPost, chunkUploadURL, chunk, &chunkUploadResponse)
		if err != nil {
			return model.Release{}, err
		}

		if statusCode != http.StatusOK {
			return model.Release{}, fmt.Errorf("invalid status code: %d, url: %s", statusCode, assetsURL)
		}

		fmt.Println(chunkUploadResponse)
	}

	fmt.Println("Chunk upload finished...")

	var (
		uploadFinishedURL = fmt.Sprintf("%s/upload/finished/%s?token=%s",
			fileAssetsResponse.UploadDomain,
			fileAssetsResponse.PackageAssetID,
			fileAssetsResponse.URLEncodedToken)
		finishedResponse interface{}
	)

	statusCode, err = api.Client.jsonRequest(http.MethodPost, uploadFinishedURL, nil, &finishedResponse)
	if err != nil {
		return model.Release{}, err
	}

	if statusCode != http.StatusOK {
		return model.Release{}, fmt.Errorf("invalid status code: %d, url: %s", statusCode, assetsURL)
	}

	//patch release

	var (
		releasePatchURL = fmt.Sprintf("%s/v0.1/apps/%s/%s/uploads/releases/%s",
			baseURL,
			opts.App.Owner,
			opts.App.AppName,
			fileAssetsResponse.ReleaseID)
		releaseBody = struct {
			UploadStatus string `json:"upload_status"`
		}{
			UploadStatus: "uploadFinished",
		}
		releasePatchResponse interface{}
	)

	body, err := api.Client.MarshallContent(releaseBody)
	if err != nil {
		return model.Release{}, err
	}

	statusCode, err = api.Client.jsonRequest(http.MethodPatch, releasePatchURL, body, &releasePatchResponse)
	if err != nil {
		return model.Release{}, err
	}

	fmt.Println(releasePatchResponse)

	if statusCode != http.StatusOK {
		return model.Release{}, fmt.Errorf("invalid status code: %d, url: %s", statusCode, assetsURL)
	}

	uploadStatus := "commited"
	releaseDistinctID := -1
	for ok := true; ok; ok = uploadStatus != "readyToBePublished" {
		var (
			getURL = fmt.Sprintf("%s/v0.1/apps/%s/%s/uploads/releases/%s",
				baseURL,
				opts.App.Owner,
				opts.App.AppName,
				fileAssetsResponse.ReleaseID)
			getResponse struct {
				ID                string `json:"id"`
				ReleaseDistinctID int    `json:"release_distinct_id,omitempty"`
				UploadStatus      string `json:"upload_status"`
			}
		)

		statusCode, err = api.Client.jsonRequest(http.MethodGet, getURL, nil, &getResponse)
		if err != nil {
			return model.Release{}, err
		}

		fmt.Println(getResponse)

		if statusCode != http.StatusOK {
			return model.Release{}, fmt.Errorf("invalid status code: %d, url: %s", statusCode, assetsURL)
		}

		uploadStatus = getResponse.UploadStatus

		if uploadStatus == "readyToBePublished" {
			releaseDistinctID = getResponse.ReleaseDistinctID
		}
	}

	return api.GetReleaseOnAppByID(opts.App, releaseDistinctID)
}
