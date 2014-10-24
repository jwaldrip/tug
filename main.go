package main

import (
	"fmt"
	"os"

	"github.com/nitrous-io/tug/Godeps/_workspace/src/github.com/mgutz/ansi"
)

var commands = []*Command{
	cmdInit,
	cmdStart,
	cmdShell,
	cmdRun,
	cmdVersion,
	cmdHelp,
}

type StringSet []string

func (ss *StringSet) Set(value string) error {
	*ss = append(*ss, value)
	return nil
}

func (ss *StringSet) String() string {
	return "[]"
}

func handlePanic() {
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "error: an unhandled exception has occurred\n")
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func debug(format string, a ...interface{}) {
	if os.Getenv("DDEBUG") == "true" {
		banner := ansi.ColorCode("yellow+h:black") + "tug" + ansi.ColorCode("green:black") + ":" + ansi.ColorCode("reset")
		fmt.Printf(fmt.Sprintf("%s %s", banner, format), a...)
	}
}

func die(err error) {
	banner := ansi.ColorCode("red+h:black") + "tug" + ansi.ColorCode("green:black") + ":" + ansi.ColorCode("reset")
	fmt.Fprintf(os.Stderr, fmt.Sprintf("%s %v\n", banner, err))
	os.Exit(1)
}

func fail(format string, a ...interface{}) {
	banner := ansi.ColorCode("red+h:black") + "tug" + ansi.ColorCode("green:black") + ":" + ansi.ColorCode("reset")
	fmt.Printf(fmt.Sprintf("%s %s", banner, format), a...)
}

func message(format string, a ...interface{}) {
	banner := ansi.ColorCode("blue+h:black") + "tug" + ansi.ColorCode("green:black") + ":" + ansi.ColorCode("reset")
	fmt.Printf(fmt.Sprintf("%s %s", banner, format), a...)
}

func main() {
	defer handlePanic()
	if DockerPs().Run() != nil {
		fmt.Println("docker unavailable")
		os.Exit(0)
	}

	args := os.Args[1:]
	if len(args) < 1 {
		usage()
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Run != nil {
			cmd.Flag.Usage = func() {
				cmd.printUsage()
			}
			if err := cmd.Flag.Parse(args[1:]); err != nil {
				os.Exit(2)
			}
			cmd.Run(cmd, cmd.Flag.Args())
			return
		}
	}
	usage()
}
