package appcenter

import (
	"strings"

	"github.com/bitrise-io/appcenter/client"
	"github.com/bitrise-io/appcenter/model"
)

// ReleaseAPI ...
type ReleaseAPI struct {
	api     client.API
	release model.Release
	opts    model.ReleaseOptions
}

// AddGroup ...
func (r ReleaseAPI) AddGroup(g model.Group) error {
	return r.api.AddReleaseToGroup(g, r.release.ID, r.opts)
}

// AddGroupsToRelease ...
func (r ReleaseAPI) AddGroupsToRelease(groupNames []string) error {
	if len(groupNames) > 0 {
		for _, groupName := range groupNames {
			if len(strings.TrimSpace(groupName)) == 0 {
				continue
			}
			group, err := r.api.GetGroupByName(groupName, r.opts.App)
			if err != nil {
				return err
			}

			err = r.AddGroup(group)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// AddStore ...
func (r ReleaseAPI) AddStore(s model.Store) error {
	return r.api.AddReleaseToStore(s, r.release.ID, r.opts)
}

// AddTester ...
func (r ReleaseAPI) AddTester(email string, mandatoryUpdate, notifyTesters bool) error {
	return r.api.AddTesterToRelease(email, r.release.ID, r.opts)
}

// SetReleaseNote ...
func (r ReleaseAPI) SetReleaseNote(releaseNote string) error {
	return r.api.SetReleaseNoteOnRelease(releaseNote, r.release.ID, r.opts)
}

// UploadSymbol - build and version is required for Android and optional for iOS
func (r ReleaseAPI) UploadSymbol(filePath string) error {
	return r.api.UploadSymbolToRelease(filePath, r.release, r.opts)
}
