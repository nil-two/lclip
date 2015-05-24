package lclip

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

func DefaultPath() (string, error) {
	h, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(h, ".lclip.json"), nil
}

type Clipboard struct {
	storage map[string]string
}

func NewClipboard(path string) (*Clipboard, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c := &Clipboard{}
	if err = json.NewDecoder(f).Decode(&c.storage); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Clipboard) Get(label string) string {
	return c.storage[label]
}

func (c *Clipboard) Set(label, data string) {
	c.storage[label] = data
}
