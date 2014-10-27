package commands

import "github.com/nitrous-io/tug/Godeps/_workspace/src/github.com/jwaldrip/odin/cli"

// Commands for tug
var Commands = []*cli.SubCommand{
	cmdInit,
	cmdStart,
	cmdShell,
	cmdRun,
}
