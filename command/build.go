package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nitrous-io/tug/Godeps/_workspace/src/github.com/jwaldrip/odin/cli"
	"github.com/nitrous-io/tug/docker"
	"github.com/nitrous-io/tug/helpers"
	"github.com/nitrous-io/tug/tugfile"
)

var Build = cli.NewSubCommand("build", "Build docker images for production", runBuild)

func init() {
	Build.SetLongDescription(`
Build docker images for production.

Examples:

  tug build
  tug build -p myapp
	`)
	Build.DefineStringFlag("prefix", "", "prefix for docker image tags")
	Build.AliasFlag('p', "prefix")
}

func runBuild(c cli.Command) {
	prefix := c.Flag("prefix").Get()
	fmt.Printf("prefix %+v\n", prefix)

	tf, err := tugfile.New("./Tugfile")

	if err != nil {
		helpers.Die(err)
	}

	if tf == nil {
		helpers.Die(fmt.Errorf("no Tugfile found"))
	}

	if !tf.Docker {
		for _, process := range tf.Processes {
			if process.Adapter == "local" {
				helpers.Die(fmt.Errorf("no Dockerfile found"))
			}
		}
	}

	if prefix == "" {
		abs, _ := filepath.Abs(tf.Root)
		prefix = filepath.Base(abs)
	}

	for _, process := range tf.Processes {
		tag := fmt.Sprintf("%s.%s", prefix, process.Name)
		helpers.Message("creating %s\n", tag)
		switch process.Adapter {
		case "docker":
			cmd := docker.Tag(process.Command, tag)
			cmd.Run()
		case "local":
			cmd := docker.Build(tf.Root, tag)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		}
	}
}
