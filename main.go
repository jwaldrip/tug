package main

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

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

func checkDockerHost() {
	var dockerURL *url.URL
	var err error
	var conn net.Conn
	var status string

	dockerURL, err = url.Parse(os.Getenv("DOCKER_HOST"))
	if err == nil {
		conn, err = net.Dial(dockerURL.Scheme, dockerURL.Host)
		if err == nil {
			conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
			conn.SetReadDeadline(time.Now().Add(2 * time.Second))
			fmt.Fprintf(conn, "GET /_ping HTTP/1.0\r\n\r\n")
			status, err = bufio.NewReader(conn).ReadString('\n')
			if !strings.Contains(status, "HTTP/1.0 200 OK") {
				err = fmt.Errorf("invalid status: %s", status)
			}
		}
	}

	if err != nil {
		fmt.Printf("invalid docker host: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	defer handlePanic()
	if os.Getenv("DOCKER_HOST") == "" {
		os.Setenv("DOCKER_HOST", "unix:///var/lib/docker.sock")
	}

	checkDockerHost()

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
