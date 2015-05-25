package lclip

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func createStorageFile(path string) error {
	empty := []byte("{}\n")
	return ioutil.WriteFile(path, empty, 0644)
}

func DefaultPath() (string, error) {
	h, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(h, ".lclip.db"), nil
}

type Clipboard struct {
	path    string
	storage map[string][]byte
}

func NewClipboard(path string) (*Clipboard, error) {
	if !exists(path) {
		if err := createStorageFile(path); err != nil {
			return nil, err
		}
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c := &Clipboard{path: path}
	if err = json.NewDecoder(f).Decode(&c.storage); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Clipboard) Get(label string) ([]byte, error) {
	return c.storage[label], nil
}

func (c *Clipboard) Set(label string, data []byte) error {
	c.storage[label] = data
	return nil
}

func (c *Clipboard) Labels() []string {
	a := make([]string, 0, len(c.storage))
	for label, _ := range c.storage {
		a = append(a, label)
	}
	return a
}

func (c *Clipboard) Close() error {
	f, err := os.Create(c.path)
	if err != nil {
		return err
	}
	return json.NewEncoder(f).Encode(&c.storage)
}
