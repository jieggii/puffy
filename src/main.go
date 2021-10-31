package main

import (
	"encoding/json"
	"log"
	"net/http"
	"puffy/src/config"
	"puffy/src/json_objects"
	"puffy/src/utils"
	"strconv"
	"strings"
)

func handleRequest(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	var event json_objects.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	knownRepo := false
	for _, repo := range cfg.Repos {
		if repo.Name == event.Repository.FullName {
			// todo: validate payload signature
			// (https://docs.github.com/en/developers/webhooks-and-events/webhooks/securing-your-webhooks)
			command := repo.Exec
			pid, err := utils.Execute(command)
			if err != nil {
				log.Println(err)
			} else {
				log.Println("got new push event from "+repo.Name+". Started process '"+repo.Exec+"'. PID:", pid)
			}
			knownRepo = true
			break
		}
	}
	if !knownRepo {
		log.Println("warning: got push event from unknown repo:", event.Repository.FullName)
	}
}

func main() {
	cfg := config.LoadConfig()
	var repoNames []string
	for _, repo := range cfg.Repos {
		repoNames = append(repoNames, repo.Name)
	}
	log.Println("Serving GitHub repositories:", strings.Join(repoNames[:], ", "))

	mux := http.NewServeMux()
	mux.HandleFunc(cfg.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, r, &cfg)
	})

	addr := cfg.Host + ":" + strconv.Itoa(cfg.Port)
	log.Println("Started puffy at " + addr + cfg.Endpoint)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
