package commands

import (
	"os"

	"github.com/nitrous-io/tug/Godeps/_workspace/src/github.com/jwaldrip/odin/cli"
	"github.com/nitrous-io/tug/docker"
)

var cmdInit = cli.NewSubCommand("init", "Initialize tug.", runInit)

func init() {
	cmdInit.SetLongDescription(`
Initialize tug.

Examples:

  tug init
  `)
}

func runInit(c cli.Command) {
	cmd := docker.Pull("nitrousio/docker-forward")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	cmd = docker.Pull("ddollar/docker-gateway")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
