package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nitrous-io/tug/Godeps/_workspace/src/github.com/jwaldrip/odin/cli"
	"github.com/nitrous-io/tug/docker"
)

var Shell = cli.NewSubCommand("shell", "Open a shell in a running container.", runShell)

func init() {
	Shell.SetLongDescription(`
Open a shell on a running container.

Examples:

  tug shell web
  `)
	Shell.DefineParams("container")
}

func runShell(c cli.Command) {
	wd, _ := os.Getwd()
	abs, _ := filepath.Abs(wd)
	tag := fmt.Sprintf("%s.%s", filepath.Base(abs), c.Param("container").String())
	cmd := docker.ExecInteractive(tag, "bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
