package configuration

import (
	"github.com/BurntSushi/toml"
	"log"
)

// Info from config file
type Config struct {
	Dsn     string
	Blacklistdir string
	Category []string
	Blockcategory []string
}

// Reads info from config file
func ReadConfig(configfile string) Config {
	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}
	return config
}
