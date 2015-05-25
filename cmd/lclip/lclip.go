package main

import (
	"fmt"
	"io/ioutil"

	"github.com/gonuts/commander"
	"github.com/kusabashira/lclip"
	"github.com/yuya-takeyama/argf"
)

var cmd_set = &commander.Command{
	UsageLine: "set LABEL [FILE]...",
	Short:     "set text to LABEL",
	Run: func(cmd *commander.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("no specify LABEL")
		}
		label := args[0]

		path, err := lclip.DefaultPath()
		if err != nil {
			return err
		}
		c, err := lclip.NewClipboard(path)
		if err != nil {
			return err
		}

		r, err := argf.From(args[1:])
		if err != nil {
			return err
		}
		src, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}

		c.Set(label, src)
		return c.Close()
	},
}

func main() {
}
