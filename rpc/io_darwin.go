package rpc

import (
	"io"
	"os/exec"
)

// SetupRemoteIO - starts the threads that are responsible for handling remote IO.
func SetupRemoteIO(remoteIn io.WriteCloser, remoteOut io.ReadCloser, cmd *exec.Cmd) (*Remote, error) {
	// set up the communication channel.
	output := make(chan []byte)

	// set up the remote data.
	remote := Remote{
		outgoing:  output,
		process:   cmd,
		listeners: make(map[string]RemoteListener)}

	// set up the reader and writer.
	go remoteWriter(output, remoteIn)
	go remoteReader(remote.handler, remoteOut)

	return &remote, nil
}
