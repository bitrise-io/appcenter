package appcenter

import (
	"fmt"

	"github.com/bitrise-io/appcenter/client"
	"github.com/bitrise-io/appcenter/commander"
	"github.com/bitrise-io/appcenter/model"
)

// AppAPI ...
type AppAPI struct {
	API             client.API
	CommandExecutor commander.CommandExecutor
	ReleaseOptions  model.ReleaseOptions
	CLIParams       model.CLIParams
}

// CreateApplicationAPI ...
func CreateApplicationAPI(api client.API, releaseOptions model.ReleaseOptions, cliParams model.CLIParams) AppAPI {
	return AppAPI{
		API:             api,
		ReleaseOptions:  releaseOptions,
		CommandExecutor: commander.CommandExecutor{},
		CLIParams:       cliParams,
	}
}

// NewRelease ...
// Uploads the artifact with the AppCenter CLI does the following:
// 1) Uploads the artifact and sets the first given group as "default" group.
// 2) Fetches the releases and gets the latest because it is the recent uploaded release.
func (a AppAPI) NewRelease() (model.Release, error) {
	releaseID, err := a.API.CreateRelease(a.ReleaseOptions)
	if err != nil {
		return model.Release{},
			fmt.Errorf("failed to create new release on app: %s, owner: %s, %v",
				a.ReleaseOptions.App.AppName,
				a.ReleaseOptions.App.Owner,
				err)
	}

	return a.API.GetAppReleaseDetails(a.ReleaseOptions.App, releaseID)
}

func (a AppAPI) createCLICommandArgs(opts model.ReleaseOptions, cliParams model.CLIParams) []string {
	appName := opts.App.Owner + "/" + opts.App.AppName
	commandArgs := []string{
		"distribute",
		"release",
		"-a", appName,
		"-f", opts.FilePath,
		"-g", opts.GroupNames[0],
		"--token", cliParams.APIToken,
	}

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
func (a AppAPI) Groups(name string) (model.Group, error) {
	return a.API.GetGroupByName(name, a.ReleaseOptions.App)
}

// Stores ...
func (a AppAPI) Stores(name string) (model.Store, error) {
	return a.API.GetStore(name, a.ReleaseOptions.App)
}
