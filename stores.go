package appcenter

import (
	"fmt"
	"net/http"
)

// Store ...
type Store struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Track         string `json:"track"`
	IntuneDetails struct {
		TargetAudience struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"target_audience"`
		AppCategory struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"app_category"`
	} `json:"intune_details"`
	ServiceConnectionID string `json:"service_connection_id"`
	CreatedBy           string `json:"created_by"`
}

// AddRelease ...
func (s Store) AddRelease(r Release) error {
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
