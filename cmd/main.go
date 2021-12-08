package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"puffy/cmd/config"
	"puffy/cmd/json_objects"
	"puffy/cmd/utils"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

func resloveRepo(repoName string, cfg *config.Config) *config.Repo {
	for _, repo := range cfg.Repos {
		if repo.Name == repoName {
			return &repo
		}
	}
	return nil
}

func selectShell(repo *config.Repo, cfg *config.Config) string {
	if repo.Shell == "" {
		return cfg.Shell
	} else {
		return repo.Shell
	}
}

func selectWorkdir(repo *config.Repo, cfg *config.Config) string {
	if repo.Workdir == "" {
		return cfg.Workdir
	} else {
		return repo.Workdir
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read request body", http.StatusBadRequest)
		log.Println("Warning: could not read request body (request from " + utils.GetIP(r) + ")")
		return
	}
	var event json_objects.Event
	event_decode_err := json.NewDecoder(
		ioutil.NopCloser(bytes.NewReader(body)),
	).Decode(&event)

	if event_decode_err != nil {
		var pingEvent json_objects.PingEvent
		ping_event_decode_err := json.NewDecoder(
			ioutil.NopCloser(bytes.NewReader(body)),
		).Decode(&pingEvent)

		if ping_event_decode_err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			log.Println(ping_event_decode_err)
			log.Println("Warning: received invalid request body (request from " + utils.GetIP(r) + ")")
			return
		}
		if resloveRepo(pingEvent.Repository.FullName, cfg) == nil {
			http.Error(w, "Unknown repository", http.StatusForbidden)
			log.Println("Warning: received ping event from unknown repository: " + pingEvent.Repository.FullName)
			return
		}
		w.Write([]byte("pong!"))
		log.Println("Received ping event from " + pingEvent.Repository.FullName)
		return
	}
	repo := resloveRepo(event.Repository.FullName, cfg)
	if repo == nil {
		http.Error(w, "Unknown repository", http.StatusForbidden)
		log.Println("Warning: received push event from unknown repository: " + event.Repository.FullName)
		return
	}
	w.Write([]byte("The event will be handled"))
	log.Println("Received push event from " + repo.Name)

	shell := selectShell(repo, cfg)
	workdir := selectWorkdir(repo, cfg)

	strPID, err := utils.ExecuteCommand(shell, workdir, repo.Exec, repo.Name)
	if err != nil {
		strError := err.Error()
		http.Error(w, strError, http.StatusInternalServerError)
		log.Println(strError)
		return
	}
	log.Println("Spawned process for " + repo.Name + " (PID: " + strPID + ")")
}

func getRepoNames(cfg *config.Config) []string {
	var repoNames []string
	for _, repo := range cfg.Repos {
		repoNames = append(repoNames, repo.Name)
	}
	return repoNames
}

func startServer(c *cli.Context) error {
	cfg := config.LoadConfig(c.String("config"))

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
		Version: "v1.1.0",
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
