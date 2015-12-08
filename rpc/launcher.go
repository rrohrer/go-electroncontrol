package rpc

import "os/exec"

// Launch - launches electron from executableLocation with the arguments
// provided.
func Launch(executableLocation string, arguments ...string) (*Remote, error) {

	// build the command to execute launching electron.
	cmd := exec.Command(executableLocation, arguments...)

	// start Electron Shell.
	err := cmd.Start()
	if nil != err {
		return nil, err
	}

	// hook into stdin and stdout
	remoteStdin, err := cmd.StdinPipe()
	if nil != err {
		return nil, err
	}

	remoteStdout, err := cmd.StdoutPipe()
	if nil != err {
		return nil, err
	}
	// returns an activated Remote struct that can send and recieve commands.
	remote, err := SetupRemoteIO(remoteStdin, remoteStdout, cmd)
	if nil != err {
		return nil, err
	}

	return remote, nil
}
