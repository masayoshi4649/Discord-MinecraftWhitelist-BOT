package main

import (
	"github.com/BurntSushi/toml"
)

// conftoml ... conf.toml
type tomlconf struct {
	Bot struct {
		Token string
	}
	Discord struct {
		Channel []string
	}
}

// getConf ... Get config
func getConf() tomlconf {
	var c tomlconf
	_, err := toml.DecodeFile("./conf.toml", &c)
	if err != nil {
		panic(err)
	}

	return c
}
