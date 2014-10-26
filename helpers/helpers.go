package helpers

import (
	"fmt"
	"os"

	"github.com/nitrous-io/tug/Godeps/_workspace/src/github.com/mgutz/ansi"
)

func Debug(format string, a ...interface{}) {
	if os.Getenv("DDEBUG") == "true" {
		banner := ansi.ColorCode("yellow+h:black") + "tug" + ansi.ColorCode("green:black") + ":" + ansi.ColorCode("reset")
		fmt.Printf(fmt.Sprintf("%s %s", banner, format), a...)
	}
}

func Die(err error) {
	banner := ansi.ColorCode("red+h:black") + "tug" + ansi.ColorCode("green:black") + ":" + ansi.ColorCode("reset")
	fmt.Fprintf(os.Stderr, fmt.Sprintf("%s %v\n", banner, err))
	os.Exit(1)
}

func Fail(format string, a ...interface{}) {
	banner := ansi.ColorCode("red+h:black") + "tug" + ansi.ColorCode("green:black") + ":" + ansi.ColorCode("reset")
	fmt.Printf(fmt.Sprintf("%s %s", banner, format), a...)
}

func Message(format string, a ...interface{}) {
	banner := ansi.ColorCode("blue+h:black") + "tug" + ansi.ColorCode("green:black") + ":" + ansi.ColorCode("reset")
	fmt.Printf(fmt.Sprintf("%s %s", banner, format), a...)
}
