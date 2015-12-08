package shell

import "github.com/rrohrer/go-electroncontrol/rpc"

// Electron - the containing data structure that wraps an instance of electron shell.
type Electron struct {
	iD     uint
	remote *rpc.Remote
}

// ID'd used for electron instances.
var currentID uint

func getNextID() (id uint) {
	id = currentID
	currentID++
	return
}

// New - returns a new Electron instance.
func New(electronLocation string, args ...string) (*Electron, error) {
	remote, err := rpc.Launch(electronLocation, args...)
	if nil != err {
		return nil, err
	}
	return &Electron{getNextID(), remote}, nil
}

// Close - shutdown an electron instance.
func (e *Electron) Close() {
	e.remote.Close()
}
