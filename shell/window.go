package shell

import (
	"encoding/json"
	"errors"
	"time"
)

// Window handle to an electron window.
type Window struct {
	WindowID int
	electron *Electron
}

// WindowOptions - creation options for creating a window with various properties.
type WindowOptions struct {
	Width            int    `json:"width,omitempty"`
	Height           int    `json:"height,omitempty"`
	X                int    `json:"x,omitempty"`
	Y                int    `json:"y,omitempty"`
	UseContentSize   bool   `json:"useContentSize"`
	Center           bool   `json:"center"`
	MinWidth         int    `json:"minWidth,omitempty"`
	MinHeight        int    `json:"minHeight,omitempty"`
	MaxWidth         int    `json:"maxWidth,omitempty"`
	MaxHeight        int    `json:"maxHeight,omitempty"`
	Resizable        bool   `json:"resizable"`
	AlwaysOnTop      bool   `json:"alwaysOnTop"`
	Fullscreen       bool   `json:"fullscreen"`
	SkipTaskbar      bool   `json:"skipTaskbar"`
	Kiosk            bool   `json:"kiosk"`
	Title            string `json:"title,omitempty"`
	Show             bool   `json:"show"`
	Frame            bool   `json:"frame"`
	AcceptFirstMouse bool   `json:"acceptFirstMouse"`
}

// CreateWindow - creates a window on the remote shell. Takes a WindowOptions which
// loosely maps to Electron's window options.
func (e *Electron) CreateWindow(options WindowOptions) (*Window, error) {

	// convert WindowOptions to JSON and send the command to the remote shell.
	jsonData, err := json.Marshal(options)
	if nil != err {
		return nil, err
	}

	e.Command("window_create", jsonData)

	// wait for a response from the remote shell for 5 seconds or timeout.
	var response []byte
	select {
	case r := <-windowCreationResponses:
		response = r
	case <-time.After(time.Second * 5):
		return nil, errors.New("WindowCreate timed out.")
	}

	// turn the JSON data into a WindowID.
	window := Window{electron: e}
	err = json.Unmarshal(response, &window)
	if nil != err {
		return nil, err
	}
	return &window, nil
}

// Here are the channels for synchronous window ops
var windowChannelsInitialized bool
var windowCreationResponses chan []byte
var windowLoadResponses chan []byte

// Here are the callbacks that pipe to channels for synchronous window ops.
func windowCreationCallback(data []byte) {
	windowCreationResponses <- data
}

func windowLoadCompletionCallback(data []byte) {
	windowLoadResponses <- data
}

// InitializeWindowCallbacks - sets up callbacks and channels for synchronous window operations.
func InitializeWindowCallbacks(electron *Electron) {
	if !windowChannelsInitialized {
		windowChannelsInitialized = true

		windowCreationResponses = make(chan []byte)
		windowLoadResponses = make(chan []byte)
	}

	electron.Listen("window_create_response", windowCreationCallback)
	electron.Listen("window_load_complete", windowLoadCompletionCallback)
}

type loadURLCommand struct {
	WindowID int
	URL      string
}

// LoadURL - Commands the window to load a URL.
func (w *Window) LoadURL(location string) error {
	// create the JSON for the loadURL command.
	c := loadURLCommand{WindowID: w.WindowID, URL: location}
	jsonData, err := json.Marshal(c)
	if nil != err {
		return err
	}

	w.electron.Command("window_load_url", jsonData)

	select {
	case <-windowLoadResponses:
		return nil
	case <-time.After(time.Second * 30):
		return errors.New("LoadURL timed out.")
	}
}
