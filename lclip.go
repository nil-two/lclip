package main

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

func DefaultPath() (string, error) {
	h, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(h, ".lclip.json"), nil
}

type Clipboard struct {
	path    string
	storage map[string][]byte
}

func NewClipboard(path string) (*Clipboard, error) {
	if !exists(path) {
		err := ioutil.WriteFile(path, []byte(`{}`), 0644)
		if err != nil {
			return nil, err
		}
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c := &Clipboard{path: path}
	if err := json.NewDecoder(f).Decode(&c.storage); err != nil {
		return nil, err
	}
	return c, nil
}

func NewClipboardWithDefaultPath() (*Clipboard, error) {
	path, err := DefaultPath()
	if err != nil {
		return nil, err
	}
	return NewClipboard(path)
}

func (c *Clipboard) Get(label string) []byte {
	return c.storage[label]
}

func (c *Clipboard) Set(label string, data []byte) {
	c.storage[label] = data
}

func (c *Clipboard) Labels() []string {
	labels := make([]string, 0, len(c.storage))
	for label, _ := range c.storage {
		labels = append(labels, label)
	}
	return labels
}

func (c *Clipboard) Delete(label string) {
	delete(c.storage, label)
}

func (c *Clipboard) Close() error {
	f, err := os.Create(c.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(&c.storage)
}
