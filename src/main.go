package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"puffy/src/config"
	"puffy/src/json_objects"
	"puffy/src/utils"
	"strconv"
	"strings"
)

func handleRequest(writer http.ResponseWriter, request *http.Request, cfg *config.Config) {
	var event json_objects.Event
	err := json.NewDecoder(request.Body).Decode(&event)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
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
				log.Println("Got new push event from "+repo.Name+". Started process '"+repo.Exec+"'. PID:", pid)
			}
			knownRepo = true
			break
		}
	}
	if !knownRepo {
		log.Println("Warning: got push event from unknown repo:", event.Repository.FullName)
	}
}

func main() {
	displayVersion := flag.Bool("version", false, "display puffy version and exit")

	flag.Parse()

	if *displayVersion == true {
		fmt.Println("puffy version 0.0.1")
		return
	}

	cfg := config.LoadConfig()

	var repoNames []string
	for _, repo := range cfg.Repos {
		repoNames = append(repoNames, repo.Name)
	}
	log.Println("Serving GitHub repositories:", strings.Join(repoNames[:], ", "))

	mux := http.NewServeMux()
	mux.HandleFunc(cfg.Endpoint, func(writer http.ResponseWriter, request *http.Request) {
		handleRequest(writer, request, &cfg)
	})

	addr := cfg.Host + ":" + strconv.Itoa(cfg.Port)
	log.Println("Started puffy at " + addr + cfg.Endpoint)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
