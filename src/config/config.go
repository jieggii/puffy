package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"puffy/src/utils"
	"strconv"
)

type Repo struct {
	Name   string
	Secret string
	Exec   string
}

type Config struct {
	Host     string
	Port     int
	Endpoint string
	Repos    []Repo `toml:"repo"`
}

func validateConfig(meta toml.MetaData, config *Config) {
	keys := meta.Keys()
	undecoded := meta.Undecoded()
	if len(undecoded) != 0 {
		log.Fatal("error: unexpected", undecoded, "keys in the config file")
	}
	if !utils.KeyIsPresent("port", keys) {
		log.Fatal("error: required field \"port\" is not set in the config file")
	}
	if len(config.Repos) == 0 {
		log.Fatal("error: no repositories specified in the config file")
	}
	for i, repo := range config.Repos {
		if repo.Name == "" {
			log.Fatal("error: missing required field \"name\" in the repo #" + strconv.Itoa(i+1) + " in the config file")
		}
		if repo.Exec == "" {
			log.Fatal("error: missing required field \"exec\" in repo with name \"" + repo.Name + "\" in the config file")
		}
		if repo.Secret == "" {
			log.Fatal("error: missing required secret in repo with name \"" + repo.Name + "\" in the config file")
		}
	}
}

func prepareConfig(meta toml.MetaData, config *Config) *Config {
	keys := meta.Keys()
	if !utils.KeyIsPresent("host", keys) {
		log.Println("setting host to \"0.0.0.0\" as it is not specified in the config file")
		config.Host = "0.0.0.0"
	}
	if !utils.KeyIsPresent("endpoint", keys) {
		log.Println("setting endpoint to \"/\" as it is not specified in the config file")
		config.Endpoint = "/"
	}
	var repoNames []string
	for _, repo := range config.Repos {
		repoNames = append(repoNames, repo.Name)
	}
	log.Println("serving repositories:", repoNames)
	return config
}

func LoadConfig() Config {
	var config Config
	configPath := utils.GetEnv("PUFFY_CONFIG_PATH", "/etc/puffy/config.toml")
	log.Println("using config:", configPath)

	meta, err := toml.DecodeFile(
		configPath,
		&config,
	)
	if err != nil {
		log.Fatal(err)
	}

	prepareConfig(meta, &config)
	validateConfig(meta, &config)

	return config
}
