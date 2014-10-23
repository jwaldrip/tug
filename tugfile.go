package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/nitrous-io/tug/Godeps/_workspace/src/github.com/inconshreveable/muxado"
)

type Tugfile struct {
	BasePort  int
	Docker    bool
	LocalEnv  map[string]string
	DockerEnv map[string]string
	Gateway   string
	Name      string
	Processes []*TugfileProcess
	Root      string

	forward muxado.Session
}

type TugfileProcess struct {
	Name    string
	Adapter string
	Command string

	Ports map[string]string
	Sync  map[string]string
	Tag   string
}

var TugfileLineMatcher = regexp.MustCompile(`^([A-Za-z0-9-_]+):\s*(.*)$`)
var TugfileLineDocker = regexp.MustCompile(`^docker/(.*)$`)

func NewTugfile(filename string) (*Tugfile, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, nil
	}

	f, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	tf, err := TugfileFromReader(f)

	if err != nil {
		return nil, err
	}

	tf.Root = filepath.Dir(filename)

	if _, err = os.Stat(filepath.Join(tf.Root, "Dockerfile")); os.IsNotExist(err) {
		tf.Docker = false
	} else {
		tf.Docker = true
	}

	return tf, nil
}

func DefaultTugfile(df *Dockerfile) (*Tugfile, error) {
	return TugfileFromReader(strings.NewReader(fmt.Sprintf("cmd[public]: local/%s", df.Command)))
}

func TugfileFromReader(reader io.Reader) (*Tugfile, error) {
	tf := &Tugfile{}

	tf.BasePort = 5000
	tf.LocalEnv = make(map[string]string)
	tf.DockerEnv = make(map[string]string)
	tf.Processes = make([]*TugfileProcess, 0)
	tf.Root = ""

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		data := scanner.Text()
		matches := TugfileLineMatcher.FindStringSubmatch(data)
		if len(matches) == 0 {
			continue
		}
		adapter := "local"
		command := matches[2]
		m := TugfileLineDocker.FindStringSubmatch(command)
		if len(m) > 0 {
			adapter = "docker"
			command = m[1]
		}
		p := &TugfileProcess{Name: matches[1], Adapter: adapter, Command: command, Ports: map[string]string{}, Sync: map[string]string{}}
		tf.Processes = append(tf.Processes, p)
	}

	return tf, nil
}

func (tf *Tugfile) Build() {
	gateway, _ := DockerRun("ddollar/docker-gateway").Output()
	tf.Gateway = strings.TrimSpace(string(gateway))

	for _, process := range tf.Processes {
		switch process.Adapter {
		case "docker":
			process.Tag = process.Command
		case "local":
			if tf.Docker {
				tag, cmd := DockerBuild(tf.Root, process.Name)
				process.Tag = tag
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			}
		}
	}
}

func (tf *Tugfile) Forward() {
	DockerStop(tf.DockerName("forward"))
	local, remote := net.Pipe()
	forward := DockerRun("-i", "--privileged", "--name", tf.DockerName("forward"), "nitrousio/docker-forward")
	forward.Stdin = remote
	forward.Stdout = remote
	//forward.Stderr = os.Stdout
	forward.Start()
	tf.forward = muxado.Client(local)

	for _, process := range tf.Processes {
		switch process.Adapter {
		case "docker":
			for i, port := range DockerPorts(process.Tag) {
				ext := tf.portFor(process.Name, i)
				exts := strconv.Itoa(ext)
				process.Ports[exts] = port
				prefix := fmt.Sprintf("%s_PORT_%s_TCP", strings.ToUpper(process.Name), port)

				go tf.forwardPort(ext, fmt.Sprintf("%s:%d", tf.Gateway, ext))

				tf.DockerEnv[prefix] = fmt.Sprintf("tcp://%s:%s", tf.Gateway, exts)
				tf.DockerEnv[fmt.Sprintf("%s_ADDR", prefix)] = tf.Gateway
				tf.DockerEnv[fmt.Sprintf("%s_PORT", prefix)] = exts
				tf.DockerEnv[fmt.Sprintf("%s_PROTO", prefix)] = "tcp"
				tf.LocalEnv[prefix] = fmt.Sprintf("tcp://%s:%s", "127.0.0.1", exts)
				tf.LocalEnv[fmt.Sprintf("%s_ADDR", prefix)] = "127.0.0.1"
				tf.LocalEnv[fmt.Sprintf("%s_PORT", prefix)] = exts
				tf.LocalEnv[fmt.Sprintf("%s_PROTO", prefix)] = "tcp"
			}
		}
	}
}

