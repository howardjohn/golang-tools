package imports

import (
	"os"
	"testing"

	"github.com/howardjohn/golang-tools/internal/testenv"
)

func TestMain(m *testing.M) {
	testenv.ExitIfSmallMachine()
	os.Exit(m.Run())
}
