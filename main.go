package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nitrous-io/tug/Godeps/_workspace/src/github.com/jwaldrip/odin/cli"
	"github.com/nitrous-io/tug/command"
	"github.com/nitrous-io/tug/docker"
	"github.com/nitrous-io/tug/helpers"
)

var tug = cli.New(VERSION, "Docker development workflow", cli.ShowUsage)

func init() {
	tug.AddSubCommands(
		command.Start,
		command.Shell,
		command.Run,
		command.Build,
		command.Deploy,
	)
}

func handlePanic() {
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "error: an unhandled exception has occurred\n")
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func checkDocker() {
	done := make(chan error)
	go func() {
		done <- docker.Ps().Run()
	}()
	select {
	case err := <-done:
		if err != nil {
			helpers.Die(fmt.Errorf("docker unavailable"))
		}
	case <-time.After(2 * time.Second):
		helpers.Die(fmt.Errorf("docker unavailable"))
	}
}

func main() {
	defer handlePanic()
	checkDocker()
	tug.Start()
}
