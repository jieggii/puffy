package utils

import (
	"errors"
	"github.com/BurntSushi/toml"
	"log"
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

func Execute(command string) (int, error) {
	tokens := strings.Split(command, " ")
	cmd := exec.Command(tokens[0], tokens[1:]...)
	err := cmd.Start()
	if err != nil {
		return 0, errors.New("Error: could not run command '" + command + "'")
	}
	go func() {
		err = cmd.Wait()
		if err != nil {
			log.Println("Error: command '"+command+"' finished with error:", err)
		}
	}()
	return cmd.Process.Pid, nil
}
