package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var cmdShell = &Command{
	Run:   runShell,
	Usage: "shell NAME",
	Short: "Open a shell on a running container",
	Long: `
Open a shell on a running container.

Examples:

  tug shell web
`,
}

func runShell(c *Command, args []string) {
	if len(args) < 1 {
		fail("must specify a container name\n")
		return
	}
	wd, _ := os.Getwd()
	abs, _ := filepath.Abs(wd)
	tag := fmt.Sprintf("%s.%s", filepath.Base(abs), args[0])
	cmd := DockerExecInteractive(tag, "bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
