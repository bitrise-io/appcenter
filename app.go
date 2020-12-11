package appcenter

import (
	"fmt"

	"github.com/bitrise-io/appcenter/commander"
	"github.com/bitrise-io/appcenter/model"
)

// AppAPI ...
type AppAPI struct {
	api             API
	commandExecutor commander.CommandExecutor
}

// NewRelease ...
// Uploads the artifact with the AppCenter CLI does the following:
// 1) Uploads the artifact and sets the first given group as "default" group.
// 2) Fetches the releases and gets the latest because it is the recent uploaded release.
// 3) Fetches the lastest release full data.
// 4) Sets the remaining groups on the release with the API. Because AppCenter CLI is not able to set the groups in one command.
func (a AppAPI) NewRelease(opts model.ReleaseOptions) (model.Release, error) {
	//upload the artifact with the AappCenter CLI
	commandArgs := a.createCLICommandArgs(opts)
	str, err := a.commandExecutor.ExecuteCommand("appcenter", commandArgs...)
	if err != nil {
		return model.Release{}, fmt.Errorf("Failed to create AppCenter release: %s", str)
	}

	fmt.Println(fmt.Sprintf("Command execution result: %s", str))

	//fetch releases and find the latest
	release, err := a.api.GetLatestReleases(opts.App)
	if err != nil {
		return model.Release{}, err
	}

	return release, nil
}

func (a AppAPI) createCLICommandArgs(opts model.ReleaseOptions) []string {
	appName := opts.App.Owner + "/" + opts.App.AppName
	commandArgs := []string{"distribute", "release", "-a", appName, "-f", opts.FilePath, "-g", opts.GroupNames[0]}

	if len(opts.BuildNumber) != 0 {
		commandArgs = append(commandArgs, "--build-number")
		commandArgs = append(commandArgs, opts.BuildNumber)
	}

	if len(opts.BuildVersion) != 0 {
		commandArgs = append(commandArgs, "--build-version")
		commandArgs = append(commandArgs, opts.BuildVersion)
	}

	return commandArgs
}

// Groups ...
func (a AppAPI) Groups(name string, app model.App) (model.Group, error) {
	return a.api.GetGroupByName(name, app)
}

// Stores ...
func (a AppAPI) Stores(name string, app model.App) (model.Store, error) {
	return a.api.GetStore(name, app)
}
