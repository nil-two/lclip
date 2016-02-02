package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/yuya-takeyama/argf"
)

var (
	name    = "lclip"
	version = "0.6.0"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, `
usage: %[1]s <operation> [...]
operations:
  %[1]s {-h --help}                      # show this help message
  %[1]s {-v --version}                   # print the version
  %[1]s {-l --labels}                    # list labels
  %[1]s {-g --get}    <label>            # paste text from label
  %[1]s {-s --set}    <label> [file(s)]  # copy text to label
  %[1]s {-d --delete} <label(s)>         # delete label(s)
`[1:], name)
}

func printVersion() {
	fmt.Fprintln(os.Stderr, version)
}

func cmdGet(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no specify label")
	}
	label := args[0]

	c, err := NewClipboardWithDefaultPath()
	if err != nil {
		return err
	}
	defer c.Close()

	dst := c.Get(label)
	if _, err = os.Stdout.Write(dst); err != nil {
		return err
	}
	return nil
}

func cmdSet(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no specify label")
	}
	label := args[0]

	c, err := NewClipboardWithDefaultPath()
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

	c.Set(label, src)
	return nil
}

func cmdLabels() error {
	c, err := NewClipboardWithDefaultPath()
	if err != nil {
		return err
	}
	defer c.Close()

	labels := c.Labels()
	sort.Strings(labels)
	for _, label := range labels {
		fmt.Println(label)
	}
	return nil
}

func cmdDelete(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no specify label")
	}
	labels := args

	c, err := NewClipboardWithDefaultPath()
	if err != nil {
		return err
	}
	defer c.Close()

	for _, label := range labels {
		c.Delete(label)
	}
	return nil
}

func _main() int {
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
	flag.Usage = printUsage
	flag.Parse()

	switch {
	case flag.NFlag() == 0:
		printUsage()
		return 2
	case flag.NFlag() > 1:
		fmt.Fprintf(os.Stderr, "%s: conflicting command specified\n", name)
		return 2
	}

	switch {
	case isHelp:
		printUsage()
		return 0
	case isVersion:
		printVersion()
		return 0
	}

	var err error
	switch {
	case isGet:
		err = cmdGet(flag.Args())
	case isSet:
		err = cmdSet(flag.Args())
	case isLabels:
		err = cmdLabels()
	case isDelete:
		err = cmdDelete(flag.Args())
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", name, err)
		return 1
	}
	return 0
}

func main() {
	e := _main()
	os.Exit(e)
}
