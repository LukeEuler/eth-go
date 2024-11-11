package main

import (
	"flag"

	"github.com/LukeEuler/dolly/log"

	eg "github.com/LukeEuler/eth-go"
	"github.com/LukeEuler/eth-go/config"
)

func main() {
	log.AddConsoleOut(4)

	configFile := flag.String("c", "config.toml", "set the config file path")
	flag.Parse()
	log.Entry.Infof("config file: %s", *configFile)
	config.New(*configFile)

	eg.Lucky()
	eg.NewKeys()
	eg.Derivation()
	eg.GetBalance()
	eg.Transfer()
}
