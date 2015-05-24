package lclip

type Clipboard struct {
}

func NewClipboard(path string) (*Clipboard, error) {
	return &Clipboard{}, nil
}

func (c *Clipboard) Get(label string) string {
	return ""
}
