package shell

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// WindowEventCallback - function signature for window callbacks.
type WindowEventCallback func([]byte)

// WindowOnCloseCallback - function called when window is closed.
type WindowOnCloseCallback func()

// Window handle to an electron window.
type Window struct {
	WindowID       int
	electron       *Electron
	listeners      map[string]WindowEventCallback
	closedCallback WindowOnCloseCallback
	sync.RWMutex
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
	BackgroundColor  string `json:"backgroundColor,omitempty"`
}

// windowIDCommand - used for the many messages that involve only a WindowID
type windowIDCommand struct {
	WindowID int
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
	window := Window{electron: e, listeners: make(map[string]WindowEventCallback)}
	err = json.Unmarshal(response, &window)
	if nil != err {
		return nil, err
	}

	// lock electron so that the active window map can be written to.
	e.Lock()
	defer e.Unlock()

	// store the window.
	e.activeWindows[window.WindowID] = &window
	return &window, nil
}

// Close - used to shutdown a window.
func (w *Window) Close() {
	// lock electron so that the active window map can be written to.
	w.electron.Lock()
	defer w.electron.Unlock()

	// remove this from the active window list.
	delete(w.electron.activeWindows, w.WindowID)

	wID := windowIDCommand{w.WindowID}
	jsonData, _ := json.Marshal(wID)

	w.electron.Command("window_close", jsonData)
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

type windowListenCallbackPartial struct {
	WindowID  int
	MessageID string
	Message   json.RawMessage
}

func (e *Electron) windowListenCallback(data []byte) {
	e.RLock()
	defer e.RUnlock()

	// pull the partial message out of data
	partialData := windowListenCallbackPartial{}
	err := json.Unmarshal(data, &partialData)
	if nil != err {
		return
	}

	// pass the data on to the callback.
	if key, ok := e.activeWindows[partialData.WindowID]; ok && nil != key {
		key.RLock()
		defer key.RUnlock()
		if key1, ok1 := key.listeners[partialData.MessageID]; ok1 && nil != key1 {
			go key1(partialData.Message)
		}
	}
}

func (e *Electron) windowClosedCallback(data []byte) {
	e.Lock()
	defer e.Unlock()

	wID := windowIDCommand{}
	err := json.Unmarshal(data, &wID)
	if nil != err {
		return
	}

	if key, ok := e.activeWindows[wID.WindowID]; ok && nil != key {
		if nil != key.closedCallback {
			go key.closedCallback()
		}

		delete(e.activeWindows, wID.WindowID)
	}
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
	electron.Listen("window_get_subscribed_message", electron.windowListenCallback)
	electron.Listen("window_closed", electron.windowClosedCallback)
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

// listenToMessageCommand - struct that describes the JSON that will be seent to
// subscribee to messages.
type listenToMessageCommand struct {
	WindowID  int
	MessageID string
}

// Listen - register's a callback on the WINDOW CONTENT (webpage displayed by Electron)
func (w *Window) Listen(messageID string, callback WindowEventCallback) error {
	w.Lock()
	defer w.Unlock()

	// create the listen command JSON string.
	listenCommand := listenToMessageCommand{w.WindowID, messageID}
	jsonCommand, err := json.Marshal(listenCommand)
	if nil != err {
		return err
	}

	// save the callback.
	w.listeners[messageID] = callback

	// subscribe
	w.electron.Command("window_subscribe_message", jsonCommand)

	return nil
}

// sendMessageCommand - filled to send a message to the WINDOW CONTENT.
type sendMessageCommand struct {
	WindowID  int
	MessageID string
	Message   string
}

// Message - send a message to the WINDOW CONTENT (webpage displayed by Electron).
func (w *Window) Message(messageID string, message []byte) error {
	// create the command.
	messageCommand := sendMessageCommand{w.WindowID, messageID, string(message)}
	jsonCommand, err := json.Marshal(messageCommand)
	if nil != err {
		return err
	}

	w.electron.Command("window_send_message", jsonCommand)
	return nil
}

// OpenDevTools - opens a developer console on the window in question.
func (w *Window) OpenDevTools() {
	commandData := windowIDCommand{w.WindowID}
	jsonCommandData, _ := json.Marshal(commandData)
	w.electron.Command("window_open_dev_tools", jsonCommandData)
}

// CloseDevTools - closes developer tools on the window in question.
func (w *Window) CloseDevTools() {
	commandData := windowIDCommand{w.WindowID}
	jsonCommandData, _ := json.Marshal(commandData)
	w.electron.Command("window_close_dev_tools", jsonCommandData)
}

// OnClosed - called when the window is closed, either by ELECTRON or the user.
func (w *Window) OnClosed(callback WindowOnCloseCallback) {
	w.closedCallback = callback
}
