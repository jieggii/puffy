package utils

import (
	"errors"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

func KeyIsPresent(name string, keys []toml.Key) bool {
	for _, key := range keys {
		if key.String() == name {
			return true
		}
	}
	return false
}

func ExecuteCommand(command string, repoName string) (string, error) {
	tokens := strings.Split(command, " ")

	cmd := exec.Command(tokens[0], tokens[1:]...)
	err := cmd.Start()
	if err != nil {
		return "", errors.New("Error: could not spawn process for " + repoName + ": " + err.Error())
	}

	strPID := strconv.Itoa(cmd.Process.Pid)
	go func() {
		err = cmd.Wait()
		if err != nil {
			log.Println("Error: process for " + repoName + "(PID: " + strPID + ") finished with error: " + err.Error())
		} else {
			log.Println("Process for " + repoName + " (PID: " + strPID + ") was sucessfully finished")
		}
	}()

	return strPID, nil
}

func GetIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