func (tf *Tugfile) Start(port int) {
	var wg sync.WaitGroup

	for i, process := range tf.Processes {
		wg.Add(1)
		go tf.startProcess(process, tf.portFor(process.Name, i), &wg)
	}

	wg.Wait()
}

func (tf *Tugfile) LongestName() int {
	longest := 0
	for _, process := range tf.Processes {
		if len(process.Name) > longest {
			longest = len(process.Name)
		}
	}
	return longest
}

func (tf *Tugfile) DockerName(name string) string {
	return fmt.Sprintf("%s.%s", tf.Name, name)
}

func (tf *Tugfile) DockerRun(process *TugfileProcess, command string) *exec.Cmd {
	args := []string{"--privileged", "-i", "--name", tf.DockerName(process.Name)}

	for k, v := range tf.DockerEnv {
		args = append(args, "-e")
		args = append(args, fmt.Sprintf("%s=%s", k, v))
	}

	for local, remote := range process.Sync {
		l, _ := filepath.Abs(local)
		r, _ := filepath.Abs(remote)
		args = append(args, "-v")
		args = append(args, fmt.Sprintf("%s:%s", l, r))
	}

	for external, internal := range process.Ports {
		args = append(args, "-p")
		args = append(args, fmt.Sprintf("%s:%s", external, internal))
	}

	args = append(args, process.Tag)

	if command != "" {
		args = append(args, "bash")
		args = append(args, "-c")
		args = append(args, command)
	}

	return DockerRun(args...)
}

func (tf *Tugfile) startProcess(process *TugfileProcess, port int, wg *sync.WaitGroup) {
	r, w := io.Pipe()
	switch process.Adapter {
	case "docker":
		DockerStop(tf.DockerName(process.Name))
		cmd := tf.DockerRun(process, "")
		cmd.Stdout = w
		cmd.Stderr = w
		cmd.Start()
	case "local":
		if tf.Docker {
			DockerStop(tf.DockerName(process.Name))
			df, _ := NewDockerfile(filepath.Join(tf.Root, "Dockerfile"))
			for _, add := range df.Add {
				process.Sync[add.Local] = add.Remote
			}
			cmd := tf.DockerRun(process, process.Command)
			cmd.Stdout = w
			cmd.Stderr = w
			cmd.Start()
		} else {
			cmd := exec.Command("bash", "-c", process.Command)
			cmd.Env = tf.envAsArray(tf.LocalEnv)
			cmd.Env = append(cmd.Env, fmt.Sprintf("PORT=%d", port))
			cmd.Stdout = w
			cmd.Stderr = w
			cmd.Start()
		}
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Printf(fmt.Sprintf("%%-%ds | %%s\n", tf.LongestName()), process.Name, scanner.Text())
	}
	wg.Done()
}

func (tf *Tugfile) forwardPort(port int, dest string) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		die(err)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()

		if err != nil {
			panic(err)
		}

		go tf.forwardConnection(conn, dest)
	}
}

func (tf *Tugfile) forwardConnection(local net.Conn, dest string) {
	remote, _ := tf.forward.Open()
	remote.Write([]byte("connect\n"))
	remote.Write([]byte(fmt.Sprintf("%s\n", dest)))
	remote.Read(make([]byte, 2))
	PipeStreams(remote, local)
}

func (tf *Tugfile) portFor(name string, idx int) int {
	for psidx, ps := range tf.Processes {
		if ps.Name == name {
			return tf.BasePort + (psidx * 100) + idx
		}
	}
	return -1
}

func (tf *Tugfile) envAsArray(in map[string]string) (out []string) {
	for _, pair := range os.Environ() {
		out = append(out, pair)
	}
	for name, val := range in {
		out = append(out, fmt.Sprintf("%s=%s", name, val))
	}
	return
}
