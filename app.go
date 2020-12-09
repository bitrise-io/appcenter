package appcenter

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bitrise-io/go-utils/command"
)

// App ...
type App struct {
	client      Client
	owner, name string
}

// ExecuteCommand ...
func ExecuteCommand(stringCommand string, args ...string) (string, error) {
	cmd := command.New(stringCommand, args...)
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s failed: %s", cmd.PrintableCommandArgs(), out)
	}
	return out, nil
}

// NewRelease ... for more info: https://docs.microsoft.com/en-us/appcenter/distribution/uploading#uploading-using-the-apis
func (a App) NewRelease(filePath string, opts ReleaseOptions) (Release, error) {
	//set the firt group of the app (the CLI not able to set multiple groups)
	appName := a.owner + "/" + a.name
	commandArgs := []string{"distribute", "release", "--app", appName, "-f", filePath, "-g", opts.GroupNames[0]}

	if len(opts.BuildNumber) != 0 {
		commandArgs = append(commandArgs, "--build-number")
		commandArgs = append(commandArgs, opts.BuildNumber)
	}

	if len(opts.BuildVersion) != 0 {
		commandArgs = append(commandArgs, "--build-version")
		commandArgs = append(commandArgs, opts.BuildVersion)
	}

	str, err := ExecuteCommand("appcenter", commandArgs...)
	if err != nil {
		return Release{}, fmt.Errorf("Failed to create AppCenter release: %s", str)
	}

	fmt.Println(str)

	//fetch releases and find the latest
	var (
		getURL      = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases", baseURL, a.owner, a.name)
		getResponse []Release
	)

	statusCode, err := a.client.jsonRequest(http.MethodGet, getURL, nil, &getResponse)
	if err != nil {
		return Release{}, err
	}

	if statusCode != http.StatusOK {
		return Release{}, fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, getURL, getResponse)
	}

	latestReleaseID := getResponse[0].ID

	var (
		releaseShowURL = fmt.Sprintf("%s/v0.1/apps/%s/%s/releases/%s", baseURL, a.owner, a.name, strconv.Itoa(latestReleaseID))
		release        Release
	)

	statusCode, err = a.client.jsonRequest(http.MethodGet, releaseShowURL, nil, &release)
	if err != nil {
		return Release{}, err
	}

	if statusCode != http.StatusOK {
		return Release{}, fmt.Errorf("invalid status code: %d, url: %s, body: %v", statusCode, getURL, getResponse)
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
