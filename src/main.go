package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"puffy/src/config"
	"puffy/src/json_objects"
	"puffy/src/utils"
	"strconv"
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
			pid := utils.Execute(repo.Exec)
			fmt.Println("got push event from", repo.Name, ". Executing command", repo.Exec, ". PID:", pid)
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

	mux := http.NewServeMux()
	mux.HandleFunc(cfg.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, r, &cfg)
	})

	log.Println("starting puffy on " + cfg.Host + ":" + strconv.Itoa(cfg.Port))
	err := http.ListenAndServe(cfg.Host+":"+strconv.Itoa(cfg.Port), mux)
	if err != nil {
		log.Fatal(err)
	}
}
