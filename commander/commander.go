package commander

import (
	"fmt"

	"github.com/bitrise-io/go-utils/command"
)

// CommandExecutor ...
type CommandExecutor struct{}

// ExecuteCommand ...
func (c CommandExecutor) ExecuteCommand(stringCommand string, args ...string) (string, error) {
	cmd := command.New(stringCommand, args...)
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s failed: %s", cmd.PrintableCommandArgs(), out)
	}
	return out, nil
}
