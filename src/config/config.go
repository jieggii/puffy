package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"puffy/src/utils"
)

type Repo struct {
	Name   string
	Secret string
	Exec   string
}

type Config struct {
	path     string
	Host     string
	Port     int
	Endpoint string
	Repos    []Repo `toml:"repo"`
}

func validateConfig(meta toml.MetaData, config *Config) {
	keys := meta.Keys()
	undecoded := meta.Undecoded()
	if len(undecoded) != 0 {
		log.Fatal("error: unexpected ", undecoded, "keys in the config file (" + config.path + ")")
	}
	if !utils.KeyIsPresent("port", keys) {
		log.Fatal("error: required variable \"port\" is not set")
	}
	if len(config.Repos) == 0 {
		log.Fatal("error: no repositories specified in the config file (" + config.path + ")")
	}
	for i, repo := range config.Repos {
		if repo.Name == "" {
			log.Fatal("error: missing required field \"name\" in the repo #", i+1)
		}
		if repo.Exec == "" {
			log.Fatal("error: missing required field \"exec\" in repo with name \"" + repo.Name + "\"")
		}
		if repo.Secret == "" {
			log.Fatal("error: missing required secret in repo with name \"" + repo.Name + "\"")
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
		log.Println("setting to endpoint  \"/\" as it is not specified in the config file")
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
	config.path = utils.GetEnv("PUFFY_CONFIG_PATH", "/etc/puffy/config.toml")
	log.Println("using config:", config.path)

	meta, err := toml.DecodeFile(
		config.path,
		&config,
	)
	if err != nil {
		log.Fatal(err)
	}

	prepareConfig(meta, &config)
	validateConfig(meta, &config)

	return config
}
