package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func DockerBuild(path, name string) (string, *exec.Cmd) {
	abs, _ := filepath.Abs(path)
	tag := fmt.Sprintf("%s.%s", filepath.Base(abs), name)
	return tag, exec.Command("docker", "build", "--rm", "-t", tag, path)
}

func DockerExecInteractive(tag string, command ...string) *exec.Cmd {
	args := []string{"exec", "-it", tag}
	for _, part := range command {
		args = append(args, part)
	}
	return exec.Command("docker", args...)
}

func DockerInspect(tag string, format string) []byte {
	output, _ := exec.Command("docker", "inspect", "-f", format, tag).Output()
	if string(output) == "" {
		cmd := DockerPull(tag)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

		output, _ = exec.Command("docker", "inspect", "-f", format, tag).Output()
	}
	return output
}

func DockerPorts(tag string) []string {
	output := DockerInspect(tag, "{{range $k, $v := .Config.ExposedPorts}}{{$k}}{{end}}")
	ports := []string{}
	for _, exposed := range strings.Split(string(output), " ") {
		ports = append(ports, strings.Split(exposed, "/")[0])
	}
	return ports
}

func DockerPs() *exec.Cmd {
	return exec.Command("docker", "ps")
}

func DockerPull(tag string) *exec.Cmd {
	return exec.Command("docker", "pull", tag)
}

func DockerRun(args ...string) *exec.Cmd {
	runargs := []string{"run"}
	for _, arg := range args {
		runargs = append(runargs, arg)
	}
	return exec.Command("docker", runargs...)
}

func DockerStop(tag string) {
	pull := exec.Command("docker", "stop", tag)
	pull.Run()

	stop := exec.Command("docker", "rm", tag)
	stop.Run()
}
