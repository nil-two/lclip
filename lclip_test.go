package lclip

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var tempPath string

func TestMain(m *testing.M) {
	f, err := ioutil.TempFile(os.TempDir(), "test")
	if err != nil {
		fmt.Fprintln(os.Stderr, "lclip_test:", err)
		return
	}
	f.Close()
	tempPath = f.Name()

	e := m.Run()
	defer os.Exit(e)

	os.Remove(tempPath)
}

func TestGetText(t *testing.T) {
	tempData := []byte(`
{"foo": "bar"}
`[1:])
	if err := ioutil.WriteFile(tempPath, tempData, 0644); err != nil {
		t.Fatal(err)
	}

	c, err := NewClipboard(tempPath)
	if err != nil {
		t.Errorf("NewClipboard returns %q; want nil", err)
	}
	src := "foo"
	actual := c.Get(src)
	expect := "bar"
	if actual != expect {
		t.Errorf("Get(%q) = %q; want %q",
			src, actual, expect)
	}
}
