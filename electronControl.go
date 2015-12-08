package ElectronControl

import "github.com/rrohrer/go-electroncontrol/shell"

// Shell - the containing data structure that wraps an instance of electron shell.
type Shell shell.Instance

// electronCommandArguments - sets the arguments to pass into the electron command.
var electronCommandArguments []string

// electronExeName - the final step of the command to execute to start electron.
var electronExeName = "electron.exe"

// electronPath - system path that electron resides in.
var electronPath string

// SetPath - sets the path that the system should look for electron in.
func SetPath(path string) {
	electronPath = path
}

// SetExecutableName - sets the name of the executeable that is going to get launched.
func SetExecutableName(name string) {
	electronExeName = name
}

// SetCommandArguments - sets the command arguments that get passed into electron on launch.
func SetCommandArguments(args ...string) {
	electronCommandArguments = args
}

// NewShell - launches electronExeName and opens a connection.
func NewShell() Shell {
	return Shell{}
}
