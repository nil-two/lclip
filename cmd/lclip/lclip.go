package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/kusabashira/lclip"
	"github.com/yuya-takeyama/argf"
)

func usage() {
	os.Stderr.WriteString(`
cli interface for labeled clipboard.

Commands:
	lclip [ -g | --get ] LABEL           # Paste text from LABEL
	lclip [ -s | --set ] LABEL [FILE]... # Copy text to LABEL
	lclip [ -l | --labels ]              # List labels
	lclip [ -d | --delete ] [LABEL]...   # Delete LABEL(s)

	lclip [ -h | --help ]                # Show this help message
	lclip [ -v | --version ]             # Print the version
`[1:])
}

func version() {
	os.Stderr.WriteString(`
v0.2.0
`[1:])
}

func cmd_get(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no specify LABEL")
	}
	label := args[0]

	c, err := lclip.NewClipboardWithDefaultPath()
	if err != nil {
		return err
	}
	defer c.Close()

	dst, err := c.Get(label)
	if err != nil {
		return err
	}
	if _, err = os.Stdout.Write(dst); err != nil {
		return err
	}
	return nil
}

func cmd_set(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no specify LABEL")
	}
	label := args[0]

	c, err := lclip.NewClipboardWithDefaultPath()
	if err != nil {
		return err
	}
	defer c.Close()

	r, err := argf.From(args[1:])
	if err != nil {
		return err
	}
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return c.Set(label, src)
}

func cmd_labels() error {
	c, err := lclip.NewClipboardWithDefaultPath()
	if err != nil {
		return err
	}
	defer c.Close()

	labels, err := c.Labels()
	if err != nil {
		return err
	}
	sort.Strings(labels)
	for _, label := range labels {
		fmt.Println(label)
	}
	return nil
}

func cmd_delete(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no specify LABEL")
	}
	labels := args

	c, err := lclip.NewClipboardWithDefaultPath()
	if err != nil {
		return err
	}
	defer c.Close()

	for _, label := range labels {
		if err = c.Delete(label); err != nil {
			return err
		}
	}
	return nil
}

func _main() error {
	var isHelp, isVersion bool
	flag.BoolVar(&isHelp, "h", false, "")
	flag.BoolVar(&isHelp, "help", false, "")
	flag.BoolVar(&isVersion, "v", false, "")
	flag.BoolVar(&isVersion, "version", false, "")

	var isGet, isSet, isLabels, isDelete bool
	flag.BoolVar(&isGet, "g", false, "")
	flag.BoolVar(&isGet, "get", false, "")
	flag.BoolVar(&isSet, "s", false, "")
	flag.BoolVar(&isSet, "set", false, "")
	flag.BoolVar(&isLabels, "l", false, "")
	flag.BoolVar(&isLabels, "labels", false, "")
	flag.BoolVar(&isDelete, "d", false, "")
	flag.BoolVar(&isDelete, "delete", false, "")
	flag.Usage = usage
	flag.Parse()

	if flag.NFlag() > 1 {
		return fmt.Errorf("cannot specify more than one command")
	}
	switch {
	case isHelp:
		usage()
		return nil
	case isVersion:
		version()
		return nil
	case isGet:
		return cmd_get(flag.Args())
	case isSet:
		return cmd_set(flag.Args())
	case isLabels:
		return cmd_labels()
	case isDelete:
		return cmd_delete(flag.Args())
	}
	usage()
	return nil
}

func main() {
	if err := _main(); err != nil {
		fmt.Fprintln(os.Stderr, "lclip:", err)
		os.Exit(1)
	}
}
