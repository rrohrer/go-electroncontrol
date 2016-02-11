package rpc

import (
	"errors"
	"io"
	"net"
	"os/exec"
	"time"

	"gopkg.in/natefinch/npipe.v2"
)

// pipeName - named pipe that allows a connection to Electron.
var pipeName = `\\.\pipe\ElectronControl`

// keep track of whether is initiazed or not.
var isInitialized bool

// listener that is always running while the RPC library isn't shutdown.
var remoteIOListener net.Listener

// channel of incoming connections so that it's easy to block until Electron is up
// and running.
var incomingConnections chan net.Conn

// SetupRemoteIO - Windows implementation. Cannot use STDIO because Atom blocks it.
func SetupRemoteIO(remoteIn io.WriteCloser, remoteOut io.ReadCloser, cmd *exec.Cmd) (*Remote, error) {
	// set up the communication channel.
	output := make(chan []byte)

	// wait for the net connection to happen from Electron.
	// time out if it doesn't connect after 5 seconds.
	var conn net.Conn
	select {
	case c := <-incomingConnections:
		conn = c
	case <-time.After(time.Second * 5):
		return nil, errors.New("Connection to client timed out.")
	}

	// set up the remote data.
	remote := Remote{
		outgoing:      output,
		process:       cmd,
		listeners:     make(map[string]RemoteListener),
		conn:          conn,
		shutdownWrite: make(chan bool)}

	// set up the remote reader and writer threads.
	go RemoteReader(remote.Handler, conn)
	go RemoteWriter(output, conn, remote.shutdownWrite)

	// start the read and write threads.
	return &remote, nil
}

// InitializeIO - if this hasn't been called before, start accepting connections on
// the pipe that pipeName refers to.
func InitializeIO() error {
	// check to see if this is the first time this function is called.
	if isInitialized {
		return nil
	}

	// create a listener on the named pipe.
	remoteL, err := npipe.Listen(pipeName)
	if nil != err {
		return err
	}
	remoteIOListener = remoteL

	// create the listener channel.
	incomingConnections = make(chan net.Conn)

	// listen until Closed.
	go acceptLoop()
	isInitialized = true
	return nil
}

// ShutdownIO - stops accepting connections.
func ShutdownIO() {
	isInitialized = false
	remoteIOListener.Close()
}

// acceptLoop - loops until remoteIOListener is closed.
func acceptLoop() {
	for {
		conn, err := remoteIOListener.Accept()
		if nil != err {
			return
		}

		incomingConnections <- conn
	}
}
