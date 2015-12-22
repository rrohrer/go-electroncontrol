package shell

import (
	"errors"
	"sync"

	"github.com/rrohrer/go-electroncontrol/rpc"
)

// Electron - the containing data structure that wraps an instance of electron shell.
type Electron struct {
	iD            uint
	activeWindows map[int]*Window
	remote        *rpc.Remote
	sync.RWMutex
}

// ID'd used for electron instances.
var currentID uint

func getNextID() (id uint) {
	id = currentID
	currentID++
	return
}

// New - returns a new Electron instance.
func New(electronLocation, workingDir string, args ...string) (*Electron, error) {
	// launch the remote instance of electron.
	remote, err := rpc.Launch(electronLocation, workingDir, args...)
	if nil != err {
		return nil, err
	}

	// create the electron instance.
	electron := &Electron{iD: getNextID(), remote: remote, activeWindows: make(map[int]*Window)}

	// setup the responders for callbacks to windows.
	InitializeWindowCallbacks(electron)
	return electron, nil
}

// Close - shutdown an electron instance.
func (e *Electron) Close() {
	e.remote.Close()
}

// Command - Sends a remote command to the Electron instance.
func (e *Electron) Command(commandID string, commandBody []byte) error {
	if nil == e.remote {
		return errors.New("Called Command on a nil remote.")
	}
	return e.remote.Command(commandID, commandBody)
}

// Listen - registers a callback to occur on a specified remote command.
func (e *Electron) Listen(commandID string, listener rpc.RemoteListener) error {
	if nil == e.remote {
		return errors.New("Called Listen on nil remote.")
	}

	e.remote.Listen(commandID, listener)
	return nil
}
