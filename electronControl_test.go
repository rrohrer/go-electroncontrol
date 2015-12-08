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
	remote, err := rpc.Launch("electron", "c:/sandbox/electron/electroncontrol/app/")
	if nil != err {
		t.Error(err)
	} else {
		t.Log("Successfully launched.")
	}
	<-time.After(time.Second * 5)
	remote.Command("window_create", []byte("{'width':800, 'height':800}"))
	<-time.After(time.Second * 5)
}
