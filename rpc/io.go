package rpc

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os/exec"
	"sync"
)

// Remote - holds a connection
type Remote struct {
	outgoing  chan<- []byte
	process   *exec.Cmd
	listeners map[string]RemoteListener
	sync.RWMutex
	conn net.Conn
}

// RemoteListener - function signature that allows things to listen to commands
// sent from the remote process.
type RemoteListener func([]byte)

// wrapper for the JSON commands that are sent over stdin/stdout.
type command struct {
	CommandID  string
	ComandBody string
}

// Command - queues a command to be sent to the remote.
func (r *Remote) Command(commandID string, commandBody []byte) error {
	fmt.Println(r)
	// take the whole command as a JSON string.
	data, err := json.Marshal(command{commandID, string(commandBody)})
	if nil != err {
		return err
	}

	// base64 encode the string.
	base64Data := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(base64Data, data)

	// send the string to the Remote.
	r.outgoing <- base64Data
	return nil
}

// Listen - registers a callback to occur on a specified remote commad.
func (r *Remote) Listen(commandID string, listener RemoteListener) {
	// lock for systems that touch listeners map.
	// lock for reading and writing because this function writes.
	r.Lock()
	defer r.Unlock()

	r.listeners[commandID] = listener
}

// Handler - handles a callback from the remote process.
func (r *Remote) Handler(remoteData []byte) {
	// lock for functions touch the listeners map.
	// only lock for reading because this function only reads.
	r.RLock()
	defer r.RUnlock()

	// base64 decode the message
	data := make([]byte, base64.StdEncoding.DecodedLen(len(remoteData)))
	_, err := base64.StdEncoding.Decode(data, remoteData)
	if nil != err {
		return
	}

	// unpack the command + command body into the containing struct.
	cmd := command{}
	err = json.Unmarshal(data, cmd)
	if nil != err {
		return
	}

	// call the callback if it exists.
	if key, ok := r.listeners[cmd.CommandID]; ok {
		key([]byte(cmd.ComandBody))
	}
}

// RemoteWriter - writes to the stdIN of the process that was passed in.
func RemoteWriter(output <-chan []byte, remoteIn io.WriteCloser) {
	for {
		data := <-output

		_, err := remoteIn.Write(data)
		if nil != err {
			return
		}

		_, err = remoteIn.Write([]byte("\n"))
		if nil != err {
			return
		}
	}
}

// RemoteReader - reads from stdOUT of the process that was passed in.
func RemoteReader(callback func([]byte), remoteOut io.ReadCloser) {

	// make a new reader
	reader := bufio.NewReader(remoteOut)

	for {
		// read a line from the stdout of the remote program.
		line, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			return
		}

		fmt.Println(line)

		// asyncronously send it to the client.
		go callback(line)
	}
}
