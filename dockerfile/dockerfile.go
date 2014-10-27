package dockerfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type DockerfileAdd struct {
	Local  string
	Remote string
}

func (dfa *DockerfileAdd) String() string {
	return fmt.Sprintf("{Local:%s Remote:%s}", dfa.Local, dfa.Remote)
}

type Dockerfile struct {
	Add     []*DockerfileAdd
	Command string
	Expose  []string
	Workdir string
}

func New(filename string) (*Dockerfile, error) {
	df := &Dockerfile{}

	f, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		data := scanner.Text()

		line := strings.Trim(string(data), " \t\n")
		parts := strings.Split(line, " ")

		switch strings.ToUpper(parts[0]) {
		case "ADD":
			if len(parts) == 3 {
				df.Add = append(df.Add, &DockerfileAdd{Local: parts[1], Remote: parts[2]})
			}
		case "CMD":
			df.Command = strings.Join(parts[1:len(parts)], " ")
		case "EXPOSE":
			if len(parts) == 2 {
				df.Expose = append(df.Expose, parts[1])
			}
		case "WORKDIR":
			df.Workdir = parts[1]
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return df, nil
}
