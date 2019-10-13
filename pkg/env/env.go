package env

import (
	"io/ioutil"
	"os"
)

// ReadSecretEnv Attempts to get the value of an environment variable, else reads content of a file with a path stored in the same environment variable ending with _FILE.
func ReadSecretEnv(name string) string {
	env := os.Getenv(name)
	if env != "" {
		return env
	}

	file := os.Getenv(name + "_FILE")
	if file == "" {
		return ""
	}

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}

	return string(bytes)
}
