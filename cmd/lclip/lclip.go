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

	lclip [ -h | --help ]                # Show this help message
	lclip [ -v | --version ]             # Print the version
`[1:])
}

func version() {
	os.Stderr.WriteString(`
v0.1.0
`[1:])
}

func cmd_get(args []string) error {
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
	return c.Close()
}

func cmd_set(args []string) error {
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
}

func cmd_labels() error {
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
}

func _main() error {
	var isHelp, isVersion bool
	flag.BoolVar(&isHelp, "h", false, "")
	flag.BoolVar(&isHelp, "help", false, "")
	flag.BoolVar(&isVersion, "v", false, "")
	flag.BoolVar(&isVersion, "version", false, "")

	var isGet, isSet, isLabels bool
	flag.BoolVar(&isGet, "g", false, "")
	flag.BoolVar(&isGet, "get", false, "")
	flag.BoolVar(&isSet, "s", false, "")
	flag.BoolVar(&isSet, "set", false, "")
	flag.BoolVar(&isLabels, "l", false, "")
	flag.BoolVar(&isLabels, "labels", false, "")
	flag.Usage = usage
	flag.Parse()
	switch {
	case isHelp:
		usage()
		return nil
	case isVersion:
		version()
		return nil
	}

	if (isGet && isSet) || (isSet && isLabels) || (isLabels && isGet) {
		return fmt.Errorf("cannot specify more than one command")
	}
	switch {
	case isGet:
		return cmd_get(flag.Args())
	case isSet:
		return cmd_set(flag.Args())
	case isLabels:
		return cmd_labels()
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
