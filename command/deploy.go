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

var Deploy = cli.NewSubCommand("deploy", "Deploy to a production service", runDeploy)

func init() {
	Deploy.SetLongDescription(`
Deploy to a production service

Examples:

  tug deploy tutum.co/myuser/myapp
	`)
	Deploy.DefineParams("prefix")
}

func runDeploy(c cli.Command) {
	prefix := c.Param("prefix").Get()

	tf, err := tugfile.New("./Tugfile")

	if err != nil {
		helpers.Die(err)
	}

	if tf == nil {
		helpers.Die(fmt.Errorf("no Tugfile found"))
	}

	if !tf.HasDockerfile {
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
		helpers.Message("deploying %s\n", tag)
		switch process.Adapter {
		case "docker":
			cmd := docker.Tag(process.Command, tag)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
			cmd = docker.Push(tag)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		case "local":
			cmd := docker.Build(tf.Root, tag)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
			cmd = docker.Push(tag)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		}
	}
}
