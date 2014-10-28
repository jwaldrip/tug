package tugfile

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
	"github.com/nitrous-io/tug/docker"
	"github.com/nitrous-io/tug/dockerfile"
	"github.com/nitrous-io/tug/helpers"
)

type Tugfile struct {
	BasePort  int
	Docker    bool
	Env       map[string]string
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

	Env   map[string]string
	Ports map[string]string
	Sync  map[string]string
	Tag   string
}

var TugfileLineMatcher = regexp.MustCompile(`^([A-Za-z0-9-_]+):\s*(.*)$`)
var TugfileLineDocker = regexp.MustCompile(`^docker/(.*)$`)

func New(filename string) (*Tugfile, error) {
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

	if _, err = os.Stat(filepath.Join(tf.Root, ".env")); !os.IsNotExist(err) {
		env, _ := ReadEnv(filepath.Join(tf.Root, ".env"))
		for key, val := range env {
			tf.Env[key] = val
		}
	}

	return tf, nil
}

func Default(df *dockerfile.Dockerfile) (*Tugfile, error) {
	return TugfileFromReader(strings.NewReader(fmt.Sprintf("cmd[public]: local/%s", df.Command)))
}

func TugfileFromReader(reader io.Reader) (*Tugfile, error) {
	tf := &Tugfile{}

	tf.BasePort = 5000
	tf.Env = make(map[string]string)
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
		p := &TugfileProcess{Name: matches[1], Adapter: adapter, Command: command, Env: map[string]string{}, Ports: map[string]string{}, Sync: map[string]string{}}
		tf.Processes = append(tf.Processes, p)
	}

	return tf, nil
}

func (tf *Tugfile) Build() {
	gateway, _ := docker.Run("ddollar/docker-gateway").Output()
	tf.Gateway = strings.TrimSpace(string(gateway))

	for _, process := range tf.Processes {
		switch process.Adapter {
		case "docker":
			process.Tag = process.Command
		case "local":
			if tf.Docker {
				abs, _ := filepath.Abs(tf.Root)
				process.Tag = fmt.Sprintf("%s.%s", filepath.Base(abs), process.Name)
				cmd := docker.Build(tf.Root, process.Tag)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			}
		}
	}
}

func (tf *Tugfile) ResolveLinks() {
	docker.Stop(tf.DockerName("forward"))
	local, remote := net.Pipe()
	forward := docker.Run("-i", "--privileged", "--name", tf.DockerName("forward"), "nitrousio/docker-forward")
	forward.Stdin = remote
	forward.Stdout = remote
	forward.Start()
	tf.forward = muxado.Client(local)

	links := make(map[string]string)

	for psidx, ps := range tf.Processes {
		switch ps.Adapter {
		case "docker":
			for portidx, port := range docker.Ports(ps.Tag) {
				ext := tf.BasePort + (psidx * 100) + portidx
				prefix := fmt.Sprintf("%s_PORT_%s_TCP", strings.ToUpper(ps.Name), port)
				links[prefix] = strconv.Itoa(ext)
				ps.Ports[strconv.Itoa(ext)] = port
				go tf.forwardPort(ext, fmt.Sprintf("%s:%d", tf.Gateway, ext))
			}
		case "local":
			if tf.Docker {
				for portidx, port := range docker.Ports(ps.Tag) {
					ext := tf.BasePort + (psidx * 100) + portidx
					prefix := fmt.Sprintf("%s_PORT_%s_TCP", strings.ToUpper(ps.Name), port)
					links[prefix] = strconv.Itoa(ext)
					ps.Ports[strconv.Itoa(ext)] = port
					go tf.forwardPort(ext, fmt.Sprintf("%s:%d", tf.Gateway, ext))
				}
			} else {
				port := tf.BasePort + (psidx * 100)
				prefix := fmt.Sprintf("%s_PORT_%d_TCP", strings.ToUpper(ps.Name), port)
				links[prefix] = strconv.Itoa(port)
			}
		}
		if !(ps.Adapter == "local" && tf.Docker) {
		}
	}

	for _, ps := range tf.Processes {
		for prefix, port := range links {
			if ps.Adapter == "local" && !tf.Docker {
				ps.Env[prefix] = fmt.Sprintf("tcp://127.0.0.1:%s", port)
				ps.Env[fmt.Sprintf("%s_ADDR", prefix)] = "127.0.0.1"
				ps.Env[fmt.Sprintf("%s_PORT", prefix)] = port
				ps.Env[fmt.Sprintf("%s_PROTO", prefix)] = "tcp"
			} else {
				ps.Env[prefix] = fmt.Sprintf("tcp://%s:%s", tf.Gateway, port)
				ps.Env[fmt.Sprintf("%s_ADDR", prefix)] = tf.Gateway
				ps.Env[fmt.Sprintf("%s_PORT", prefix)] = port
				ps.Env[fmt.Sprintf("%s_PROTO", prefix)] = "tcp"
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

	for k, v := range tf.Env {
		args = append(args, "-e")
		args = append(args, fmt.Sprintf("%s=%s", k, v))
	}

	for k, v := range process.Env {
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

	return docker.Run(args...)
}

func (tf *Tugfile) startProcess(process *TugfileProcess, port int, wg *sync.WaitGroup) {
	r, w := io.Pipe()
	switch process.Adapter {
	case "docker":
		docker.Stop(tf.DockerName(process.Name))
		cmd := tf.DockerRun(process, "")
		cmd.Stdout = w
		cmd.Stderr = w
		cmd.Start()
	case "local":
		bootstrapCommand := fmt.Sprintf("if [ -x bin/bootstrap ]; then bin/bootstrap; fi; %s", process.Command)
		if tf.Docker {
			docker.Stop(tf.DockerName(process.Name))
			df, _ := dockerfile.New(filepath.Join(tf.Root, "Dockerfile"))
			for _, add := range df.Add {
				process.Sync[add.Local] = add.Remote
			}
			cmd := tf.DockerRun(process, bootstrapCommand)
			cmd.Stdout = w
			cmd.Stderr = w
			cmd.Start()
		} else {
			cmd := exec.Command("bash", "-c", bootstrapCommand)
			cmd.Env = os.Environ()
			for _, item := range tf.envAsArray(tf.Env) {
				cmd.Env = append(cmd.Env, item)
			}
			for _, item := range tf.envAsArray(process.Env) {
				cmd.Env = append(cmd.Env, item)
			}
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
		helpers.Die(err)
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
	for name, val := range in {
		out = append(out, fmt.Sprintf("%s=%s", name, val))
	}
	return
}
