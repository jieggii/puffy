package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

func handleRequest(w http.ResponseWriter, r *http.Request, cfg *Config) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Warning: could not read request body (request from " + getIP(r) + ")")
		http.Error(w, "Could not read request body", http.StatusBadRequest)
		return
	}
	var event Event
	event_decode_err := json.NewDecoder(
		ioutil.NopCloser(bytes.NewReader(body)),
	).Decode(&event)

	if event_decode_err != nil {
		var pingEvent PingEvent
		ping_event_decode_err := json.NewDecoder(
			ioutil.NopCloser(bytes.NewReader(body)),
		).Decode(&pingEvent)

		if ping_event_decode_err != nil {
			log.Println("Warning: received invalid request body (request from " + getIP(r) + ")")
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if resloveRepo(pingEvent.Repository.FullName, cfg) == nil {
			log.Println("Warning: received ping event from unknown repository: " + pingEvent.Repository.FullName)
			http.Error(w, "Unknown repository", http.StatusForbidden)
			return
		}
		log.Println("Received ping event from " + pingEvent.Repository.FullName)
		w.Write([]byte("pong!"))
		return
	}
	repo := resloveRepo(event.Repository.FullName, cfg)
	if repo == nil {
		log.Println("Warning: received push event from unknown repository: " + event.Repository.FullName)
		http.Error(w, "Unknown repository", http.StatusForbidden)
		return
	}
	log.Println("Received push event from " + repo.Name)

	shell := selectShell(repo, cfg)
	workdir := selectWorkdir(repo, cfg)

	strPID, err := executeCommand(shell, workdir, repo.Exec, repo.Name)
	if err != nil {
		strError := err.Error()
		log.Println(strError)
		http.Error(w, strError, http.StatusInternalServerError)
		return
	}
	log.Println("Spawned process for " + repo.Name + " (PID: " + strPID + ")")
	w.Write([]byte("The event will be handled"))
}

func startServer(c *cli.Context) error {
	cfg := loadConfig(c.String("config"))

	repoNames := getRepoNames(cfg)
	log.Println("Serving GitHub repositories:", strings.Join(repoNames[:], ", "))

	addr := cfg.Host + ":" + strconv.Itoa(cfg.Port)
	log.Println("Started puffy at " + addr + cfg.Endpoint)

	mux := http.NewServeMux()
	mux.HandleFunc(cfg.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, r, cfg)
	})
	return http.ListenAndServe(addr, mux)
}

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "show puffy version",
	}
	app := &cli.App{
		Name:    "puffy",
		Version: "v2.0.1",
		Usage:   "simple GitHub webhook listener for push events",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "/etc/puffy/config.toml",
				Usage:   "path to puffy config file",
				EnvVars: []string{"PUFFY_CONFIG_PATH"},
			},
		},
		Action: func(c *cli.Context) error {
			return startServer(c)
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
