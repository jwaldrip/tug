package main

import (
	"os"
	"path/filepath"
)

var cmdStart = &Command{
	Run:   runStart,
	Usage: "start [-v] [-c command] [-p local:container] [-s local:container] [-t tag]",
	Short: "Start the application",
	Long: `
Start the application in the current directory.

Examples:

  tug start
  tug start -f 5000:3000 -s .:/app
  tug start -c "make run"
`,
}

var flagForward StringSet
var flagHost string
var flagPort int
var flagSync StringSet
var flagTag string
var flagVerbose bool

var ignores map[string]bool

func init() {
	wd, _ := os.Getwd()

	cmdStart.Flag.StringVar(&flagHost, "h", "", "nitro host to use")
	cmdStart.Flag.Var(&flagForward, "f", "local:remote port to forward")
	cmdStart.Flag.IntVar(&flagPort, "p", 5000, "base port")
	cmdStart.Flag.Var(&flagSync, "s", "local:remote file or directory to sync")
	cmdStart.Flag.StringVar(&flagTag, "t", filepath.Base(wd), "docker tag to use")
	cmdStart.Flag.BoolVar(&flagVerbose, "v", false, "show verbose output")

	ignores = make(map[string]bool)
}

func runStart(c *Command, args []string) {
	tf, err := NewTugfile("./Tugfile")

	if err != nil {
		die(err)
	}

	if tf == nil {
		df, _ := NewDockerfile("./Dockerfile")
		tf, err = DefaultTugfile(df)

		if err != nil {
			die(err)
		}
	}

	wd, _ := os.Getwd()
	abs, _ := filepath.Abs(wd)
	tf.Name = filepath.Base(abs)

	tf.Build()
	tf.Forward()
	tf.Start(flagPort)
}
