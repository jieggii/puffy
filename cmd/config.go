package main

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
)

type Repo struct {
	Name    string
	Workdir string
	Shell   string
	Exec    string
}

type Config struct {
	Host     string
	Port     int
	Endpoint string
	Workdir  string
	Shell    string
	Repos    []Repo
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

func getRepoNames(cfg *Config) []string {
	var repoNames []string
	for _, repo := range cfg.Repos {
		repoNames = append(repoNames, repo.Name)
	}
	return repoNames
}

func keyIsPresent(name string, keys []toml.Key) bool {
	for _, key := range keys {
		if key.String() == name {
			return true
		}
	}
	return false
}

func validateConfig(meta toml.MetaData, config *Config) {
	keys := meta.Keys()
	undecoded := meta.Undecoded()
	if len(undecoded) != 0 {
		log.Fatal("Fatal: unexpected ", undecoded, " keys in the config")
	}
	if !keyIsPresent("port", keys) {
		log.Fatal("Fatal: required field 'port' is not set in the config")
	}
	if len(config.Repos) == 0 {
		log.Fatal("Fatal: no repositories specified in the config")
	}
	if !pathExists(config.Shell) {
		log.Fatal("Fatal: 'shell' specified in the config (" + config.Shell + ") does not exist")
	}
	if !pathExists(config.Workdir) {
		log.Fatal("Fatal: 'workdir' specified in the config (" + config.Workdir + ") does not exist")
	}
	for i, repo := range config.Repos {
		if repo.Name == "" {
			log.Fatal("Fatal: missing required field 'name' in the repo #" + strconv.Itoa(i+1) + " in the config")
		}
		if repo.Exec == "" {
			log.Fatal("Fatal: missing required field 'exec' in repo with name '" + repo.Name + "' in the config")
		}
		if repo.Workdir != "" && !pathExists(repo.Workdir) {
			log.Fatal("Fatal: 'workdir' specified in the config (" + repo.Workdir + ") of repo " + repo.Name + " does not exist")
		}
		if repo.Shell != "" && !pathExists(repo.Shell) {
			log.Fatal("Fatal: 'shell' specified in the config (" + repo.Shell + ") of repo " + repo.Name + " does not exist")
		}
	}
}

func prepareConfig(meta toml.MetaData, config *Config) *Config {
	keys := meta.Keys()

	if !keyIsPresent("host", keys) {
		log.Println("Setting 'host' to '0.0.0.0' as it is not specified in the config")
		config.Host = "0.0.0.0"
	}
	if !keyIsPresent("endpoint", keys) {
		log.Println("Setting 'endpoint' to '/' as it is not specified in the config")
		config.Endpoint = "/"
	}
	if !keyIsPresent("workdir", keys) {
		log.Println("Setting 'workdir' to '/' as it is not specified in the config")
		config.Workdir = "/"
	}
	if !keyIsPresent("shell", keys) {
		log.Println("Setting 'shell' to '/usr/bin/sh' as it is not specified in the config")
		config.Shell = "/usr/bin/sh"
	}
	return config
}

func loadConfig(configPath string) *Config {
	var config Config
	log.Println("Using config file:", configPath)
	meta, err := toml.DecodeFile(
		configPath,
		&config,
	)
	if err != nil {
		log.Fatal(err)
	}
	prepareConfig(meta, &config)
	validateConfig(meta, &config)
	return &config
}
