package ElectronControl

import (
	"os"
	"testing"
	"time"

	"github.com/rrohrer/go-electroncontrol/rpc"
	"github.com/rrohrer/go-electroncontrol/shell"
)

var electronAppLocation = "c:/sandbox/electron/electroncontrol/app/"

func TestMain(m *testing.M) {
	// do test setup here.

	// execute tests and quit.
	os.Exit(m.Run())
}

func TestBasicLaunch(t *testing.T) {
	err := rpc.Initialize()
	if nil != err {
		t.Error(err)
	}
	defer rpc.Shutdown()

	remote, err := rpc.Launch("electron", electronAppLocation)
	if nil != err {
		t.Error(err)
		return
	}
	t.Log("Successfully launched.")

	remote.Command("window_create", []byte("{\"width\":\"1100\", \"height\":\"1100\"}"))
	<-time.After(time.Second * 1)
	remote.Close()
	<-time.After(time.Second * 1)
}

func TestBasicAPI(t *testing.T) {
	err := Initialize()
	if nil != err {
		t.Error(err)
	}

	SetCommandArguments(electronAppLocation)

	electron, err := New()
	if nil != err {
		t.Error(err)
		return
	}

	electron.CreateWindow(shell.WindowOptions{Width: 700, Height: 700, Frame: false, Show: true})
	<-time.After(time.Second * 1)

	defer electron.Close()
}
