package main

import (
	"errors"
	"log"
	"net/http"
	"os/exec"
	"strconv"
)

func resloveRepo(repoName string, cfg *Config) *Repo {
	for _, repo := range cfg.Repos {
		if repo.Name == repoName {
			return &repo
		}
	}
	return nil
}

func selectShell(repo *Repo, cfg *Config) string {
	if repo.Shell == "" {
		return cfg.Shell
	} else {
		return repo.Shell
	}
}

func selectWorkdir(repo *Repo, cfg *Config) string {
	if repo.Workdir == "" {
		return cfg.Workdir
	} else {
		return repo.Workdir
	}
}

func executeCommand(shell string, workdir string, command string, repoName string) (string, error) {
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

func getIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
