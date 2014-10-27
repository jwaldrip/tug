package command

import (
	"os"
	"path/filepath"

	"github.com/nitrous-io/tug/Godeps/_workspace/src/github.com/jwaldrip/odin/cli"
	"github.com/nitrous-io/tug/dockerfile"
	"github.com/nitrous-io/tug/helpers"
	"github.com/nitrous-io/tug/tugfile"
)

var Start = cli.NewSubCommand("start", "Start the application", runStart)

var flagPort int
var flagVerbose bool

var ignores map[string]bool

func init() {
	Start.SetLongDescription(`
Start the application in the current directory.

Examples:

  tug start
  tug start -p 9000 -v
	`)

	Start.DefineIntFlagVar(&flagPort, "port", 5000, "base port")
	Start.AliasFlag('p', "port")
	Start.DefineBoolFlagVar(&flagVerbose, "verbose", false, "show verbose output")
	Start.AliasFlag('v', "verbose")

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
