package lclip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/mitchellh/go-homedir"
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

func TestDefaultPath(t *testing.T) {
	h, err := homedir.Dir()
	if err != nil {
		t.Fatal(err)
	}
	expect := filepath.Join(h, ".lclip.json")
	actual, err := DefaultPath()
	if err != nil {
		t.Errorf("DefaultPath returns %q; want nil", err)
	}
	if actual != expect {
		t.Errorf("DefaultPath = %q; want %q",
			actual, expect)
	}
}

func TestCreateStorageFileIfNotExists(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "test")
	if err != nil {
		t.Fatal(err)
	}
	if err = f.Close(); err != nil {
		t.Fatal(err)
	}

	path := f.Name()
	if err = os.Remove(path); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)

	c, err := NewClipboard(path)
	if err != nil {
		t.Errorf("NewClipboard returns %q; want nil", err)
	}
	defer c.Close()

	expect := []byte("{}\n")
	actual, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("got %q; want %q",
			actual, expect)
	}
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
	defer c.Close()
	for _, test := range indexTestsGetText {
		expect := test.Dst
		actual := c.Get(test.Src)
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
	{Label: "abc", Data: "def"},
	{Label: "", Data: ""},
}

func TestSetText(t *testing.T) {
	if err := ioutil.WriteFile(tempPath, []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}

	c, err := NewClipboard(tempPath)
	if err != nil {
		t.Errorf("NewClipboard returns %q; want nil", err)
	}
	defer c.Close()
	for _, test := range indexTestsSetText {
		c.Set(test.Label, test.Data)
		expect := test.Data
		actual := c.Get(test.Label)
		if actual != expect {
			t.Errorf("after Set(%q, %q), Get(%q) = %q; want %q",
				test.Label, test.Data,
				test.Label, actual, expect)
		}
	}
}

var indexTestsLabels = [][]string{
	{"foo", "bar", "baz"},
	{"hoge", "piyo", "fuga"},
}

func TestListLabels(t *testing.T) {
	empty := []byte(`{}`)
	for _, labels := range indexTestsLabels {
		if err := ioutil.WriteFile(tempPath, empty, 0644); err != nil {
			t.Fatal(err)
		}

		c, err := NewClipboard(tempPath)
		if err != nil {
			t.Errorf("NewClipboard returns %q; want nil", err)
		}
		for _, label := range labels {
			c.Set(label, "")
		}

		expect := append(make([]string, 0, len(labels)), labels...)
		actual := c.Labels()
		sort.Strings(expect)
		sort.Strings(actual)
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("got %q; want %q", actual, expect)
		}
		if err := c.Close(); err != nil {
			t.Error("Close returns %q; want nil", err)
		}
	}
}

func TestSaveText(t *testing.T) {
	if err := ioutil.WriteFile(tempPath, []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}

	k, v := "key", "value"
	{
		c, err := NewClipboard(tempPath)
		if err != nil {
			t.Errorf("NewClipboard returns %q; want nil", err)
		}
		c.Set(k, v)
		if err := c.Close(); err != nil {
			t.Error("Close returns %q; want nil", err)
		}
	}
	{
		c, err := NewClipboard(tempPath)
		if err != nil {
			t.Errorf("NewClipboard returns %q; want nil", err)
		}
		defer c.Close()
		expect := v
		actual := c.Get(k)
		if actual != expect {
			t.Errorf("Get(%q) = %q; want %q",
				k, actual, expect)
		}
	}
}
