package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/gonuts/commander"
	"github.com/kusabashira/lclip"
	"github.com/yuya-takeyama/argf"
)

var cmd_get = &commander.Command{
	UsageLine: "get LABEL",
	Short:     "get text from LABEL",
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

		dst := c.Get(label)
		if _, err = os.Stdout.Write(dst); err != nil {
			return err
		}
		if _, err = os.Stdout.Write([]byte("\n")); err != nil {
			return err
		}
		return c.Close()
	},
}

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

var cmd_labels = &commander.Command{
	UsageLine: "labels",
	Short:     "list labels",
	Run: func(cmd *commander.Command, args []string) error {
		path, err := lclip.DefaultPath()
		if err != nil {
			return err
		}
		c, err := lclip.NewClipboard(path)
		if err != nil {
			return err
		}

		labels := c.Labels()
		sort.Strings(labels)
		for _, label := range labels {
			fmt.Println(label)
		}
		return c.Close()
	},
}

var command = &commander.Command{
	UsageLine: "lclip",
	Short:     "cli interface for labeled clipboard",
	Subcommands: []*commander.Command{
		cmd_get,
		cmd_set,
		cmd_labels,
	},
}

func main() {
	if err := command.Dispatch(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "lclip:", err)
		os.Exit(1)
	}
}
