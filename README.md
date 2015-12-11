# go-electroncontrol
Golang bindings for electroncontrol

## Usage

```go
err := Initialize()
defer Shutdown()

SetCommandArguments("c:/")

electron, err := New()

window, err := electron.CreateWindow(shell.WindowOptions{
    Width: 700,
    Height: 700,
    Frame: false,
    Show: true})
defer window.Close()

// load a webpage
err = window.LoadURL("http://google.com")

// send an event to the app
data, err := json.Marshal(myData)
window.Message("my_event", data)

// create a callback and listen to an event from the BrowserWindow
cb := func (data []byte) {json.UnMarshal(data, &myData)}
window.Listen("my_event_from_webview", cb)

```

## API
To use this you need an install of electron, as well as the [electroncontrol](https://github.com/rrohrer/electroncontrol) app.

Set up the Path and command line so that go-electroncontrol can control your shell.

I would suggest you use something like [go-bindata](https://github.com/jteeuwen/go-bindata/) to package electron WITH your GO app,
and then decompress on first startup.

### ElectronControl
**SetPath**: Sets the path of electron(.exe)

---

**SetExecutableName**: Set the name of electron. (in case you want to rename it)

---

**SetCommandArguments**: VA ARGS that are the command switches to pass into electron.

---

**New**: Loads an instance of electron. Returns a `*shell.Electron` that can be used.

---

**Initialize**: Set up ElectronControl and get it ready.

---

**Shutdown**: Shut down.

---

### Shell
**Electron::Command**: Sends a command to that instance of Electron.  Not intended
to be called directly.

---

**Electron::Listen**: Listens to a command from electron.  Not intended to be called
directly.

---

**Electron::CreateWindow**: Create a new window on the instance of electron.  Takes a `shell.WindowOptions` which (fits a golang case convention) but has the same names as used
in electron's `new BrowserWindow(options)`.

---

**Window::Close**: Closes an active window.

---

**Window::LoadURL**: Load a URL in the window. Blocks until the URL load is complete.

---

**Window::Listen**: Allows a callback to be sent to a BrowserWindow.  This, along with **Window::Message** allow direct communication with the webview.

---

**Window::OpenDevTools**: Opens the webkit developer tools pannel.

---

**Window::CloseDevTools**: Closes the webkit developer tools panel.

---

**Window::OnClosed**: Allows a callback to be registered for when a user closes a window.
