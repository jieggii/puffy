package config

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

func KeyIsPresent(name string, keys []toml.Key) bool {
	for _, key := range keys {
		if key.String() == name {
			return true
		}
	}
	return false
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

func validateConfig(meta toml.MetaData, config *Config) {
	keys := meta.Keys()
	undecoded := meta.Undecoded()
	if len(undecoded) != 0 {
		log.Fatal("Fatal: unexpected ", undecoded, " keys in the config file")
	}
	if !KeyIsPresent("port", keys) {
		log.Fatal("Fatal: required field 'port' is not set in the config file")
	}
	if len(config.Repos) == 0 {
		log.Fatal("Fatal: no repositories specified in the config file")
	}
	if !pathExists(config.Shell) {
		log.Fatal("Fatal: 'shell' specified in the config file (" + config.Shell + ") does not exist")
	}
	if !pathExists(config.Workdir) {
		log.Fatal("Fatal: 'workdir' specified in the config file )" + config.Workdir + ") does not exist")
	}
	for i, repo := range config.Repos {
		if repo.Name == "" {
			log.Fatal("Fatal: missing required field 'name' in the repo #" + strconv.Itoa(i+1) + " in the config file")
		}
		if repo.Exec == "" {
			log.Fatal("Fatal: missing required field 'exec' in repo with name '" + repo.Name + "' in the config file")
		}
		if repo.Workdir != "" && !pathExists(repo.Workdir) {
			log.Fatal("Fatal: 'workdir' specified in the config file (" + repo.Workdir + ") of repo " + repo.Name + " does not exist")
		}
		if repo.Shell != "" && !pathExists(repo.Shell) {
			log.Fatal("Fatal: 'shell' specified in the config file (" + repo.Shell + ") of repo " + repo.Name + " does not exist")
		}
	}
}

func prepareConfig(meta toml.MetaData, config *Config) *Config {
	keys := meta.Keys()

	if !KeyIsPresent("host", keys) {
		log.Println("Setting 'host' to '0.0.0.0' as it is not specified in the config")
		config.Host = "0.0.0.0"
	}
	if !KeyIsPresent("endpoint", keys) {
		log.Println("Setting 'endpoint' to '/' as it is not specified in the config")
		config.Endpoint = "/"
	}
	if !KeyIsPresent("workdir", keys) {
		log.Println("Setting 'workdir' to '/' as it is not specified in the config")
		config.Workdir = "/"
	}
	if !KeyIsPresent("shell", keys) {
		log.Println("Setting 'shell' to '/usr/bin/sh' as it is not specified in the config")
		config.Shell = "/usr/bin/sh"
	}
	return config
}

func LoadConfig(configPath string) *Config {
	var config Config
	log.Println("Using config:", configPath)
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
