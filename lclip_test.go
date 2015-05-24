package lclip

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	e := m.Run()
	os.Exit(e)
}
