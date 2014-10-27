package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nitrous-io/tug/Godeps/_workspace/src/github.com/jwaldrip/odin/cli"
	"github.com/nitrous-io/tug/docker"
)

var Run = cli.NewSubCommand("run", "Run a command on an running container.", runRun)

func init() {
	Run.SetLongDescription(`
Run a command on an running container.

Examples:

  tug run web rake db:create
	`)
	Run.DefineParams("container", "command")
}

func runRun(c cli.Command) {
	wd, _ := os.Getwd()
	abs, _ := filepath.Abs(wd)
	tag := fmt.Sprintf("%s.%s", filepath.Base(abs), c.Param("container").String())
	args := append([]string{c.Param("command").String()}, c.Args().Strings()...)
	cmd := docker.ExecInteractive(tag, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func runSocket(name string) string {
	wd, _ := os.Getwd()
	socket := filepath.Join(wd, ".tug", fmt.Sprintf("%s.sock", name))
	os.MkdirAll(filepath.Dir(socket), 0700)
	return socket
}
