package utils

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"os/exec"
	"strings"
)

func KeyIsPresent(name string, keys []toml.Key) bool {
	for _, key := range keys {
		if key.String() == name {
			return true
		}
	}
	return false
}

func GetEnv(key string, fallback string) string {
	value, isPresent := os.LookupEnv(key)
	if isPresent {
		return value
	} else {
		return fallback
	}
}

func Execute(command string) int {
	tokens := strings.Split(command, " ")

	var cmd *exec.Cmd
	if len(tokens) == 1 {
		cmd = exec.Command(tokens[0])
	} else {
		cmd = exec.Command(tokens[0], tokens[1:]...)
	}

	err := cmd.Start()
	if err != nil {
		return 0
	}
	pid := cmd.Process.Pid
	go func() {
		err = cmd.Wait()
		if err != nil {
			fmt.Printf("command \"" + command + "\" finished with error:", err)
		}
	}()
	return pid
}
