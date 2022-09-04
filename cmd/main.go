package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/csaldiasdev/distws/internal/agent"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const defaultPortValue = 0

func main() {
	RepoRpcPort := flag.Uint("repoRpcPort", defaultPortValue, "a valid port")
	HubRpcPort := flag.Uint("hubRpcPort", defaultPortValue, "a valid port")
	RaftPort := flag.Uint("raftPort", defaultPortValue, "a valid port")
	HttpPort := flag.Uint("httpPort", defaultPortValue, "a valid port")
	SerfPort := flag.Uint("serfPort", defaultPortValue, "a valid port")
	SerfMember := flag.String("member", "", "a serf address")

	prettylog := flag.Bool("prettylog", false, "a boolean")

	flag.Parse()

	if *prettylog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if *RepoRpcPort == defaultPortValue {
		log.Error().
			Uint("configValue", *RepoRpcPort).
			Msg("Invalid o missing RepoRpcPort value")

		os.Exit(1)
	}

	if *HubRpcPort == defaultPortValue {
		log.Error().
			Uint("configValue", *HubRpcPort).
			Msg("Invalid o missing HubRpcPort value")

		os.Exit(1)
	}

	if *RaftPort == defaultPortValue {
		log.Error().
			Uint("configValue", *RaftPort).
			Msg("Invalid o missing RaftPort value")

		os.Exit(1)
	}

	if *HttpPort == defaultPortValue {
		log.Error().
			Uint("configValue", *HttpPort).
			Msg("Invalid o missing HttpPort value")

		os.Exit(1)
	}

	if *SerfPort == defaultPortValue {
		log.Error().
			Uint("configValue", *SerfPort).
			Msg("Invalid o missing SerfPort value")

		os.Exit(1)
	}

	config := agent.AgentConfiguration{
		RepositoryRpcPort: *RepoRpcPort,
		HubRpcPort:        *HubRpcPort,
		RaftPort:          *RaftPort,
		SerfPort:          *SerfPort,
		SerfMembers:       []string{*SerfMember},
	}

	a, err := agent.NewAgent(config)

	if err != nil {
		os.Exit(1)
	}

	go a.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Info().Msg("Program exited")
}
