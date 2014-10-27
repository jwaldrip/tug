package docker

import (
	"os"
	"os/exec"
	"strings"
)

func Build(path, tag string) *exec.Cmd {
	return exec.Command("docker", "build", "--rm", "-t", tag, path)
}

func ExecInteractive(tag string, command ...string) *exec.Cmd {
	args := []string{"exec", "-it", tag}
	for _, part := range command {
		args = append(args, part)
	}
	return exec.Command("docker", args...)
}

func Inspect(tag string, format string) []byte {
	output, _ := exec.Command("docker", "inspect", "-f", format, tag).Output()
	if string(output) == "" {
		cmd := Pull(tag)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

		output, _ = exec.Command("docker", "inspect", "-f", format, tag).Output()
	}
	return output
}

func Ports(tag string) []string {
	output := Inspect(tag, "{{range $k, $v := .Config.ExposedPorts}}{{$k}}{{end}}")
	ports := []string{}
	for _, exposed := range strings.Split(string(output), " ") {
		ports = append(ports, strings.Split(exposed, "/")[0])
	}
	return ports
}

func Ps() *exec.Cmd {
	return exec.Command("docker", "ps")
}

func Pull(tag string) *exec.Cmd {
	return exec.Command("docker", "pull", tag)
}

func Run(args ...string) *exec.Cmd {
	runargs := []string{"run"}
	for _, arg := range args {
		runargs = append(runargs, arg)
	}
	return exec.Command("docker", runargs...)
}

func Stop(tag string) {
	pull := exec.Command("docker", "stop", tag)
	pull.Run()

	stop := exec.Command("docker", "rm", tag)
	stop.Run()
}

func Tag(from, to string) *exec.Cmd {
	return exec.Command("docker", "tag", "-f", from, to)
}
