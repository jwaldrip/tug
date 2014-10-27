package tugfile

import (
	"bufio"
	"os"
	"strings"
)

type Env map[string]string

func ReadEnv(filename string) (Env, error) {
	env := make(Env)

	f, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), "=", 2)
		env[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return env, nil
}
