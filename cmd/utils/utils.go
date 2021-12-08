package utils

import (
	"errors"
	"log"
	"net/http"
	"os/exec"
	"strconv"
)

func ExecuteCommand(shell string, workdir string, command string, repoName string) (string, error) {
	cmd := exec.Command(shell, "-c", "cd '"+workdir+"' && "+command)
	err := cmd.Start()
	if err != nil {
		return "", errors.New("Error: could not spawn process for " + repoName + ": " + err.Error())
	}
	strPID := strconv.Itoa(cmd.Process.Pid)
	go func() {
		err = cmd.Wait()
		if err != nil {
			log.Println("Error: process for " + repoName + " (PID: " + strPID + ") finished with error: " + err.Error())
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
