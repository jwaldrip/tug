package main

import "os"

var cmdInit = &Command{
	Run:   runInit,
	Usage: "init",
	Short: "Initialize tug",
	Long: `
Initialize tug.

Examples:

  tug attach
`,
}

func runInit(c *Command, args []string) {
	cmd := DockerPull("nitrousio/docker-forward")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	cmd = DockerPull("ddollar/docker-gateway")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
