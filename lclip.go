package lclip

import (
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/go-homedir"
	"github.com/naoina/genmai"
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

type Variable struct {
	Label string `db:"pk"`
	Data  []byte
}

type Clipboard struct {
	path string
	db   *genmai.DB
}

func NewClipboard(path string) (*Clipboard, error) {
	db, err := genmai.New(&genmai.SQLite3Dialect{}, path)
	if err != nil {
		return nil, err
	}
	if err = db.CreateTableIfNotExists(&Variable{}); err != nil {
		return nil, err
	}
	return &Clipboard{
		path: path,
		db:   db,
	}, nil
}

func NewClipboardWithDefaultPath() (*Clipboard, error) {
	path, err := DefaultPath()
	if err != nil {
		return nil, err
	}
	return NewClipboard(path)
}

func (c *Clipboard) Path() string {
	return c.path
}

func (c *Clipboard) searchVariable(label string) (*Variable, error) {
	res := make([]Variable, 0, 1)
	err := c.db.Select(&res, c.db.Where("label", "=", label))
	if err != nil {
		return nil, err
	}
	if len(res) < 1 {
		return nil, nil
	}
	return &res[0], nil
}

func (c *Clipboard) Get(label string) []byte {
	v, err := c.searchVariable(label)
	if err != nil {
		return nil
	}
	if v == nil {
		return []byte(``)
	}
	return v.Data
}

func (c *Clipboard) Set(label string, data []byte) {
	v, err := c.searchVariable(label)
	if err != nil {
		return
	}
	if v == nil {
		_, err = c.db.Insert(&Variable{Label: label, Data: data})
		return
	}
	_, err = c.db.Update(&Variable{Label: v.Label, Data: data})
	return
}

func (c *Clipboard) Labels() []string {
	res := make([]Variable, 0)
	if err := c.db.Select(&res); err != nil {
		return nil
	}
	labels := make([]string, len(res))
	for i := 0; i < len(res); i++ {
		labels[i] = res[i].Label
	}
	return labels
}

func (c *Clipboard) Delete(label string) error {
	v, err := c.searchVariable(label)
	if err != nil {
		return err
	}
	if v == nil {
		return nil
	}
	_, err = c.db.Delete(v)
	return nil
}

func (c *Clipboard) Close() error {
	return c.db.Close()
}
