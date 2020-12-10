package appcenter

import (
	"fmt"

	"github.com/bitrise-io/appcenter/commander"
	"github.com/bitrise-io/appcenter/model"
)

// App ...
type App struct {
	client          Client
	owner, name     string
	commandExecutor commander.CommandExecutor
}

// NewRelease ...
// Uploads the artifact with the AppCenter CLI does the following:
// 1) Uploads the artifact and sets the first given group as "default" group.
// 2) Fetches the releases and gets the latest because it is the recent uploaded release.
// 3) Fetches the lastest release full data.
// 4) Sets the remaining groups on the release with the API. Because AppCenter CLI is not able to set the groups in one command.
func (a App) NewRelease(opts ReleaseOptions) (Release, error) {
	//upload the artifact with the AappCenter CLI
	commandArgs := a.createCLICommandArgs(opts)
	str, err := a.commandExecutor.ExecuteCommand("appcenter", commandArgs...)
	if err != nil {
		return Release{}, fmt.Errorf("Failed to create AppCenter release: %s", str)
	}

	fmt.Println(fmt.Sprintf("Command execution result: %s", str))

	//fetch releases and find the latest
	api := API{Client: a.client}
	release, err := api.GetLatestReleases(a)
	if err != nil {
		return Release{}, err
	}

	release.app = a

	// set the groups on the app
	if len(opts.GroupNames) > 1 {
		for _, groupName := range opts.GroupNames[1:] {
			if len(groupName) == 0 {
				continue
			}
			group, err := a.Groups(groupName)
			if err != nil {
				return Release{}, err
			}

			release.AddGroup(group, opts.Mandatory, opts.NotifyTesters)
		}
	}

	return release, nil
}

func (a App) createCLICommandArgs(opts ReleaseOptions) []string {
	appName := a.owner + "/" + a.name
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
func (a App) Groups(name string) (model.Group, error) {
	api := API{Client: a.client}

	return api.GetGroupByName(name, a)
}

// Stores ...
func (a App) Stores(name string) (model.Store, error) {
	api := API{Client: a.client}

	return api.GetStore(name, a)
}
