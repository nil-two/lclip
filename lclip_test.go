package lclip

import (
	"bytes"
	"encoding/json"
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

type GetTextTest struct {
	Src string
	Dst string
}

var indexTestsGetText = []GetTextTest{
	{Src: "foo", Dst: "bar"},
	{Src: "hoge", Dst: "piyo"},
}

func TestGetText(t *testing.T) {
	w := bytes.NewBuffer(make([]byte, 0))
	m := make(map[string]string)
	for _, test := range indexTestsGetText {
		m[test.Src] = test.Dst
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(tempPath, w.Bytes(), 0644); err != nil {
		t.Fatal(err)
	}

	c, err := NewClipboard(tempPath)
	if err != nil {
		t.Errorf("NewClipboard returns %q; want nil", err)
	}
	for _, test := range indexTestsGetText {
		actual := c.Get(test.Src)
		expect := test.Dst
		if actual != expect {
			t.Errorf("Get(%q) = %q; want %q",
				test.Src, actual, expect)
		}
	}
}

type SetTextTest struct {
	Label string
	Data  string
}

var indexTestsSetText = []SetTextTest{
	{Label: "a", Data: "aaa"},
}

func TestSetText(t *testing.T) {
	if err := ioutil.WriteFile(tempPath, []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}

	c, err := NewClipboard(tempPath)
	if err != nil {
		t.Errorf("NewClipboard returns %q; want nil", err)
	}
	for _, test := range indexTestsSetText {
		c.Set(test.Label, test.Data)
		actual := c.Get(test.Label)
		expect := test.Data
		if actual != expect {
			t.Errorf("after Set(%q, %q), Get(%q) = %q; want %q",
				test.Label, test.Data,
				test.Label, actual, expect)
		}
	}
}
