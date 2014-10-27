package commands

import (
	"os"
	"path/filepath"

	"github.com/nitrous-io/tug/Godeps/_workspace/src/github.com/jwaldrip/odin/cli"
	"github.com/nitrous-io/tug/dockerfile"
	"github.com/nitrous-io/tug/helpers"
	"github.com/nitrous-io/tug/tugfile"
)

var cmdStart = cli.NewSubCommand("start", "Start the application", runStart)

var flagForward stringSet
var flagHost string
var flagPort int
var flagSync stringSet
var flagTag string
var flagVerbose bool

var ignores map[string]bool

func init() {
	cmdStart.SetLongDescription(`
Start the application in the current directory.

Examples:

  tug start
  tug start -f 5000:3000 -s .:/app
  tug start -c "make run"
	`)

	wd, _ := os.Getwd()

	cmdStart.DefineStringFlagVar(&flagHost, "host", "", "nitro host to use")
	cmdStart.AliasFlag('h', "host")
	cmdStart.DefineFlag(&flagForward, "forward", "local:remote port to forward")
	cmdStart.AliasFlag('f', "forward")
	cmdStart.DefineIntFlagVar(&flagPort, "port", 5000, "base port")
	cmdStart.AliasFlag('p', "port")
	cmdStart.DefineFlag(&flagSync, "sync", "local:remote file or directory to sync")
	cmdStart.AliasFlag('s', "sync")
	cmdStart.DefineStringFlagVar(&flagTag, "tag", filepath.Base(wd), "docker tag to use")
	cmdStart.AliasFlag('t', "tag")
	cmdStart.DefineBoolFlagVar(&flagVerbose, "verbose", false, "show verbose output")
	cmdStart.AliasFlag('v', "verbose")

	ignores = make(map[string]bool)
}

func runStart(c cli.Command) {
	tf, err := tugfile.New("./Tugfile")

	if err != nil {
		helpers.Die(err)
	}

	if tf == nil {
		df, _ := dockerfile.New("./Dockerfile")
		tf, err = tugfile.Default(df)

		if err != nil {
			helpers.Die(err)
		}
	}

	wd, _ := os.Getwd()
	abs, _ := filepath.Abs(wd)
	tf.Name = filepath.Base(abs)

	tf.Build()
	tf.ResolveLinks()
	tf.Start(flagPort)
}
