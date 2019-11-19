package appcenter

import (
	"fmt"
	"net/http"
)

// ReleaseOptions ...
type ReleaseOptions struct {
	BuildVersion string `json:"build_version,omitempty"`
	BuildNumber  string `json:"build_number,omitempty"`
	ReleaseID    int    `json:"release_id,omitempty"`
}

// Release ...
type Release struct {
	app                App
	ID                 int    `json:"id,omitempty"`
	Version            string `json:"version,omitempty"`
	Origin             string `json:"origin,omitempty"`
	ShortVersion       string `json:"short_version,omitempty"`
	Enabled            bool   `json:"enabled,omitempty"`
	UploadedAt         string `json:"uploaded_at,omitempty"`
	DestinationType    string `json:"destination_type,omitempty"`
	DistributionGroups []struct {
		ID       string `json:"id,omitempty"`
		Name     string `json:"name,omitempty"`
		IsLatest bool   `json:"is_latest,omitempty"`
	} `json:"distribution_groups,omitempty"`
	DistributionStores []struct {
		ID               string `json:"id,omitempty"`
		Name             string `json:"name,omitempty"`
		Type             string `json:"type,omitempty"`
		PublishingStatus string `json:"publishing_status,omitempty"`
		IsLatest         bool   `json:"is_latest,omitempty"`
	} `json:"distribution_stores,omitempty"`
	Destinations []struct {
		ID               string `json:"id,omitempty"`
		Name             string `json:"name,omitempty"`
		IsLatest         bool   `json:"is_latest,omitempty"`
		Type             string `json:"type,omitempty"`
		PublishingStatus string `json:"publishing_status,omitempty"`
		DestinationType  string `json:"destination_type,omitempty"`
		DisplayName      string `json:"display_name,omitempty"`
	} `json:"destinations,omitempty"`
	Build struct {
		BranchName    string `json:"branch_name,omitempty"`
		CommitHash    string `json:"commit_hash,omitempty"`
		CommitMessage string `json:"commit_message,omitempty"`
	} `json:"build,omitempty"`
	IsExternalBuild bool `json:"is_external_build,omitempty"`
}

// SetGroup ...
func (r Release) SetGroup(g Group, mandatoryUpdate, notifyTesters bool) error {
	var (
		postURL     = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases/%d/groups", baseURL, r.app.owner, r.app.name, r.ID)
		postRequest = struct {
			ID              string `json:"id"`
			MandatoryUpdate bool   `json:"mandatory_update"`
			NotifyTesters   bool   `json:"notify_testers"`
		}{
			ID:              g.ID,
			MandatoryUpdate: mandatoryUpdate,
			NotifyTesters:   notifyTesters,
		}
	)

	statusCode, err := r.app.client.jsonRequest(http.MethodPost, postURL, postRequest, nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusCreated {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, postURL)
	}

	return nil
}

// SetStore ...
func (r Release) SetStore(s Store) error {
	var (
		postURL     = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases/%d/stores", baseURL, r.app.owner, r.app.name, r.ID)
		postRequest = struct {
			ID string `json:"id"`
		}{
			ID: s.ID,
		}
	)

	statusCode, err := r.app.client.jsonRequest(http.MethodPost, postURL, postRequest, nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusCreated {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, postURL)
	}

	return nil
}

// SetReleaseNote ...
func (r Release) SetReleaseNote(releaseNote string) error {
	var (
		putURL     = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases/%d", baseURL, r.app.owner, r.app.name, r.ID)
		putRequest = struct {
			ReleaseNotes string `json:"release_notes"`
		}{
			ReleaseNotes: releaseNote,
		}
	)

	statusCode, err := r.app.client.jsonRequest(http.MethodPut, putURL, putRequest, nil)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d, url: %s", statusCode, putURL)
	}

	return nil
}
