package ElectronControl

import (
	"os"
	"testing"
	"time"

	"github.com/rrohrer/go-electroncontrol/rpc"
)

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

	remote, err := rpc.Launch("electron", "c:/sandbox/electroncontrol/app/")
	if nil != err {
		t.Error(err)
	} else {
		t.Log("Successfully launched.")
	}
	remote.Command("window_create", []byte("{'width':800, 'height':800}"))
	<-time.After(time.Second * 1)
	remote.Close()
	<-time.After(time.Second * 1)
}

func TestBasicAPI(t *testing.T) {
	err := Initialize()
	if nil != err {
		t.Error(err)
	}

	SetCommandArguments("c:/sandbox/electroncontrol/app/")

	electron, err := New()
	if nil != err {
		t.Error(err)
	}

	defer electron.Close()
}
