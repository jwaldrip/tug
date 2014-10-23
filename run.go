package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var cmdRun = &Command{
	Run:   runRun,
	Usage: "run CONTAINER COMMAND",
	Short: "Run a command on a running container",
	Long: `
Run a command on an running container.

Examples:

  tug run web rake db:create
`,
}

func runRun(c *Command, args []string) {
	if len(args) < 1 {
		fail("must specify a container name\n")
		return
	}
	if len(args) < 2 {
		fail("must specify a command\n")
		return
	}
	wd, _ := os.Getwd()
	abs, _ := filepath.Abs(wd)
	tag := fmt.Sprintf("%s.%s", filepath.Base(abs), args[0])
	cmd := DockerExecInteractive(tag, args[1:]...)
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
