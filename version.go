package main

import (
	"fmt"
)

var VERSION = "dev"

var cmdVersion = &Command{
	Run:   runVersion,
	Usage: "version",
	Short: "Display current version",
	Long: `
Display current version

Examples:

	tug version
`,
}

func init() {
}

func runVersion(cmd *Command, args []string) {
	fmt.Println(VERSION)
}
