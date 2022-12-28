package command

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// TODO use this env to understand test running and
	// use inMem fs for testing cli commands and flags
	os.Setenv("GO_ENV", "testing")
	os.Exit(m.Run())
}
