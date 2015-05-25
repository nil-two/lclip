package main

import (
	"github.com/gonuts/commander"
)

var cmd_set = &commander.Command{
	UsageLine: "set LABEL",
	Short:     "set text to LABEL",
	Run: func(cmd *commander.Command, args []string) error {
		return nil
	},
}

func main() {
}
