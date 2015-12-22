package ElectronControl

import (
	"path/filepath"

	"github.com/rrohrer/go-electroncontrol/rpc"
	"github.com/rrohrer/go-electroncontrol/shell"
)

// electronCommandArguments - sets the arguments to pass into the electron command.
var electronCommandArguments []string

// electronExeName - the final step of the command to execute to start electron.
var electronExeName = "electron.exe"

// electronPath - system path that electron resides in.
var electronPath string

// electronWorkingDir - workingDir that electron comes from.
var electronWorkingDir string

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

// SetWorkingDir - set the working dir of electron.exe
func SetWorkingDir(dir string) {
	electronWorkingDir = dir
}

// New - launches electronExeName and opens a connection.
func New() (*shell.Electron, error) {

	return shell.New(filepath.Join(electronPath, electronExeName), electronWorkingDir, electronCommandArguments...)
}

// Initialize - Sets up ElectronControl
func Initialize() error {
	return rpc.Initialize()
}

// Shutdown - closes down ElectronControl.
func Shutdown() {
	rpc.Shutdown()
}
