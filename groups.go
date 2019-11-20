package appcenter

import (
	"fmt"
	"net/http"
)

// Group ...
type Group struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Origin      string `json:"origin"`
	IsPublic    bool   `json:"is_public"`
}

// AddRelease ...
func (g Group) AddRelease(r Release, mandatoryUpdate, notifyTesters bool) error {
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
